package i18n

import (
	"fmt"
	"regexp"
)

type Locale struct {
	Id       string
	Language *Language
}

// Matches strings like aa-CC, and aa-Bbbb-CC.
var languageTagRegexp = regexp.MustCompile(`([a-z]{2,3}(?:[_\-][A-Z][a-z]{3})?)[_\-][A-Z]{2}`)

// NewLocale searches s for a valid language tag as defined by RFC 5646.
// http://en.wikipedia.org/wiki/IETF_language_tag
// It returns an error if s doesn't contain exactly one language tag or
// if the language represented by the tag is not supported by this package.
func NewLocale(s string) (*Locale, error) {
	matches := languageTagRegexp.FindAllStringSubmatch(s, -1)
	if count := len(matches); count != 1 {
		return nil, fmt.Errorf("%d locales found in string %s", count, s)
	}
	id, languageCode := matches[0][0], matches[0][1]
	language := LanguageWithCode(languageCode)
	if language == nil {
		return nil, fmt.Errorf("unknown language code %s", languageCode)
	}
	return &Locale{id, language}, nil
}

func mustNewLocale(s string) *Locale {
	locale, err := NewLocale(s)
	if err != nil {
		panic(err)
	}
	return locale
}
