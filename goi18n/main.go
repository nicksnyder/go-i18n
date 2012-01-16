package main

import (
	"flag"
	"fmt"
	"os"
)

const commands = `
Commands:
  extract - Extracts strings from Go source files.
`

func usage() {
	errorf("Usage: %s command", os.Args[0])
	flag.PrintDefaults()
	errorf(commands)
	os.Exit(2)
}

func main() {
	flag.Usage = usage
	flag.Parse()

	switch flag.Arg(0) {
	case "extract":
		extract()
	default:
		usage()
	}
}

func errorf(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", a...)
}
