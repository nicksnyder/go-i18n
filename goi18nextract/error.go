package main

import (
	"flag"
	"fmt"
	"os"
)

func Printf(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", a...)
}

func Exitf(format string, a ...interface{}) {
	Printf(format, a...)
	flag.Usage()
}
