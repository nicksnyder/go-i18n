go-i18n [![Build Status](https://secure.travis-ci.org/nicksnyder/go-i18n.png?branch=master)](http://travis-ci.org/nicksnyder/go-i18n)
=======

go-i18n is a Go [package](#i18n-package) and a [command](#goi18n-command) that can be used to translate Go programs into multiple lanaguages.

Requires Go 1.2.

Features
--------

* Implements [CLDR plural rules](http://cldr.unicode.org/index/cldr-spec/plural-rules).
* Uses [text/template](http://golang.org/pkg/text/template/) for parameter substitution.
* Translation files are simple JSON.

i18n package
------------

The i18n package provides runtime APIs for looking up translated strings.

```go
import "github.com/nicksnyder/go-i18n/i18n"
```

##### Loading translations

Load translation files during your program's initialization.
The name of a translation file must contain a [language tag](http://en.wikipedia.org/wiki/IETF_language_tag) that is supported by i18n.

```go
i18n.MustLoadTranslationFile("path/to/fr-FR.all.json")
```

##### Selecting a locale

Tfunc returns a function that can lookup the translation of a string for that locale.
It accepts multiple locale parameters so you can gracefully fallback to other locales.

```go
userLocale = "ar-AR"       // e.g. from user preference, accept header, or language cookie
defaultLocale = "en-US"    // known valid locale
T, err := i18n.Tfunc(userLocale, defaultLocale)
```

##### Loading a string translation

Use the translation function to fetch the translation of a string.

```go
fmt.Println(T("Hello world"))
```

Usually it is a good idea to identify strings by a generic id rather than the English translation, but the rest of this document will continue to use the English translation for readability.

```go
T("program_greeting")
```

##### Strings with variables

You can have variable substitutions in your string using [text/template](http://golang.org/pkg/text/template/) syntax.

```go
T("Hello {{.Person}}", map[string]interface{}{
	"Person": "Bob",
}))
```

##### Plural strings

Each language handles pluralization differently. A few examples:
* English treats one as singular and all other numbers as plural (e.g. 0 cats, 1 cat, 2 cats).
* French treats zero and one as singular and all other numbers as pural (e.g. 0 chat, 1 chat, 2 chats)
* Arabic has six different plural forms!

The translation function handles [all of this logic](http://www.unicode.org/cldr/charts/latest/supplemental/language_plural_rules.html) for you.

```go
T("You have {{.Count}} unread emails", 2)
```

With variable substitutions:

```go
fmt.Println(T("{{.Person}} has {{.Count}} unread emails", 2, map[string]interface{}{
	"Person": "Bob",
}))
```

Sentences with multiple plural components can be supported with nesting.

```go
fmt.Println(T("{{.Person}} has {{.Count}} unread email in the past {{.Timeframe}}.", 3, map[string]interface{}{
	"Person":    "Bob",
	"Timeframe": T("{{.Count}} days", 2),
}))
```

A complete example is [here](i18n/example_test.go).

goi18n command
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
* English (`en`)
* French (`fr`)
* Japanese (`ja`)
* Spanish (`es`)

More languages are straightforward to add:

1. Lookup the language's [CLDR plural rules](http://www.unicode.org/cldr/charts/latest/supplemental/language_plural_rules.html).
2. Add the language to [language.go](i18n/language.go):

    ```go
    RegisterLanguage(&Language{
        Code:             "en",
        Name:             "English",
        PluralCategories: newSet(One, Other),
        IntFunc: func(i int64) PluralCategory {
            if i == 1 {
                return One
            }
            return Other
        },
        FloatFunc: func(f float64) PluralCategory {
            return Other
        },
    })
    ```

3. Add a test to [language_test.go](i18n/language_test.go)
4. Submit a pull request!

License
-------
go-i18n is available under the MIT license. See the [LICENSE](LICENSE) file for more info.
