# go-i18n ![Build status](https://github.com/nicksnyder/go-i18n/workflows/Build/badge.svg) [![Report card](https://goreportcard.com/badge/github.com/nicksnyder/go-i18n)](https://goreportcard.com/report/github.com/nicksnyder/go-i18n) [![codecov](https://codecov.io/gh/nicksnyder/go-i18n/branch/master/graph/badge.svg)](https://codecov.io/gh/nicksnyder/go-i18n) [![Sourcegraph](https://sourcegraph.com/github.com/nicksnyder/go-i18n/-/badge.svg)](https://sourcegraph.com/github.com/nicksnyder/go-i18n?badge)

go-i18n is a Go [package](#package-i18n) and a [command](#command-goi18n) that helps you translate Go programs into multiple languages.

- Supports [pluralized strings](http://cldr.unicode.org/index/cldr-spec/plural-rules) for all 200+ languages in the [Unicode Common Locale Data Repository (CLDR)](https://www.unicode.org/cldr/charts/28/supplemental/language_plural_rules.html).
  - Code and tests are [automatically generated](https://github.com/nicksnyder/go-i18n/tree/main/v2/internal/plural/codegen) from [CLDR data](http://cldr.unicode.org/index/downloads).
- Supports strings with named variables using [text/template](http://golang.org/pkg/text/template/) syntax.
- Supports message files of any format (e.g. JSON, TOML, YAML).

## Package i18n
[![GoDoc](https://godoc.org/github.com/nicksnyder/go-i18n?status.svg)](https://godoc.org/github.com/nicksnyder/go-i18n/v2/i18n)

The i18n package provides support for looking up messages according to a set of locale preferences.

```go
import "github.com/nicksnyder/go-i18n/v2/i18n"
```

Create a Bundle to use for the lifetime of your application.

```go
bundle := i18n.NewBundle(language.English)
```

Load translations into your bundle during initialization.

```go
bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
bundle.LoadMessageFile("es.toml")
```

Create a Localizer to use for a set of language preferences.

```go
func(w http.ResponseWriter, r *http.Request) {
    lang := r.FormValue("lang")
    accept := r.Header.Get("Accept-Language")
    localizer := i18n.NewLocalizer(bundle, lang, accept)
}
```

Use the Localizer to lookup messages.

```go
localizer.Localize(&i18n.LocalizeConfig{
    DefaultMessage: &i18n.Message{
        ID: "PersonCats",
        One: "{{.Name}} has {{.Count}} cat.",
        Other: "{{.Name}} has {{.Count}} cats.",
    },
    TemplateData: map[string]interface{}{
        "Name": "Nick",
        "Count": 2,
    },
    PluralCount: 2,
}) // Nick has 2 cats.
```

## Command goi18n
[![GoDoc](https://godoc.org/github.com/nicksnyder/go-i18n?status.svg)](https://godoc.org/github.com/nicksnyder/go-i18n/v2/goi18n)

The goi18n command manages message files used by the i18n package.

```
go get -u github.com/nicksnyder/go-i18n/v2/goi18n
goi18n -help
```

### Extracting messages

Use `goi18n extract` to extract all i18n.Message struct literals in Go source files to a message file for translation.

```toml
# active.en.toml
[PersonCats]
description = "The number of cats a person has"
one = "{{.Name}} has {{.Count}} cat."
other = "{{.Name}} has {{.Count}} cats."
```

### Translating a new language

1. Create an empty message file for the language that you want to add (e.g. `translate.es.toml`).
2. Run `goi18n merge active.en.toml translate.es.toml` to populate `translate.es.toml` with the messages to be translated.

   ```toml
   # translate.es.toml
   [HelloPerson]
   hash = "sha1-5b49bfdad81fedaeefb224b0ffc2acc58b09cff5"
   other = "Hello {{.Name}}"
   ```

3. After `translate.es.toml` has been translated, rename it to `active.es.toml`.

   ```toml
   # active.es.toml
   [HelloPerson]
   hash = "sha1-5b49bfdad81fedaeefb224b0ffc2acc58b09cff5"
   other = "Hola {{.Name}}"
   ```

4. Load `active.es.toml` into your bundle.

   ```go
   bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
   bundle.LoadMessageFile("active.es.toml")
   ```

### Translating new messages

If you have added new messages to your program:

1. Run `goi18n extract` to update `active.en.toml` with the new messages.
2. Run `goi18n merge active.*.toml` to generate updated `translate.*.toml` files.
3. Translate all the messages in the `translate.*.toml` files.
4. Run `goi18n merge active.*.toml translate.*.toml` to merge the translated messages into the active message files.

## For more information and examples:

- Read the [documentation](https://godoc.org/github.com/nicksnyder/go-i18n/v2).
- Look at the [code examples](https://github.com/nicksnyder/go-i18n/blob/main/v2/i18n/example_test.go) and [tests](https://github.com/nicksnyder/go-i18n/blob/main/v2/i18n/localizer_test.go).
- Look at an example [application](https://github.com/nicksnyder/go-i18n/tree/main/v2/example).

## License

go-i18n is available under the MIT license. See the [LICENSE](LICENSE) file for more info.
