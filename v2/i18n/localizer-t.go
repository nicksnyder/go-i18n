package i18n

import (
	"html/template"
	"regexp"

	"github.com/zengabor/go-i18n/v2/internal/plural"
	"golang.org/x/text/language"
)

func (l *Localizer) Tags() []language.Tag {
	return l.tags
}

func (l *Localizer) Language() string {
	if len(l.tags) > 0 {
		return l.tags[0].String()
	}
	return l.bundle.defaultLanguage.String()
}

func (l *Localizer) DefaultLanguageTag() language.Tag {
	return l.bundle.defaultLanguage
}

func (l *Localizer) Translate(messageID string, templateData map[string]interface{}) (string, error) {
	return l.Localize(&LocalizeConfig{MessageID: messageID, TemplateData: templateData})
}

// To be used in template FuncMap, like tmpl.Funcs(template.FuncMap{"T": localizer.T()})
func (l *Localizer) T() func(messageID string, args ...interface{}) (template.HTML, error) {
	return func(messageID string, args ...interface{}) (template.HTML, error) {
		lc := LocalizeConfig{MessageID: messageID}
		td := make(map[string]interface{})
		al := len(args)
		for i, n := range l.referenceArgNames(messageID) {
			if i >= al {
				break
			}
			td[n] = args[i]
			if n == "PluralCount" {
				lc.PluralCount = args[i]
			}
		}
		lc.TemplateData = td
		s, err := l.Localize(&lc)
		return template.HTML(s), err
	}
}

var argMatcher = regexp.MustCompile(`{{\s?\.(\w+)\s?}}`)

// Collects the names of arguments in the "Other" template of the default language of the bundle
func (l *Localizer) referenceArgNames(id string) []string {
	tags := []language.Tag{l.bundle.defaultLanguage}
	_, t := l.matchTemplate(id, nil, language.NewMatcher(tags), tags)
	if t == nil {
		return nil
	}
	var r []string
	if o, ok := t.PluralTemplates[plural.Other]; ok {
		for _, m := range argMatcher.FindAllStringSubmatch(o.Src, -1) {
			r = append(r, m[1])
		}
	}
	return r
}
