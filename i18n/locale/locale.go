// Package locale parses locale strings.
package locale

import (
	"fmt"
	"github.com/nicksnyder/go-i18n/i18n/language"
	"regexp"
	"strings"
)

// Locale is a language and a geographic region (e.g. en-US, en-GB).
type Locale struct {
	ID       string
	Language *language.Language
}

// tagMatcher matches language tags (e.g. zh-CN).
var tagMatcher = regexp.MustCompile(`^([a-z]{2})[_\-]([A-Z]{2})$`)

// tagSplitter matches characters not found in language tags.
var tagSplitter = regexp.MustCompile(`[^a-zA-Z_\-]+`)

// New searches s for a valid language tag (RFC 5646)
// of the form xx-YY or xx_YY where
// xx is a 2 character language code and
// YY is a 2 character country code.
//
// It returns an error if s doesn't contain exactly one language tag or
// if the language represented by the tag is not supported by this package.
func New(s string) (*Locale, error) {
	parts := tagSplitter.Split(s, -1)
	var id, lc string
	count := 0
	for _, part := range parts {
		if tag := tagMatcher.FindStringSubmatch(part); tag != nil {
			count += 1
			id, lc = tag[0], tag[1]
		}
	}
	if count != 1 {
		return nil, fmt.Errorf("%d locales found in string %s", count, s)
	}
	id = strings.Replace(id, "_", "-", -1)
	lang := language.LanguageWithID(id)
	if lang == nil {
		lang = language.LanguageWithID(lc)
	}
	if lang == nil {
		return nil, fmt.Errorf("unknown language %s", id)
	}
	return &Locale{id, lang}, nil
}

// MustNew is similar to New except that it panics if an error happens.
func MustNew(s string) *Locale {
	locale, err := New(s)
	if err != nil {
		panic(err)
	}
	return locale
}

func (l *Locale) String() string {
	return l.ID
}
