package i18n

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	//	"launchpad.net/goyaml"
	"path/filepath"
)

type Bundle struct {
	translations map[string]map[string]Translation
}

type TranslateFunc func(translationId string, args ...interface{}) string

var defaultBundle = NewBundle()

func MustLoadTranslationFile(filename string) {
	defaultBundle.MustLoadTranslationFile(filename)
}

func LoadTranslationFile(filename string) error {
	return defaultBundle.LoadTranslationFile(filename)
}

func Add(locale *Locale, translations ...Translation) {
	defaultBundle.Add(locale, translations...)
}

func MustTfunc(localeId string, localeIds ...string) TranslateFunc {
	return defaultBundle.MustTfunc(localeId, localeIds...)
}

func Tfunc(localeId string, localeIds ...string) (TranslateFunc, error) {
	return defaultBundle.Tfunc(localeId, localeIds...)
}

func NewBundle() *Bundle {
	return &Bundle{
		translations: make(map[string]map[string]Translation),
	}
}

func (b *Bundle) MustLoadTranslationFile(filename string) {
	if err := b.LoadTranslationFile(filename); err != nil {
		panic(err)
	}
}

func (b *Bundle) LoadTranslationFile(filename string) error {
	locale, err := NewLocale(filename)
	if err != nil {
		return err
	}

	translations, err := parseTranslationFile(filename)
	if err != nil {
		return err
	}

	b.Add(locale, translations...)
	return nil
}

func parseTranslationFile(filename string) ([]Translation, error) {
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

	translations := make([]Translation, 0, len(translationsData))
	for i, translationData := range translationsData {
		t, err := NewTranslation(translationData)
		if err != nil {
			return nil, fmt.Errorf("unable to parse translation #%d in %s because %s\n%v", i, filename, err, translationData)
		}
		translations = append(translations, t)
	}
	return translations, nil
}

func (b *Bundle) Add(locale *Locale, translations ...Translation) {
	if b.translations[locale.Id] == nil {
		b.translations[locale.Id] = make(map[string]Translation, len(translations))
	}
	currentTranslations := b.translations[locale.Id]
	for _, newTranslation := range translations {
		if currentTranslation := currentTranslations[newTranslation.Id()]; currentTranslation != nil {
			currentTranslations[newTranslation.Id()] = currentTranslation.Merge(newTranslation)
		} else {
			currentTranslations[newTranslation.Id()] = newTranslation
		}
	}
}

func (b *Bundle) Translations() map[string]map[string]Translation {
	return b.translations
}

func (b *Bundle) MustTfunc(localeId string, localeIds ...string) TranslateFunc {
	tf, err := b.Tfunc(localeId, localeIds...)
	if err != nil {
		panic(err)
	}
	return tf
}

func (b *Bundle) Tfunc(localeId string, localeIds ...string) (tf TranslateFunc, err error) {
	var locale *Locale
	locale, err = NewLocale(localeId)
	if err != nil {
		for _, localeId := range localeIds {
			locale, err = NewLocale(localeId)
			if err == nil {
				break
			}
		}
	}
	return func(translationId string, args ...interface{}) string {
		return b.translate(locale, translationId, args...)
	}, err
}

func (b *Bundle) translate(locale *Locale, translationId string, args ...interface{}) string {
	if locale == nil {
		return translationId
	}

	translations := b.translations[locale.Id]
	if translations == nil {
		return translationId
	}

	translation := translations[translationId]
	if translation == nil {
		return translationId
	}

	var count interface{}
	if len(args) > 0 && isNumber(args[0]) {
		count = args[0]
		args = args[1:]
	}

	pluralCategory, _ := locale.Language.PluralCategory(count)
	template := translation.Template(pluralCategory)
	if template == nil {
		return translationId
	}

	var data map[string]interface{}
	if len(args) > 0 {
		data, _ = args[0].(map[string]interface{})
	}

	if isNumber(count) {
		if data == nil {
			data = map[string]interface{}{"Count": count}
		} else if _, ok := data["Count"]; !ok {
			data["Count"] = count
		}
	}

	s := template.Execute(data)
	if s == "" {
		return translationId
	}
	return s
}

func isNumber(n interface{}) bool {
	switch n.(type) {
	case int, int8, int16, int32, int64, float32, float64:
		return true
	}
	return false
}
