package main

import (
	"flag"
	"fmt"
	"json"
	"os"
)

var extractFlags *flag.FlagSet

func extractUsage() {
	errorf("Usage: %s extract file.go ...", os.Args[0])
	extractFlags.PrintDefaults()
	os.Exit(2)
}

func extract() {
	extractFlags = flag.NewFlagSet("extract", flag.ExitOnError)
	extractFlags.Usage = extractUsage
	format := extractFlags.String("format", "json", "The format used to output messages. Supported formats: json")
	filename := extractFlags.String("output", "", "The name of the file to write to. If not specified, output is written to stdout.")
	extractFlags.Parse(flag.Args()[1:])

	files := extractFlags.Args()
	if len(files) == 0 {
		exitf("No Go source files provided")
	}

	switch *format {
	case "json":
	default:
		exitf("Unknown format: %s", *format)
	}

	e := NewExtractor()
	e.ExtractFiles(files)

	var output []byte
	switch *format {
	case "json":
		output = toJson(e.Messages())
	default:
		panic("Unknown format: " + *format)
	}

	if *filename == "" {
		fmt.Print(string(output))
		return
	}

	file, err := os.Create(*filename)
	if err != nil {
		exitf(err.String())
	}
	defer file.Close()

	if _, err := file.Write(output); err != nil {
		panic(err.String())
	}
}

func toJson(msgs []Message) []byte {
	json, err := json.MarshalIndent(msgs, "", "  ")
	if err != nil {
		panic(err.String())
	}
	return json
}
