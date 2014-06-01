// Package bundle manages translations for multiple locales.
package bundle

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	//	"launchpad.net/goyaml"
	"github.com/nicksnyder/go-i18n/i18n/locale"
	"github.com/nicksnyder/go-i18n/i18n/translation"
	"path/filepath"
)

// TranslateFunc is a copy of i18n.TranslateFunc to avoid a circular dependency.
type TranslateFunc func(translationID string, args ...interface{}) string

type Bundle struct {
	translations map[string]map[string]translation.Translation
}

func New() *Bundle {
	return &Bundle{
		translations: make(map[string]map[string]translation.Translation),
	}
}

func (b *Bundle) MustLoadTranslationFile(filename string) {
	if err := b.LoadTranslationFile(filename); err != nil {
		panic(err)
	}
}

func (b *Bundle) LoadTranslationFile(filename string) error {
	locale, err := locale.New(filename)
	if err != nil {
		return err
	}

	translations, err := parseTranslationFile(filename)
	if err != nil {
		return err
	}

	b.AddTranslation(locale, translations...)
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

func (b *Bundle) AddTranslation(locale *locale.Locale, translations ...translation.Translation) {
	if b.translations[locale.ID] == nil {
		b.translations[locale.ID] = make(map[string]translation.Translation, len(translations))
	}
	currentTranslations := b.translations[locale.ID]
	for _, newTranslation := range translations {
		if currentTranslation := currentTranslations[newTranslation.ID()]; currentTranslation != nil {
			currentTranslations[newTranslation.ID()] = currentTranslation.Merge(newTranslation)
		} else {
			currentTranslations[newTranslation.ID()] = newTranslation
		}
	}
}

func (b *Bundle) Translations() map[string]map[string]translation.Translation {
	return b.translations
}

func (b *Bundle) MustTfunc(localeID string, localeIDs ...string) TranslateFunc {
	tf, err := b.Tfunc(localeID, localeIDs...)
	if err != nil {
		panic(err)
	}
	return tf
}

func (b *Bundle) Tfunc(localeID string, localeIDs ...string) (tf TranslateFunc, err error) {
	var l *locale.Locale
	l, err = locale.New(localeID)
	if err != nil {
		for _, localeID := range localeIDs {
			l, err = locale.New(localeID)
			if err == nil {
				break
			}
		}
	}
	return func(translationID string, args ...interface{}) string {
		return b.translate(l, translationID, args...)
	}, err
}

func (b *Bundle) translate(locale *locale.Locale, translationID string, args ...interface{}) string {
	if locale == nil {
		return translationID
	}

	translations := b.translations[locale.ID]
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

	pluralCategory, _ := locale.Language.PluralCategory(count)
	template := translation.Template(pluralCategory)
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
