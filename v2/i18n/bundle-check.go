package i18n

import (
	"fmt"
	"strings"

	"golang.org/x/text/language"
)

type bundleIssues []string

func (bi bundleIssues) Error() string {
	return strings.Join(bi, ", ")
}

func (b *Bundle) Check() error {
	var issues bundleIssues
	ids := make(map[language.Tag]map[string]bool)
	tags := b.LanguageTags()
	if len(tags) < 2 {
		return nil
	}
	// map all keys
	for _, t := range tags {
		ids[t] = make(map[string]bool)
		for k := range b.messageTemplates[t] {
			ids[t][k] = true
		}
	}
	// what's missing in non-default bundles?
	defaultTag := tags[0]
	for i := 1; i < len(tags); i++ {
		for k := range ids[defaultTag] {
			if !ids[tags[i]][k] {
				issues = append(issues, fmt.Sprintf("%v doesn't have key %q", tags[i], k))
			}
		}
	}
	// what's extra in non-default bundles?
	for i := 1; i < len(tags); i++ {
		for k := range ids[tags[i]] {
			if !ids[defaultTag][k] {
				issues = append(issues, fmt.Sprintf("%v has extra key %q", tags[i], k))
			}
		}
	}
	if len(issues) > 0 {
		return &issues
	}
	return nil
}
