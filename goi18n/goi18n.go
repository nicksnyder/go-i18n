package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	flag.Usage = usage

	mergeCmd := flag.NewFlagSet("merge", flag.ExitOnError)
	mergeCmd.Usage = usageMerge
	sourceLanguage := mergeCmd.String("sourceLanguage", "en-us", "")
	outdir := mergeCmd.String("outdir", ".", "")
	format := mergeCmd.String("format", "json", "")

	constantsCmd := flag.NewFlagSet("constants", flag.ExitOnError)
	constantsCmd.Usage = usageConstants
	packageName := constantsCmd.String("package", "R", "")
	outdirConstants := constantsCmd.String("outdir", ".", "")

	if len(os.Args) == 1 {
		usage()
	}

	switch os.Args[1] {
	case "merge":
		mergeCmd.Parse(os.Args[2:])
	case "constants":
		constantsCmd.Parse(os.Args[2:])
	default:
		mergeCmd.Parse(os.Args[1:])
	}

	if mergeCmd.Parsed() {
		if len(constantsCmd.Args()) != 1 {
			fmt.Println("need at least one translation file to parse")
			usageMerge()
		}

		mc := &mergeCommand{
			translationFiles:  mergeCmd.Args(),
			sourceLanguageTag: *sourceLanguage,
			outdir:            *outdir,
			format:            *format,
		}
		if err := mc.execute(); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	} else if constantsCmd.Parsed() {
		if len(constantsCmd.Args()) != 1 {
			fmt.Println("need one translation file")
			usageConstants()
		}

		cc := &constantsCommand{
			translationFile: constantsCmd.Args()[0],
			packageName:     *packageName,
			outdir:          *outdirConstants,
		}
		if err := cc.execute(); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}
}

func usage() {
	fmt.Printf(`goi18n tools for translation files.

Usage:

    goi18n merge     [options] [files...]
    goi18n constants [options] [file]

For more details execute

    goi18n [command] -help

`)
	os.Exit(1)
}

func usageMerge() {
	fmt.Printf(`Merge translation files.

Usage:

    goi18n merge [options] [files...]

Translation files:

    A translation file contains the strings and translations for a single language.

    Translation file names must have a suffix of a supported format (e.g. .json) and
    contain a valid language tag as defined by RFC 5646 (e.g. en-us, fr, zh-hant, etc.).

    For each language represented by at least one input translation file, goi18n will produce 2 output files:

        xx-yy.all.format
            This file contains all strings for the language (translated and untranslated).
            Use this file when loading strings at runtime.

        xx-yy.untranslated.format
            This file contains the strings that have not been translated for this language.
            The translations for the strings in this file will be extracted from the source language.
            After they are translated, merge them back into xx-yy.all.format using goi18n.

Merging:

    goi18n will merge multiple translation files for the same language.
    Duplicate translations will be merged into the existing translation.
    Non-empty fields in the duplicate translation will overwrite those fields in the existing translation.
    Empty fields in the duplicate translation are ignored.

Adding a new language:

    To produce translation files for a new language, create an empty translation file with the
    appropriate name and pass it in to goi18n.

Options:

    -sourceLanguage tag
        goi18n uses the strings from this language to seed the translations for other languages.
        Default: en-us

    -outdir directory
        goi18n writes the output translation files to this directory.
        Default: .

    -format format
        goi18n encodes the output translation files in this format.
        Supported formats: json, yaml
        Default: json

`)
	os.Exit(1)
}

func usageConstants() {
	fmt.Printf(`Generate constant file from translation file.

Usage:

    goi18n constants [options] [file]

Translation files:

    A translation file contains the strings and translations for a single language.

    Translation file names must have a suffix of a supported format (e.g. .json) and
    contain a valid language tag as defined by RFC 5646 (e.g. en-us, fr, zh-hant, etc.).

Options:

    -package name
        goi18n generates the constant file under the package name.
        Default: R

    -outdir directory
        goi18n writes the constant file to this directory.
        Default: .

`)
	os.Exit(1)
}
