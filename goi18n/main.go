package main

import (
	"flag"
	"fmt"
)

func main() {
	flag.Parse()
	e := NewExtractor()
	e.ExtractFiles(flag.Args())
	for _, m := range e.Messages() {
		fmt.Println(m.Content, m.Context)
	}
}
