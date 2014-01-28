go-i18n 
=======

go-i18n is a Go package (i18n) and a command (goi18n) that can be used to translate Go programs into multiple lanaguages.

Features
--------

* Implements [CLDR plural rules](http://cldr.unicode.org/index/cldr-spec/plural-rules).
* Uses [text/template](http://golang.org/pkg/text/template/) for parameter substitution.
* Translation files are simple JSON.

i18n package
------------

The i18n package provides runtime APIs for looking up translated strings.

A simple example:

```go
package main

import (
	"fmt"
	"github.com/nicksnyder/go-i18n/pkg/i18n"
)

func main() {
	i18n.MustLoadTranslationFile("path/to/fr-FR.all.json")
	
	T, err := i18n.Tfunc("ar-AR")
	if err != nil {
		T = i18n.MustTfunc("en-US") // fallback
	}
	
	// Regular string with no substitutions.
	fmt.Println(T("Hello world"))
	
	// String with variable substitutions.
	fmt.Println(T("Hello {{.Person}}", map[string]interface{}{
		"Person": "Bob",
	}))
	
	// Plural string.
	fmt.Println(T("You have {{.Count}} unread emails", 2))
	
	// Plural string with other substitutions.
	fmt.Println(T("{{.Person}} has {{.Count}} unread emails", 2, map[string]interface{}{
		"Person": "Bob",
	}))

	// Compound plural string.
	fmt.Println(T("{{.Person}} has {{.Count}} unread email in the past {{.Timeframe}}.", 3, map[string]interface{}{
		"Person":    "Bob",
		"Timeframe": T("{{.Count}} days", 2),
	}))
}
```

Usually it is a good idea to use generic ids for translations instead of the English string.


```go
T("program_greeting")
```

A more complete example is [here](pkg/i18n/example_test.go).

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

goi18n command
--------------

The goi18n command provides functionality for managing the translation process.

A typical workflow looks like this:

1. Add a new string to your source code.

    ```go
    T("some_page_title")
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

Languages
---------

Currently supported
* Arabic
* English
* French

More languages are straightforward to add:

1. Lookup the language's [CLDR plural rules](http://www.unicode.org/cldr/charts/latest/supplemental/language_plural_rules.html).
2. Add the language to [language.go](pkg/i18n/language.go):

    ```go
	RegisterLanguage(&Language{
		Code:             "ar",
		Name:             "العربية",
		PluralCategories: newSet(Zero, One, Two, Few, Many, Other),
		IntFunc: func(i int64) PluralCategory {
			switch i {
			case 0:
				return Zero
			case 1:
				return One
			case 2:
				return Two
			default:
				mod100 := i % 100
				if mod100 >= 3 && mod100 <= 10 {
					return Few
				}
				if mod100 >= 11 {
					return Many
				}
				return Other
			}
		},
		FloatFunc: func(f float64) PluralCategory {
			return Other
		},
	})
    ```

3. Submit a pull request!
