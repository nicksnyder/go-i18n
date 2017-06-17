// Package language defines languages that implement CLDR pluralization.
package language

import (
	"fmt"
	"strings"
)

// Language is a written human language.
type Language struct {
	matchingTags []string

	// Tag uniquely identifies the language as defined by RFC 5646.
	//
	// Most language tags are a two character language code (ISO 639-1)
	// optionally followed by a dash and a two character country code (ISO 3166-1).
	// (e.g. en, pt-br)
	Tag string
	*PluralSpec
}

// NewLanguage returns new Language with specified tag and spec.
func NewLanguage(tag string, spec *PluralSpec) *Language {
	return &Language{Tag: tag, PluralSpec: spec}
}

func (l *Language) String() string {
	return l.Tag
}

// MatchingTags returns the set of language tags that map to this Language.
// e.g. "zh-hans-cn" yields {"zh", "zh-hans", "zh-hans-cn"}.
func (l *Language) MatchingTags() []string {
	if len(l.matchingTags) > 0 {
		return l.matchingTags
	}

	parts := strings.Split(l.Tag, "-")
	if len(parts) == 1 {
		l.matchingTags = parts
		return parts
	}

	matches := make([]string, 0, len(parts))
	matches = append(matches, parts[0])
	for i, part := range parts[1:] {
		matches = append(matches, matches[i]+"-"+part)
	}
	l.matchingTags = matches
	return l.matchingTags
}

// Parse returns a slice of supported languages found in src or nil if none are found.
// It can parse language tags and Accept-Language headers.
func Parse(src string) []*Language {
	var langs []*Language
	start := 0
	for end, chr := range src {
		switch chr {
		case ',', ';', '.':
			tag := NormalizeTag(strings.TrimSpace(src[start:end]))
			if spec := getPluralSpec(tag); spec != nil {
				langs = append(langs, NewLanguage(tag, spec))
			}
			start = end + 1
		}
	}
	if start > 0 {
		tag := NormalizeTag(strings.TrimSpace(src[start:]))
		if spec := getPluralSpec(tag); spec != nil {
			langs = append(langs, NewLanguage(tag, spec))
		}
		return dedupe(langs)
	}
	src = NormalizeTag(src)
	if spec := getPluralSpec(src); spec != nil {
		langs = append(langs, NewLanguage(src, spec))
	}
	return langs
}

func dedupe(langs []*Language) []*Language {
	found := make(map[string]struct{}, len(langs))
	deduped := make([]*Language, 0, len(langs))
	for _, lang := range langs {
		if _, ok := found[lang.Tag]; !ok {
			found[lang.Tag] = struct{}{}
			deduped = append(deduped, lang)
		}
	}
	return deduped
}

// MustParse is similar to Parse except it panics instead of retuning a nil Language.
func MustParse(src string) []*Language {
	langs := Parse(src)
	if len(langs) == 0 {
		panic(fmt.Errorf("unable to parse language from %q", src))
	}
	return langs
}

// Add adds support for a new language.
func Add(l *Language) {
	tag := NormalizeTag(l.Tag)
	pluralSpecs[tag] = l.PluralSpec
}

// NormalizeTag returns a language tag with all lower-case characters
// and dashes "-" instead of underscores "_"
func NormalizeTag(tag string) string {
	tag = strings.ToLower(tag)
	return strings.Replace(tag, "_", "-", -1)
}
