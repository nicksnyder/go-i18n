package extract

import (
	"flag"
	"fmt"
	"github.com/nicksnyder/go-i18n/src/pkg/msg"
	"github.com/nicksnyder/go-i18n/src/pkg/csv"
	"github.com/nicksnyder/go-i18n/src/pkg/goio"
	"github.com/nicksnyder/go-i18n/src/pkg/json"
	"io"
	"os"
	"utf8"
)

// Flags
var (
	format      string
	filename    string
	csvfieldsep string
	csvcrlf     bool
)

var flags *flag.FlagSet

func usage() {
	errorf("Usage: %s extract file.go ...\n", os.Args[0])
	flags.PrintDefaults()
	os.Exit(2)
}

func Run(args []string) {
	flags = flag.NewFlagSet("extract", flag.ExitOnError)
	flags.Usage = usage
	flags.StringVar(&format, "format", "json", "The format used to output messages. Supported formats: json, csv")
	flags.StringVar(&filename, "output", "", "The name of the file to write to. If not specified, output is written to stdout.")
	flags.StringVar(&csvfieldsep, "csvfieldsep", ",", "The field delimiter to use in csv files.")
	flags.BoolVar(&csvcrlf, "csvcrlf", false, `True to use \r\n as the line terminator in csv files. Default is \n.`)
	flags.Parse(args)

	filenames := flags.Args()
	if len(filenames) < 1 {
		usage()
	}

	var w msg.Writer
	switch format {
	case "json":
		w = json.NewWriter()
	case "csv":
		w = csv.NewWriter(getCsvFieldSep(), csvcrlf)
	default:
		exitf("Unsupported format: %s", format)
	}

	var dst io.Writer
	if filename == "" {
		dst = os.Stdout
	} else {
		file, err := os.Create(filename)
		if err != nil {
			exitf(err.String())
		}
		defer file.Close()
		dst = file
	}

	r := goio.NewReader()
	bundle := msg.NewBundle()
	for _, filename := range filenames {
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

		bundle.AddMessages(m)
	}

	w.WriteMessages(dst, bundle.Messages())
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
