package i18n

import (
	"fmt"
	"regexp"
)

type Locale struct {
	Id       string
	Language *Language
}

var localeIdRegexp = regexp.MustCompile(`([a-z]{2,3}(?:[_\-][A-Z][a-z]{3})?)[_\-][A-Z]{2}`)

func NewLocale(s string) (*Locale, error) {
	matches := localeIdRegexp.FindAllStringSubmatch(s, -1)
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

func mustNewLocale(localeId string) *Locale {
	locale, err := NewLocale(localeId)
	if err != nil {
		panic(err)
	}
	return locale
}
