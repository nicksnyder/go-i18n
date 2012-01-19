go-i18n (INCOMPLETE, EXPERIMENTAL, AND UNDER DEVELOPMENT)
=========================================================

go-i18n is a set of tools that can be used to translate Go programs into multiple lanaguages.

Go versions supported
---------------------

The project is being developed against the current release of Go:
http://golang.org/doc/devel/release.html#r60

Installation
------------

Run the following commands to install the go-i18n tools.

    goinstall -u github.com/nicksnyder/go-i18n/src/pkg/i18n
    goinstall -u github.com/nicksnyder/go-i18n/src/cmd/goi18n

i18n package
------------

The i18n package provides runtime APIs for looking up translated strings. An example is provided in example/

goi18n command
--------------

The goi18n command provides functionality for managing the translation process.

A typical workflow looks like this:

1. Write Go code using the i18n package.

		package main
		
		import (
			"flag"
			"fmt"
			"github.com/nicksnyder/go-i18n/src/pkg/i18n"
		)
		
		var (
			HelloWorld   = i18n.NewMessage("Hello world!", "This message is displayed when the program begins")
			GoodbyeWorld = i18n.NewMessage("Goodbye world.", "This message is displayed when the program ends")
		)
		
		var locale string

		func main() {
			flag.StringVar(&locale, "locale", "", "The locale to use for translated messages.")
			flag.Parse()	
			i18n.SetLocale(locale)
			fmt.Println(HelloWorld.String())
			fmt.Println(GoodbyeWorld.String())
		}

2. Extract the message from the Go source files.

		goi18n extract -format=json main.go > messages.json

	or

		goi18n extract -format=json -output=messages.json main.go

3. Get the messages translated.

4. Merge the new translations with any existing translations.

		goi18n merge existing/es_ES.json new/es_ES.json

5. Format the translations into Go source files.

		goi18n format es_ES.json fr_FR.json ... # generates es_ES.go fr_FR.go ...

6. Compile your program with the generated source files.
