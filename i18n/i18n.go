package i18n

import (
	"github.com/nicksnyder/go-i18n/i18n/bundle"
	"github.com/nicksnyder/go-i18n/i18n/locale"
	"github.com/nicksnyder/go-i18n/i18n/translation"
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
//
// It is a useful placeholder when parsing a text/template or html/template
// before the actual Tfunc is available.
func IdentityTfunc() TranslateFunc {
	return func(translationID string, args ...interface{}) string {
		return translationID
	}
}

var defaultBundle = bundle.New()

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
func AddTranslation(locale *locale.Locale, translations ...translation.Translation) {
	defaultBundle.AddTranslation(locale, translations...)
}

// MustTfunc is similar to Tfunc except it panics if an error happens.
func MustTfunc(localeID string, localeIDs ...string) TranslateFunc {
	return TranslateFunc(defaultBundle.MustTfunc(localeID, localeIDs...))
}

// Tfunc returns a TranslateFunc that will be bound to the first valid locale from its parameters.
func Tfunc(localeID string, localeIDs ...string) (TranslateFunc, error) {
	tf, err := defaultBundle.Tfunc(localeID, localeIDs...)
	return TranslateFunc(tf), err
}
