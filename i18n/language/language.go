// Package language defines languages that implement CLDR pluralization.
package language

import (
	"fmt"
	"strings"
)

// Language is a written human language.
type Language struct {
	// Tag uniquely identifies the language as defined by RFC 5646.
	//
	// Most language tags are a two character language code (ISO 639-1)
	// optionally followed by a dash and a two character country code (ISO 3166-1).
	// (e.g. en, pt-br)
	Tag string
	*PluralSpec
}

// NewLanguage returns Language with specified tag and pluralSpec.
// If pluralSpec is nil, it uses EmptyPluralSpec.
func NewLanguage(tag string, pluralSpec *PluralSpec) *Language {
	if pluralSpec == nil {
		return &Language{Tag: tag, PluralSpec: EmptyPluralSpec}
	}
	return &Language{Tag: tag, PluralSpec: pluralSpec}
}

func (l *Language) String() string {
	return l.Tag
}

// MatchingTags returns the set of language tags that map to this Language.
// e.g. "zh-hans-cn" yields {"zh", "zh-hans", "zh-hans-cn"}
// BUG: This should be computed once and stored as a field on Language for efficiency,
//      but this would require changing how Languages are constructed.
func (l *Language) MatchingTags() []string {
	parts := strings.Split(l.Tag, "-")
	var prefix, matches []string
	for _, part := range parts {
		prefix = append(prefix, part)
		match := strings.Join(prefix, "-")
		matches = append(matches, match)
	}
	return matches
}

// Parse returns a slice of supported languages found in src or nil if none are found.
// It can parse language tags and Accept-Language headers.
func Parse(src string) []*Language {
	if src == "" {
		return nil
	}

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
		return uniq(langs)
	}
	src = NormalizeTag(src)
	if spec := getPluralSpec(src); spec != nil {
		langs = append(langs, NewLanguage(src, spec))
	}
	return langs
}

func uniq(langs []*Language) []*Language {
	found := make(map[string]struct{}, len(langs))
	unique := make([]*Language, 0, len(langs))
	for _, lang := range langs {
		if _, ok := found[lang.Tag]; !ok {
			found[lang.Tag] = struct{}{}
			unique = append(unique, lang)
		}
	}
	return unique
}

// ParseFirst returns Language with first found language code in src.
// It also supports custom, non-provided in Unicode CLDR language codes.
// It can parse language tags and Accept-Language headers.
func ParseFirst(src string) *Language {
	if src == "" {
		return nil
	}

	delims := ",;."
	src = strings.TrimLeft(src, delims)

	delimIndex := strings.IndexAny(src, delims)
	if delimIndex == -1 {
		delimIndex = len(src)
	}
	tag := NormalizeTag(strings.TrimSpace(src[:delimIndex]))
	return NewLanguage(tag, getPluralSpec(tag))
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
