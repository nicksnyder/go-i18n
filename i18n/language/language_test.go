package language

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		src  string
		lang []*Language
	}{
		{"en", []*Language{NewLanguage("en", pluralSpecs["en"])}},
		{"en-US", []*Language{NewLanguage("en-us", pluralSpecs["en"])}},
		{"en_US", []*Language{NewLanguage("en-us", pluralSpecs["en"])}},
		{"en-GB", []*Language{NewLanguage("en-gb", pluralSpecs["en"])}},
		{"zh-CN", []*Language{NewLanguage("zh-cn", pluralSpecs["zh"])}},
		{"zh-TW", []*Language{NewLanguage("zh-tw", pluralSpecs["zh"])}},
		{"pt-BR", []*Language{NewLanguage("pt-br", pluralSpecs["pt"])}},
		{"pt_BR", []*Language{NewLanguage("pt-br", pluralSpecs["pt"])}},
		{"pt-PT", []*Language{NewLanguage("pt-pt", pluralSpecs["pt-pt"])}},
		{"pt_PT", []*Language{NewLanguage("pt-pt", pluralSpecs["pt-pt"])}},
		{"zh-Hans-CN", []*Language{NewLanguage("zh-hans-cn", pluralSpecs["zh"])}},
		{"zh-Hant-TW", []*Language{NewLanguage("zh-hant-tw", pluralSpecs["zh"])}},
		{"en-US-en-US", []*Language{NewLanguage("en-us-en-us", pluralSpecs["en"])}},
		{".en-US..en-US.", []*Language{NewLanguage("en-us", pluralSpecs["en"])}},
		{
			"it, xx-zz, xx-ZZ, zh, en-gb;q=0.8, en;q=0.7, es-ES;q=0.6, de-xx",
			[]*Language{
				NewLanguage("it", pluralSpecs["it"]),
				NewLanguage("zh", pluralSpecs["zh"]),
				NewLanguage("en-gb", pluralSpecs["en"]),
				NewLanguage("en", pluralSpecs["en"]),
				NewLanguage("es-es", pluralSpecs["es"]),
				NewLanguage("de-xx", pluralSpecs["de"]),
			},
		},
		{
			"it-qq,xx,xx-zz,xx-ZZ,zh,en-gb;q=0.8,en;q=0.7,es-ES;q=0.6,de-xx",
			[]*Language{
				NewLanguage("it-qq", pluralSpecs["it"]),
				NewLanguage("zh", pluralSpecs["zh"]),
				NewLanguage("en-gb", pluralSpecs["en"]),
				NewLanguage("en", pluralSpecs["en"]),
				NewLanguage("es-es", pluralSpecs["es"]),
				NewLanguage("de-xx", pluralSpecs["de"]),
			},
		},
		{"en.json", []*Language{NewLanguage("en", pluralSpecs["en"])}},
		{"en-US.json", []*Language{NewLanguage("en-us", pluralSpecs["en"])}},
		{"en-us.json", []*Language{NewLanguage("en-us", pluralSpecs["en"])}},
		{"en-xx.json", []*Language{NewLanguage("en-xx", pluralSpecs["en"])}},
		{"xx-Yyen-US", nil},
		{"en US", nil},
		{"", nil},
		{"-", nil},
		{"_", nil},
		{"-en", nil},
		{"_en", nil},
		{"-en-", nil},
		{"_en_", nil},
		{"xx", nil},
	}
	for _, test := range tests {
		lang := Parse(test.src)
		if !reflect.DeepEqual(lang, test.lang) {
			t.Errorf("Parse(%q) = %s expected %s", test.src, lang, test.lang)
		}
	}
}

func TestMatchingTags(t *testing.T) {
	tests := []struct {
		lang    *Language
		matches []string
	}{
		{NewLanguage("zh-hans-cn", nil), []string{"zh", "zh-hans", "zh-hans-cn"}},
		{NewLanguage("foo", nil), []string{"foo"}},
	}
	for _, test := range tests {
		if actual := test.lang.MatchingTags(); !reflect.DeepEqual(test.matches, actual) {
			t.Errorf("matchingTags(%q) = %q expected %q", test.lang.Tag, actual, test.matches)
		}
	}
}
