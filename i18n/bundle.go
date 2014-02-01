package i18n

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	//	"launchpad.net/goyaml"
	"path/filepath"
)

// TranslateFunc returns the translation of the string identified by translationID.
//
// If translationID is a non-plural form, then the first variadic argument may be a map[string]interface{}
// that contains template data for the string (if any).
//
// If translationID is a plural form, then the first variadic argument must be a number type
// (int, int8, int16, int32, int64, float32, float64) and the second variadic argument may be a
// map[string]interface like the non-plural form.
type TranslateFunc func(translationID string, args ...interface{}) string

// IdentityTfunc returns a TranslateFunc that always returns the translationID passed to it.
func IdentityTfunc() TranslateFunc {
	return func(translationID string, args ...interface{}) string {
		return translationID
	}
}

type bundle struct {
	translations map[string]map[string]Translation
}

var defaultBundle = newBundle()

// MustLoadTranslationFile is similar to LoadTranslationFile
// except it panics if an error happens.
func MustLoadTranslationFile(filename string) {
	defaultBundle.MustLoadTranslationFile(filename)
}

// LoadTranslationFile loads the translations from filename into memory.
//
// The locale that the translations are associated with is parsed from the filename.
//
// Generally you should load translation files once during your program's initialization.
func LoadTranslationFile(filename string) error {
	return defaultBundle.LoadTranslationFile(filename)
}

// Add adds translations for a locale.
//
// Add is useful if your translations are in a format not supported by LoadTranslationFile.
func AddTranslation(locale *Locale, translations ...Translation) {
	defaultBundle.AddTranslation(locale, translations...)
}

// MustTfunc is similar to Tfunc except it panics if an error happens.
func MustTfunc(localeID string, localeIDs ...string) TranslateFunc {
	return defaultBundle.MustTfunc(localeID, localeIDs...)
}

// Tfunc returns a TranslateFunc that will be bound to the first valid locale from its parameters.
func Tfunc(localeID string, localeIDs ...string) (TranslateFunc, error) {
	return defaultBundle.Tfunc(localeID, localeIDs...)
}

func newBundle() *bundle {
	return &bundle{
		translations: make(map[string]map[string]Translation),
	}
}

func (b *bundle) MustLoadTranslationFile(filename string) {
	if err := b.LoadTranslationFile(filename); err != nil {
		panic(err)
	}
}

func (b *bundle) LoadTranslationFile(filename string) error {
	locale, err := NewLocale(filename)
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

func (b *bundle) AddTranslation(locale *Locale, translations ...Translation) {
	if b.translations[locale.ID] == nil {
		b.translations[locale.ID] = make(map[string]Translation, len(translations))
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

func (b *bundle) Translations() map[string]map[string]Translation {
	return b.translations
}

func (b *bundle) MustTfunc(localeID string, localeIDs ...string) TranslateFunc {
	tf, err := b.Tfunc(localeID, localeIDs...)
	if err != nil {
		panic(err)
	}
	return tf
}

func (b *bundle) Tfunc(localeID string, localeIDs ...string) (tf TranslateFunc, err error) {
	var locale *Locale
	locale, err = NewLocale(localeID)
	if err != nil {
		for _, localeID := range localeIDs {
			locale, err = NewLocale(localeID)
			if err == nil {
				break
			}
		}
	}
	return func(translationID string, args ...interface{}) string {
		return b.translate(locale, translationID, args...)
	}, err
}

func (b *bundle) translate(locale *Locale, translationID string, args ...interface{}) string {
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

	pluralCategory, _ := locale.Language.pluralCategory(count)
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
	case int, int8, int16, int32, int64, float32, float64:
		return true
	}
	return false
}
