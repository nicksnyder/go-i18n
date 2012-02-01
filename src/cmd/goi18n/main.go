package main

import (
	"flag"
	"fmt"
	"github.com/nicksnyder/go-i18n/src/pkg/extract"
	"github.com/nicksnyder/go-i18n/src/pkg/format"
	"os"
)

const usageTemplate = `Usage: goi18n command [arguments]

goi18n manages translations for a Go project.

The commands are:
    extract    extracts messages from Go source files
    format     merges and formats message files

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
	case "format":
		format.Run(args[1:])
	default:
		usage()
	}
}
