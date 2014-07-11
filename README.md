go-i18n [![Build Status](https://secure.travis-ci.org/nicksnyder/go-i18n.png?branch=master)](http://travis-ci.org/nicksnyder/go-i18n)
=======

go-i18n is a Go [package](#i18n-package) and a [command](#goi18n-command) that can be used to translate Go programs into multiple languages.
* Supports pluralized strings using [CLDR plural rules](http://cldr.unicode.org/index/cldr-spec/plural-rules).
* Supports strings with named variables using [text/template](http://golang.org/pkg/text/template/) syntax.
* Translation files are simple JSON.
* [Documented](http://godoc.org/github.com/nicksnyder/go-i18n) and [tested](https://travis-ci.org/nicksnyder/go-i18n)!

Package i18n [![GoDoc](http://godoc.org/github.com/nicksnyder/go-i18n?status.png)](http://godoc.org/github.com/nicksnyder/go-i18n/i18n)
------------

The i18n package provides runtime APIs for fetching translated strings.

Command goi18n [![GoDoc](http://godoc.org/github.com/nicksnyder/go-i18n?status.png)](http://godoc.org/github.com/nicksnyder/go-i18n/goi18n)
--------------

The goi18n command provides functionality for managing the translation process.

### Installation

Make sure you have [setup GOPATH](http://golang.org/doc/code.html#GOPATH).

    go get -u github.com/nicksnyder/go-i18n/goi18n
    goi18n -help

### Workflow

A typical workflow looks like this:

1. Add a new string to your source code.

    ```go
    T("settings_title")
    ```

2. Add the string to en-US.all.json

    ```json
    [
      {
        "id": "settings_title",
        "translation": "Settings"
      }
    ]
    ```

3. Run goi18n

    ```
    goi18n path/to/*.all.json
    ```

4. Send `path/to/*.untranslated.json` to get translated.
5. Run goi18n again to merge the translations

    ```sh
    goi18n path/to/*.all.json path/to/*.untranslated.json
    ```

Translation files
-----------------

A translation file stores translated and untranslated strings.

Example:

```json
[
  {
    "id": "d_days",
    "translation": {
      "one": "{{.Count}} day",
      "other": "{{.Count}} days"
    }
  },
  {
    "id": "my_height_in_meters",
    "translation": {
      "one": "I am {{.Count}} meter tall.",
      "other": "I am {{.Count}} meters tall."
    }
  },
  {
    "id": "person_greeting",
    "translation": "Hello {{.Person}}"
  },
  {
    "id": "person_unread_email_count",
    "translation": {
      "one": "{{.Person}} has {{.Count}} unread email.",
      "other": "{{.Person}} has {{.Count}} unread emails."
    }
  },
  {
    "id": "person_unread_email_count_timeframe",
    "translation": {
      "one": "{{.Person}} has {{.Count}} unread email in the past {{.Timeframe}}.",
      "other": "{{.Person}} has {{.Count}} unread emails in the past {{.Timeframe}}."
    }
  },
  {
    "id": "program_greeting",
    "translation": "Hello world"
  },
  {
    "id": "your_unread_email_count",
    "translation": {
      "one": "You have {{.Count}} unread email.",
      "other": "You have {{.Count}} unread emails."
    }
  }
]
```

Supported languages
-------------------

* Arabic (`ar`)
* Bulgarian (`bg`)
* Catalan (`ca`)
* Chinese (simplified and traditional) (`zh`)
* Czech (`cs`)
* Danish (`da`)
* Dutch (`nl`)
* English (`en`)
* French (`fr`)
* German (`de`)
* Icelandic (`is`)
* Italian (`it`)
* Japanese (`ja`)
* Lithuanian (`lt`)
* Portuguese (`pt`)
* Portuguese (Brazilian) (`pt-BR`)
* Spanish (`es`)
* Swedish (`sv`)

Adding new languages
--------------------

It is easy to add support for additional languages:

1. Lookup the language's [CLDR plural rules](http://www.unicode.org/cldr/charts/latest/supplemental/language_plural_rules.html).
2. Add the language to [pluralspec.go](i18n/language/pluralspec.go):

    ```go
    var pluralSpecs = map[string]*PluralSpec{
        // ...
				// English
				"en": &PluralSpec{
					Plurals: newPluralSet(One, Other),
					PluralFunc: func(ops *operands) Plural {
						if ops.I == 1 && ops.V == 0 {
							return One
						}
						return Other
					},
				},
        // ...
    }
    ```

3. Add a test to [pluralspec_test.go](i18n/language/pluralspec_test.go)
4. Update this README with the new language.
5. Submit a pull request!

License
-------
go-i18n is available under the MIT license. See the [LICENSE](LICENSE) file for more info.
