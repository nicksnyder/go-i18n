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
	tags := AllLanguageTags(src)
	if len(tags) == 0 {
		return nil
	}
	langs := make([]*Language, 0, len(tags))

	for _, lang := range tags {
		if spec := getPluralSpec(lang); spec != nil {
			langs = append(langs, &Language{lang, spec})
		}
	}
	return langs
}

// AllLanguageTags returns all normalized language tags found in src, including unsupported.
func AllLanguageTags(src string) []string {
	if src == "" {
		return nil
	}

	var tags []string
	found := make(map[string]bool)

	src = strings.TrimLeft(src, ",;.")
	start := 0
	for end, chr := range src {
		switch chr {
		case ',', ';', '.':
			tag := NormalizeTag(strings.TrimSpace(src[start:end]))
			if tag != "" && !isExtension(tag) && !found[tag] {
				found[tag] = true
				tags = append(tags, tag)
			}
			start = end + 1
		}
	}
	if start > 0 {
		tag := NormalizeTag(strings.TrimSpace(src[start:]))
		if tag != "" && !isExtension(tag) && !found[tag] {
			tags = append(tags, tag)
		}
		return tags
	}
	tags = append(tags, NormalizeTag(src))
	return tags
}

func isExtension(s string) bool {
	return s == "json" || s == "yaml" || s == "yml" || s == "toml"
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
