package main

import (
	"flag"
	"fmt"
	"github.com/nicksnyder/go-i18n/pkg/i18n"
	"os"
)

func usage() {
	fmt.Printf(`goi18n formats and merges translation files.

Usage:

    goi18n [options] [files...]

Translation files:

    A translation file contains the strings and translations for a single locale (language + country).

    Translation file names must have a suffix of a supported format (e.g. .json) and
    contain a valid locale identifier (e.g. ar-EG, en-US, fr-FR, etc.).

    For each locale represented by at least one input translation file, goi18n will produce 2 output files:

        xx-XX.all.format
            This file contains all strings for the locale (translated and untranslated).

        xx-XX.untranslated.format
            This file contains the strings that have not been translated for this locale.
			The translations for the strings in this file will be extracted from the source locale.
            Get these strings translated! After they are translated, merge them back into
			xx-XX.all.format using goi18n.

    goi18n will merge multiple translation files for the same locale. 
    Duplicate translations will be merged into the existing translation.
    Non-empty fields in the duplicate translation will overwrite those fields in the existing translation.
    Empty fields in the duplicate translation are ignored.

    To produce translation files for a new locale, create an empty translation file with the
    appropriate name and pass it in to goi18n.

Options:

    -sourceLocale localeId
	    The id of the locale that strings are initially written in (e.g. xx-XX)
	    Default: en-US

    -outdir directory
        goi18n will write the output translation files to this directory.
        Default: .

    -format format
        goi18n will encode the output translation files in this format.
        Supported formats: json
        Default: json

`)
	os.Exit(1)
}

func main() {
	flag.Usage = usage
	sourceLocale := *flag.String("sourceLocale", "en-US", "")
	outdir := *flag.String("outdir", ".", "")
	format := *flag.String("format", "json", "")
	flag.Parse()
	if err := i18n.Merge(flag.Args(), sourceLocale, outdir, format); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
