package main

import (
	"flag"
	"fmt"
	"json"
	"os"
)

var extractFlags *flag.FlagSet

func extractUsage() {
	errorf("Usage: %s extract file.go...", os.Args[0])
	extractFlags.PrintDefaults()
	os.Exit(2)
}

func extract() {
	extractFlags = flag.NewFlagSet("extract", flag.ExitOnError)
	extractFlags.Usage = extractUsage
	format := extractFlags.String("format", "json", "The format used to output messages (json)")
	extractFlags.Parse(flag.Args()[1:])

	files := extractFlags.Args()
	if len(files) == 0 {
		errorf("No Go source files provided")
		extractUsage()
	}

	switch *format {
	case "json":
	default:
		errorf("Unknown format: %s", *format)
		extractUsage()
	}

	e := NewExtractor()
	e.ExtractFiles(files)
	switch *format {
	case "json":
		printJson(e.Messages())
	}
}

func printJson(msgs []Message) {
	json, err := json.MarshalIndent(msgs, "", "  ")
	if err != nil {
		fmt.Fprint(os.Stderr, err.String())
		return
	}
	fmt.Print(string(json))
}
