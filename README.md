go-i18n
=======

go-i18n is a set of tools that can be used to translate Go programs into multiple lanaguages.

This project is experimental and still under development.

Go versions supported
---------------------

The project is being developed against the current release of Go:
http://golang.org/doc/devel/release.html#r60

Installation
------------

Run the following commands to install the go-i18n tools.
    goinstall -u github.com/nicksnyder/go-i18n/i18n
    goinstall -u github.com/nicksnyder/go-i18n/goi18nextract

i18n package
------------

The i18n package provides runtime APIs for looking up translated strings. An example is provided in example/

goi18nextract command
---------------------

goi18nextract extracts messages from Go source files. Run the command with no arguments for usage information.
