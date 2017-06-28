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
		{"en", []*Language{&Language{Tag: "en", PluralSpec: pluralSpecs["en"]}}},
		{"en-US", []*Language{&Language{Tag: "en-us", PluralSpec: pluralSpecs["en"]}}},
		{"en_US", []*Language{&Language{Tag: "en-us", PluralSpec: pluralSpecs["en"]}}},
		{"en-GB", []*Language{&Language{Tag: "en-gb", PluralSpec: pluralSpecs["en"]}}},
		{"zh-CN", []*Language{&Language{Tag: "zh-cn", PluralSpec: pluralSpecs["zh"]}}},
		{"zh-TW", []*Language{&Language{Tag: "zh-tw", PluralSpec: pluralSpecs["zh"]}}},
		{"pt-BR", []*Language{&Language{Tag: "pt-br", PluralSpec: pluralSpecs["pt"]}}},
		{"pt_BR", []*Language{&Language{Tag: "pt-br", PluralSpec: pluralSpecs["pt"]}}},
		{"pt-PT", []*Language{&Language{Tag: "pt-pt", PluralSpec: pluralSpecs["pt-pt"]}}},
		{"pt_PT", []*Language{&Language{Tag: "pt-pt", PluralSpec: pluralSpecs["pt-pt"]}}},
		{"zh-Hans-CN", []*Language{&Language{Tag: "zh-hans-cn", PluralSpec: pluralSpecs["zh"]}}},
		{"zh-Hant-TW", []*Language{&Language{Tag: "zh-hant-tw", PluralSpec: pluralSpecs["zh"]}}},
		{"en-US-en-US", []*Language{&Language{Tag: "en-us-en-us", PluralSpec: pluralSpecs["en"]}}},
		{".en-US..en-US.", []*Language{&Language{Tag: "en-us", PluralSpec: pluralSpecs["en"]}}},
		{
			"it, xx-zz, xx-ZZ, zh, en-gb;q=0.8, en;q=0.7, es-ES;q=0.6, de-xx",
			[]*Language{
				&Language{Tag: "it", PluralSpec: pluralSpecs["it"]},
				&Language{Tag: "zh", PluralSpec: pluralSpecs["zh"]},
				&Language{Tag: "en-gb", PluralSpec: pluralSpecs["en"]},
				&Language{Tag: "en", PluralSpec: pluralSpecs["en"]},
				&Language{Tag: "es-es", PluralSpec: pluralSpecs["es"]},
				&Language{Tag: "de-xx", PluralSpec: pluralSpecs["de"]},
			},
		},
		{
			"it-qq,xx,xx-zz,xx-ZZ,zh,en-gb;q=0.8,en;q=0.7,es-ES;q=0.6,de-xx",
			[]*Language{
				&Language{Tag: "it-qq", PluralSpec: pluralSpecs["it"]},
				&Language{Tag: "zh", PluralSpec: pluralSpecs["zh"]},
				&Language{Tag: "en-gb", PluralSpec: pluralSpecs["en"]},
				&Language{Tag: "en", PluralSpec: pluralSpecs["en"]},
				&Language{Tag: "es-es", PluralSpec: pluralSpecs["es"]},
				&Language{Tag: "de-xx", PluralSpec: pluralSpecs["de"]},
			},
		},
		{"en.json", []*Language{&Language{Tag: "en", PluralSpec: pluralSpecs["en"]}}},
		{"en-US.json", []*Language{&Language{Tag: "en-us", PluralSpec: pluralSpecs["en"]}}},
		{"en-us.json", []*Language{&Language{Tag: "en-us", PluralSpec: pluralSpecs["en"]}}},
		{"en-xx.json", []*Language{&Language{Tag: "en-xx", PluralSpec: pluralSpecs["en"]}}},
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
		{&Language{Tag: "zh-hans-cn", PluralSpec: nil}, []string{"zh", "zh-hans", "zh-hans-cn"}},
		{&Language{Tag: "foo", PluralSpec: nil}, []string{"foo"}},
	}
	for _, test := range tests {
		if actual := test.lang.MatchingTags(); !reflect.DeepEqual(test.matches, actual) {
			t.Errorf("matchingTags(%q) = %q expected %q", test.lang.Tag, actual, test.matches)
		}
	}
}
