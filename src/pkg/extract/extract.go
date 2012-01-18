package extract

import (
	"csv"
	"flag"
	"fmt"
	"io"
	"json"
	"os"
	"utf8"
)

// Flags
var (
	format   string
	filename string
	csvdelim string
	csvcrlf  bool
)

type formatter func(io.Writer, []Message)

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
	flags.StringVar(&csvdelim, "csvdelim", ",", "The field delimiter to use for csv files.")
	flags.BoolVar(&csvcrlf, "csvcrlf", false, `True to use \r\n as the line terminator. Default is \n.`)
	flags.Parse(args)

	files := flags.Args()
	if len(files) < 1 {
		usage()
	}

	var f formatter
	switch format {
	case "json":
		f = jsonFormatter
	case "csv":
		f = csvFormatter
	default:
		exitf("Unsupported format: %s", format)
	}

	var w io.Writer
	if filename == "" {
		w = os.Stdout
	} else {
		file, err := os.Create(filename)
		if err != nil {
			exitf(err.String())
		}
		defer file.Close()
		w = file
	}

	e := NewExtractor()
	e.ExtractFiles(files)
	f(w, e.Messages())
}

func errorf(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
}

func exitf(format string, a ...interface{}) {
	errorf(format, a...)
	os.Exit(2)
}

func jsonFormatter(w io.Writer, msgs []Message) {
	json, err := json.MarshalIndent(msgs, "", "  ")
	if err != nil {
		exitf(err.String())
	}
	if _, err := w.Write(json); err != nil {
		exitf(err.String())
	}
}

func csvFormatter(w io.Writer, msgs []Message) {
	c := csv.NewWriter(w)
	defer c.Flush()
	c.Comma = csvDelim()
	c.UseCRLF = csvcrlf
	c.Write([]string{"context", "content", "translation"})
	for _, m := range msgs {
		err := c.Write([]string{m.Context, m.Content, m.Translation})
		if err != nil {
			exitf(err.String())
		}
	}
}

func csvDelim() int {
	if csvdelim == "tab" {
		return '\t'
	}
	rune, size := utf8.DecodeRuneInString(csvdelim)
	if size > 0 {
		return rune
	}
	return ','
}
