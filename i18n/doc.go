// Package i18n supports string translations with variable substitution and CLDR pluralization.
// It is intended to be used in conjunction with github.com/nicksnyder/go-i18n/goi18n,
// although that is not strictly required.
//
// Initialization
//
// Your Go program should load translations during its intialization.
//     i18n.MustLoadTranslationFile("path/to/fr-FR.all.json")
// If your translations are in a file format not supported by (Must)?LoadTranslationFile,
// then you can use the AddTranslation function to manually add translations.
//
// Fetching a translation
//
// Use Tfunc or MustTfunc to fetch a TranslateFunc that will return the translated string for a specific locale.
// The TranslateFunc will be bound to the first valid locale passed to Tfunc.
//     userLocale = "ar-AR"     // user preference, accept header, language cookie
//     defaultLocale = "en-US"  // known valid locale
//     T, err := i18n.Tfunc(userLocale, defaultLocale)
//     fmt.Println(T("Hello world"))
//
// Usually it is a good idea to identify strings by a generic id rather than the English translation,
// but the rest of this documentation will continue to use the English translation for readability.
//     T("program_greeting")
//
// Variables
//
// TranslateFunc supports strings that have variables using the text/template syntax.
//     T("Hello {{.Person}}", map[string]interface{}{
//         "Person": "Bob",
//     })
//
// Pluralization
//
// TranslateFunc supports the pluralization of strings using the CLDR pluralization rules defined here:
// http://www.unicode.org/cldr/charts/latest/supplemental/language_plural_rules.html
//     T("You have {{.Count}} unread emails", 2)
//
// Plural strings may also have variables.
//     T("{{.Person}} has {{.Count}} unread emails", 2, map[string]interface{}{
//         "Person": "Bob",
//     })
//
// Compound plural strings can be created with nesting.
//     T("{{.Person}} has {{.Count}} unread emails in the past {{.Timeframe}}.", 3, map[string]interface{}{
//         "Person":    "Bob",
//         "Timeframe": T("{{.Count}} days", 2),
//     })
//
// Templates
//
// You can use the .Funcs() method of a text/template or html/template to register a TranslateFunc
// for usage inside of that template.
package i18n
