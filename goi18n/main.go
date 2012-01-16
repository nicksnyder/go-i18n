package main

import (
	"csv"
	"flag"
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

type Formatter func(io.Writer, []Message)

func usage() {
	Printf("Usage: %s file.go ...", os.Args[0])
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	flag.Usage = usage
	flag.StringVar(&format, "format", "json", "The format used to output messages. Supported formats: json, csv")
	flag.StringVar(&filename, "output", "", "The name of the file to write to. If not specified, output is written to stdout.")
	flag.StringVar(&csvdelim, "csvdelim", ",", "The field delimiter to use for csv files.")
	flag.BoolVar(&csvcrlf, "csvcrlf", false, `True to use \r\n as the line terminator. Default is \n.`)
	flag.Parse()

	files := flag.Args()
	if len(files) == 0 {
		Exitf("No Go source files provided.")
	}

	var f Formatter
	switch format {
	case "json":
		f = jsonFormatter
	case "csv":
		f = csvFormatter
	default:
		Exitf("Unknown format: %s", format)
	}

	var w io.Writer
	if filename == "" {
		w = os.Stdout
	} else {
		file, err := os.Create(filename)
		if err != nil {
			Exitf(err.String())
		}
		defer file.Close()
		w = file
	}

	e := NewExtractor()
	e.ExtractFiles(files)
	f(w, e.Messages())
}

func jsonFormatter(w io.Writer, msgs []Message) {
	json, err := json.MarshalIndent(msgs, "", "  ")
	if err != nil {
		Exitf(err.String())
	}
	if _, err := w.Write(json); err != nil {
		Exitf(err.String())
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
			Exitf(err.String())
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
