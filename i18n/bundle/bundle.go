// Package bundle manages translations for multiple languages.
package bundle

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	//	"launchpad.net/goyaml"

	"path/filepath"

	"github.com/nicksnyder/go-i18n/i18n/language"
	"github.com/nicksnyder/go-i18n/i18n/translation"
)

// TranslateFunc is a copy of i18n.TranslateFunc to avoid a circular dependency.
type TranslateFunc func(translationID string, args ...interface{}) string

// Bundle stores the translations for multiple languages.
type Bundle struct {
	translations map[string]map[string]translation.Translation
}

// New returns an empty bundle.
func New() *Bundle {
	return &Bundle{
		translations: make(map[string]map[string]translation.Translation),
	}
}

// MustLoadTranslationFile is similar to LoadTranslationFile
// except it panics if an error happens.
func (b *Bundle) MustLoadTranslationFile(filename string) {
	if err := b.LoadTranslationFile(filename); err != nil {
		panic(err)
	}
}

// LoadTranslationFile loads the translations from filename into memory.
//
// The language that the translations are associated with is parsed from the filename (e.g. en-US.json).
//
// Generally you should load translation files once during your program's initialization.
func (b *Bundle) LoadTranslationFile(filename string) error {
	basename := filepath.Base(filename)
	langs := language.Parse(basename)
	switch l := len(langs); {
	case l == 0:
		return fmt.Errorf("no language found in %q", basename)
	case l > 1:
		return fmt.Errorf("multiple languages found in filename %q: %v; expected one", basename, langs)
	}
	translations, err := parseTranslationFile(filename)
	if err != nil {
		return err
	}
	b.AddTranslation(langs[0], translations...)
	return nil
}

func parseTranslationFile(filename string) ([]translation.Translation, error) {
	var unmarshalFunc func([]byte, interface{}) error
	switch format := filepath.Ext(filename); format {
	case ".json":
		unmarshalFunc = json.Unmarshal
	/*
		case ".yaml":
			unmarshalFunc = goyaml.Unmarshal
	*/
	default:
		return nil, fmt.Errorf("unsupported file extension %s", format)
	}
	fileBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var translationsData []map[string]interface{}
	if len(fileBytes) > 0 {
		if err := unmarshalFunc(fileBytes, &translationsData); err != nil {
			return nil, err
		}
	}

	translations := make([]translation.Translation, 0, len(translationsData))
	for i, translationData := range translationsData {
		t, err := translation.NewTranslation(translationData)
		if err != nil {
			return nil, fmt.Errorf("unable to parse translation #%d in %s because %s\n%v", i, filename, err, translationData)
		}
		translations = append(translations, t)
	}
	return translations, nil
}

// AddTranslation adds translations for a language.
//
// It is useful if your translations are in a format not supported by LoadTranslationFile.
func (b *Bundle) AddTranslation(lang *language.Language, translations ...translation.Translation) {
	if b.translations[lang.Tag] == nil {
		b.translations[lang.Tag] = make(map[string]translation.Translation, len(translations))
	}
	currentTranslations := b.translations[lang.Tag]
	for _, newTranslation := range translations {
		if currentTranslation := currentTranslations[newTranslation.ID()]; currentTranslation != nil {
			currentTranslations[newTranslation.ID()] = currentTranslation.Merge(newTranslation)
		} else {
			currentTranslations[newTranslation.ID()] = newTranslation
		}
	}
}

// Translations returns all translations in the bundle.
func (b *Bundle) Translations() map[string]map[string]translation.Translation {
	return b.translations
}

// MustTfunc is similar to Tfunc except it panics if an error happens.
func (b *Bundle) MustTfunc(languageSource string, languageSources ...string) TranslateFunc {
	tf, err := b.Tfunc(languageSource, languageSources...)
	if err != nil {
		panic(err)
	}
	return tf
}

// Tfunc returns a TranslateFunc that will be bound to the first language which
// has a non-zero number of translations in the bundle.
//
// It can parse languages from Accept-Language headers (RFC 2616).
func (b *Bundle) Tfunc(src string, srcs ...string) (TranslateFunc, error) {
	lang := b.translatedLanguage(src)
	if lang == nil {
		for _, src := range srcs {
			lang = b.translatedLanguage(src)
			if lang != nil {
				break
			}
		}
	}
	var err error
	if lang == nil {
		err = fmt.Errorf("no supported languages found %#v", append(srcs, src))
	}
	return func(translationID string, args ...interface{}) string {
		return b.translate(lang, translationID, args...)
	}, err
}

func (b *Bundle) translatedLanguage(src string) *language.Language {
	langs := language.Parse(src)
	for _, lang := range langs {
		if len(b.translations[lang.Tag]) > 0 {
			return lang
		}
	}
	return nil
}

func (b *Bundle) translate(lang *language.Language, translationID string, args ...interface{}) string {
	if lang == nil {
		return translationID
	}

	translations := b.translations[lang.Tag]
	if translations == nil {
		return translationID
	}

	translation := translations[translationID]
	if translation == nil {
		return translationID
	}

	var count interface{}
	if len(args) > 0 && isNumber(args[0]) {
		count = args[0]
		args = args[1:]
	}

	plural, _ := lang.Plural(count)
	template := translation.Template(plural)
	if template == nil {
		return translationID
	}

	var data map[string]interface{}
	if len(args) > 0 {
		data, _ = args[0].(map[string]interface{})
	}

	if isNumber(count) {
		if data == nil {
			data = map[string]interface{}{"Count": count}
		} else {
			data["Count"] = count
		}
	}

	s := template.Execute(data)
	if s == "" {
		return translationID
	}
	return s
}

func isNumber(n interface{}) bool {
	switch n.(type) {
	case int, int8, int16, int32, int64, string:
		return true
	}
	return false
}
