package main

import (
	"flag"
	"fmt"
	"github.com/nicksnyder/go-i18n/src/pkg/extract"
	"os"
)

const usageTemplate = `Usage: goi18n command [arguments]

goi18n manages translations for a Go project.

The commands are:
    extract    extracts messages from Go source files
    merge      merges multiple message files into a single file
    format     formats message files into Go source files

Use "goi18n command -help" for more information about a command.
`

func usage() {
	fmt.Print(usageTemplate)
	os.Exit(2)
}

func main() {
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		usage()
	}

	switch args[0] {
	case "extract":
		extract.Run(args[1:])
	case "merge":
		fmt.Fprintln(os.Stderr, "Not implemented (yet)")
		//merge(args[1:])
	case "format":
		fmt.Fprintln(os.Stderr, "Not implemented (yet)")
		//format(args[1:])
	default:
		usage()
	}
}
