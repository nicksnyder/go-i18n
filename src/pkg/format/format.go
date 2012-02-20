package format

import (
	"flag"
	"fmt"
	"github.com/nicksnyder/go-i18n/src/pkg/msg"
	"github.com/nicksnyder/go-i18n/src/pkg/csv"
	"github.com/nicksnyder/go-i18n/src/pkg/json"
	"os"
	"path/filepath"
	"regexp"
	"utf8"
)

// Flags
var (
	outformat   string
	outdir      string
	csvfieldsep string
	csvcrlf     bool
)

// Readers
var (
	jsonReader = json.NewReader()
	csvReader  = csv.NewReader()
)

var localeRegexp = regexp.MustCompile(`[a-z][a-z][_\-][A-Z][A-Z]`)
var bundles = make(map[string]*msg.Bundle)
var flags *flag.FlagSet

func usage() {
	errorf("Usage: %s format [OPTIONS] FILE ...\n", os.Args[0])
	flags.PrintDefaults()
	os.Exit(2)
}

func Run(args []string) {
	flags = flag.NewFlagSet("format", flag.ExitOnError)
	flags.Usage = usage
	flags.StringVar(&outformat, "format", "json", "The format used to output messages. Supported formats: json, csv")
	flags.StringVar(&outdir, "dir", ".", "The directory where output files will be written.")
	flags.StringVar(&csvfieldsep, "csvfieldsep", ",", "The field delimiter to use in csv files.")
	flags.BoolVar(&csvcrlf, "csvcrlf", false, `True to use \r\n as the line terminator in csv files. Defaults to \n.`)
	flags.Parse(args)

	filenames := flags.Args()
	if len(filenames) < 1 {
		usage()
	}

	var w msg.Writer
	switch outformat {
	case "json":
		w = json.NewWriter()
	case "csv":
		w = csv.NewWriter(getCsvFieldSep(), csvcrlf)
	default:
		exitf("Unsupported output format: %s\n", outformat)
	}

	for _, filename := range filenames {
		locale, informat := parseFilename(filename)
		if locale == "" {
			errorf("Skipping file with unknown locale: %s\n", filename)
			continue
		}

		var r msg.Reader
		switch informat {
		case "json":
			r = jsonReader
		case "csv":
			r = csvReader
		default:
			errorf("Skipping file with unsupported format: %s\n", filename)
			continue
		}

		file, err := os.Open(filename)
		if err != nil {
			errorf(err.String())
			continue
		}

		m, err := r.ReadMessages(file)
		if err != nil {
			errorf(err.String())
			continue
		}

		bundle := getBundle(locale)
		bundle.AddMessages(m)
	}

	for locale, bundle := range bundles {
		writeFile(w, "all", locale, bundle.Messages())
		writeFile(w, "untranslated", locale, bundle.UntranslatedMessages())
		writeFile(w, "translated", locale, bundle.TranslatedMessages())
	}
}

func writeFile(w msg.Writer, prefix, locale string, msgs []msg.Message) {
	filename := filepath.Join(outdir, fmt.Sprintf("%s.%s.%s", prefix, locale, outformat))
	file, err := os.Create(filename)
	if err != nil {
		errorf(err.String())
		return
	}
	defer file.Close()
	w.WriteMessages(file, msgs)
}

func parseFilename(filename string) (locale, format string) {
	matches := localeRegexp.FindAllString(filename, -1)
	if matches == nil {
		return
	}
	locale = matches[len(matches)-1]
	format = filepath.Ext(filename)[1:]
	return
}

func getBundle(locale string) *msg.Bundle {
	_, found := bundles[locale]
	if !found {
		bundles[locale] = msg.NewBundle()
	}
	return bundles[locale]
}

func errorf(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
}

func exitf(format string, a ...interface{}) {
	errorf(format, a...)
	os.Exit(2)
}

func getCsvFieldSep() int {
	if csvfieldsep == "tab" {
		return '\t'
	}
	rune, size := utf8.DecodeRuneInString(csvfieldsep)
	if size > 0 {
		return rune
	}
	return ','
}
