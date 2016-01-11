// Package nsi18n is a variation of the package i18n to support for multiple translation namespaces.
//
// Initialization
//
// Your Go program should load translations during its initialization.
//     nsi18n.CreateNamespace ("ns1","ns2","ns3")
//     nsi18n.MustLoadTranslationFile("ns1","path/to/fr-FR.ns1.json")
//     nsi18n.MustLoadTranslationFile("ns2","path/to/fr-FR.ns2.json")
//     nsi18n.MustLoadTranslationFile("ns3","path/to/fr-FR.ns3.json")
// If your translations are in a file format not supported by (Must)?LoadTranslationFile,
// then you can use the AddTranslation function to manually add translations.
//
// Fetching a translation
//
// Use Tfunc or MustTfunc to fetch a TranslateFunc that will return the translated string for a specific
// language for specific namespaces.
//     func handleRequest(w http.ResponseWriter, r *http.Request) {
//         cookieLang := r.Cookie("lang")
//         acceptLang := r.Header.Get("Accept-Language")
//         defaultLang = "en-US"  // known valid language
//         namespaces := []string{"ns2","ns3"}
//         T, err := nsi18n.Tfunc(namespaces,cookieLang, acceptLang, defaultLang)
//         fmt.Println(T("Hello world"))
//     }
//
// Usually it is a good idea to identify strings by a generic id rather than the English translation,
// but the rest of this documentation will continue to use the English translation for readability.
//     T("Hello world")     // ok
//     T("programGreeting") // better!
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
//     T("You have {{.Count}} unread emails.", 2)
//     T("I am {{.Count}} meters tall.", "1.7")
//
// Plural strings may also have variables.
//     T("{{.Person}} has {{.Count}} unread emails", 2, map[string]interface{}{
//         "Person": "Bob",
//     })
//
// Sentences with multiple plural components can be supported with nesting.
//     T("{{.Person}} has {{.Count}} unread emails in the past {{.Timeframe}}.", 3, map[string]interface{}{
//         "Person":    "Bob",
//         "Timeframe": T("{{.Count}} days", 2),
//     })
//
// Templates
//
// You can use the .Funcs() method of a text/template or html/template to register a TranslateFunc
// for usage inside of that template.
package nsi18n

import (
	"fmt"

	"github.com/nicksnyder/go-i18n/i18n/bundle"
	"github.com/nicksnyder/go-i18n/i18n/language"
	"github.com/nicksnyder/go-i18n/i18n/translation"
)

// TranslateFunc returns the translation of the string identified by translationID.
//
// If there is no translation for translationID, then the translationID itself is returned.
// This makes it easy to identify missing translations in your app.
//
// If translationID is a non-plural form, then the first variadic argument may be a map[string]interface{}
// or struct that contains template data.
//
// If translationID is a plural form, then the first variadic argument must be an integer type
// (int, int8, int16, int32, int64) or a float formatted as a string (e.g. "123.45").
// The second variadic argument may be a map[string]interface{} or struct that contains template data.
type TranslateFunc func(translationID string, args ...interface{}) string

// IdentityTfunc returns a TranslateFunc that always returns the translationID passed to it.
//
// It is a useful placeholder when parsing a text/template or html/template
// before the actual Tfunc is available.
func IdentityTfunc() TranslateFunc {
	return func(translationID string, args ...interface{}) string {
		return translationID
	}
}

var defaultbundles = make(map[string]*bundle.Bundle)

// CreateNamespace creates translation namespaces.
func CreateNamespace (namespaces ...string) {
	for _, namespace := range namespaces {
		defaultbundles[namespace] = bundle.New()
	}
}

// MustLoadTranslationFile is similar to LoadTranslationFile
// except it panics if an error happens.
func MustLoadTranslationFile(namespace string, filename string) error {
	bundle := defaultbundles[namespace]
	if bundle == nil {
		return fmt.Errorf("there is no namespace %s", namespace)
	}
	bundle.MustLoadTranslationFile(filename)
	return nil
}

// LoadTranslationFile loads the translations from filename in a specific namespace.
//
// The language that the translations are associated with is parsed from the filename (e.g. en-US.json).
//
// Generally you should load translation files once during your program's initialization.
func LoadTranslationFile(namespace string, filename string) error {
	bundle := defaultbundles[namespace]
	if bundle == nil {
		return fmt.Errorf("there is no namespace %s", namespace)
	}
	return bundle.LoadTranslationFile(filename)
}

// ParseTranslationFileBytes is similar to LoadTranslationFile except it parses the bytes in buf.
//
// It is useful for parsing translation files embedded with go-bindata.
func ParseTranslationFileBytes(namespace string, filename string, buf []byte) error {
	bundle := defaultbundles[namespace]
	if bundle == nil {
		return fmt.Errorf("there is no namespace %s", namespace)
	}
	return bundle.ParseTranslationFileBytes(filename, buf)
}

// AddTranslation adds translations for a language in a specific namespace.
//
// It is useful if your translations are in a format not supported by LoadTranslationFile.
func AddTranslation(namespace string, lang *language.Language, translations ...translation.Translation) error {
	bundle := defaultbundles[namespace]
	if bundle == nil {
		return fmt.Errorf("there is no namespace %s", namespace)
	}
	bundle.AddTranslation(lang, translations...)
	return nil
}

// LanguageTags returns the tags of all languages that have been added in a specific namespace.
func LanguageTags(namespace string) ([]string, error) {
	bundle := defaultbundles[namespace]
	if bundle == nil {
		return nil, fmt.Errorf("there is no namespace %s", namespace)
	}
	return bundle.LanguageTags(), nil
}

// LanguageTranslationIDs returns the ids of all translations that have been added for a given language in
// a specific namespace.
func LanguageTranslationIDs(namespace string, languageTag string) ([]string, error) {
	bundle := defaultbundles[namespace]
	if bundle == nil {
		return nil, fmt.Errorf("there is no namespace %s", namespace)
	}
	return bundle.LanguageTranslationIDs(languageTag), nil
}

// MustTfunc is similar to Tfunc except it panics if an error happens.
func MustTfunc(namespaces []string, languageSource string, languageSources ...string) TranslateFunc {
	var tfuncs []TranslateFunc
	for _, namespace := range namespaces {
		bundle := defaultbundles[namespace]
		if bundle == nil {
			panic (fmt.Errorf("there is no namespace %s", namespace))
		}
		tfunc, err := bundle.Tfunc(languageSource, languageSources...)
		if err != nil {
			panic (err)
		}
		tfuncs = append (tfuncs, TranslateFunc(tfunc))
	}
	return func(translationID string, args ...interface{}) string {
		str := translationID
		for _, tfunc := range tfuncs {
			str = tfunc(translationID, args...)
			if str != translationID {
				return str
			}
		}
		return str
	}
}

// Tfunc returns a TranslateFunc that will be bound to the first language which
// has a non-zero number of translations in a set of specific namespaces.
//
// The order of the namespaces in the first argument  defines the order of how translations
// are searched. In the bellow example, T searches a translation first in "ns2" and then in "ns3". It
// returns the first found translation.
//         namespaces := []string{"ns2","ns3"}
//         T, err := nsi18n.Tfunc(namespaces,"pt-BR")
//
// It can parse languages from Accept-Language headers (RFC 2616).
func Tfunc(namespaces []string, languageSource string, languageSources ...string) (TranslateFunc, error) {
	var tfuncs []TranslateFunc
	for _, namespace := range namespaces {
		bundle := defaultbundles[namespace]
		if bundle == nil {
			return nil, fmt.Errorf("there is no namespace %s", namespace)
		}
		tfunc, err := bundle.Tfunc(languageSource, languageSources...)
		if err != nil {
			return TranslateFunc(tfunc), err
		}
		tfuncs = append (tfuncs, TranslateFunc(tfunc))
	}
	return func(translationID string, args ...interface{}) string {
		str := translationID
		for _, tfunc := range tfuncs {
			str = tfunc(translationID, args...)
			if str != translationID {
				return str
			}
		}
		return str
	}, nil
}

// MustTfuncAndLanguage is similar to TfuncAndLanguage except it panics if an error happens.
func MustTfuncAndLanguage(namespaces []string, languageSource string, languageSources ...string) (TranslateFunc, []*language.Language) {
	var tfuncs []TranslateFunc
	var langs []*language.Language
	for _, namespace := range namespaces {
		bundle := defaultbundles[namespace]
		if bundle == nil {
			panic (fmt.Errorf("there is no namespace %s", namespace))
		}
		tfunc, lang, err := bundle.TfuncAndLanguage(languageSource, languageSources...)
		if err != nil {
			panic (err)
		}
		tfuncs = append (tfuncs, TranslateFunc(tfunc))
		langs = append (langs, lang)
	}
	return func(translationID string, args ...interface{}) string {
		str := translationID
		for _, tfunc := range tfuncs {
			str = tfunc(translationID, args...)
			if str != translationID {
				return str
			}
		}
		return str
	}, langs
}

// TfuncAndLanguage is similar to Tfunc except it also returns the language which TranslateFunc is bound to.
func TfuncAndLanguage(namespaces []string, languageSource string, languageSources ...string) (TranslateFunc, []*language.Language, error) {
	var tfuncs []TranslateFunc
	var langs []*language.Language
	for _, namespace := range namespaces {
		bundle := defaultbundles[namespace]
		if bundle == nil {
			return nil, langs, fmt.Errorf("there is no namespace %s", namespace)
		}
		tfunc, lang, err := bundle.TfuncAndLanguage(languageSource, languageSources...)
		if err != nil {
			return TranslateFunc(tfunc), langs, err
		}
		tfuncs = append (tfuncs, TranslateFunc(tfunc))
		langs = append (langs, lang)
	}
	return func(translationID string, args ...interface{}) string {
		str := translationID
		for _, tfunc := range tfuncs {
			str = tfunc(translationID, args...)
			if str != translationID {
				return str
			}
		}
		return str
	}, langs, nil
}
