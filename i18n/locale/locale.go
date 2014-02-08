// Package locale parses locale strings.
package locale

import (
	"fmt"
	"github.com/nicksnyder/go-i18n/i18n/language"
	"regexp"
)

// Locale is a language and a geographic region (e.g. en-US, en-GB).
type Locale struct {
	ID       string
	Language *language.Language
}

// Matches strings like aa-CC, and aa-Bbbb-CC.
var languageTagRegexp = regexp.MustCompile(`([a-z]{2,3}(?:[_\-][A-Z][a-z]{3})?)[_\-][A-Z]{2}`)

// New searches s for a valid language tag as defined by RFC 5646.
//
// It returns an error if s doesn't contain exactly one language tag or
// if the language represented by the tag is not supported by this package.
func New(s string) (*Locale, error) {
	matches := languageTagRegexp.FindAllStringSubmatch(s, -1)
	if count := len(matches); count != 1 {
		return nil, fmt.Errorf("%d locales found in string %s", count, s)
	}
	id, languageCode := matches[0][0], matches[0][1]
	language := language.LanguageWithID(languageCode)
	if language == nil {
		return nil, fmt.Errorf("unknown language code %s", languageCode)
	}
	return &Locale{id, language}, nil
}

// MustNew is similar to New except that it panics if an error happens.
func MustNew(s string) *Locale {
	locale, err := New(s)
	if err != nil {
		panic(err)
	}
	return locale
}
