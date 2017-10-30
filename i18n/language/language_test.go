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
		{"en", []*Language{{"en", pluralSpecsStore.get("en")}}},
		{"en-US", []*Language{{"en-us", pluralSpecsStore.get("en")}}},
		{"en_US", []*Language{{"en-us", pluralSpecsStore.get("en")}}},
		{"en-GB", []*Language{{"en-gb", pluralSpecsStore.get("en")}}},
		{"zh-CN", []*Language{{"zh-cn", pluralSpecsStore.get("zh")}}},
		{"zh-TW", []*Language{{"zh-tw", pluralSpecsStore.get("zh")}}},
		{"pt-BR", []*Language{{"pt-br", pluralSpecsStore.get("pt")}}},
		{"pt_BR", []*Language{{"pt-br", pluralSpecsStore.get("pt")}}},
		{"pt-PT", []*Language{{"pt-pt", pluralSpecsStore.get("pt")}}},
		{"pt_PT", []*Language{{"pt-pt", pluralSpecsStore.get("pt")}}},
		{"zh-Hans-CN", []*Language{{"zh-hans-cn", pluralSpecsStore.get("zh")}}},
		{"zh-Hant-TW", []*Language{{"zh-hant-tw", pluralSpecsStore.get("zh")}}},
		{"en-US-en-US", []*Language{{"en-us-en-us", pluralSpecsStore.get("en")}}},
		{".en-US..en-US.", []*Language{{"en-us", pluralSpecsStore.get("en")}}},
		{
			"it, xx-zz, xx-ZZ, zh, en-gb;q=0.8, en;q=0.7, es-ES;q=0.6, de-xx",
			[]*Language{
				{"it", pluralSpecsStore.get("it")},
				{"zh", pluralSpecsStore.get("zh")},
				{"en-gb", pluralSpecsStore.get("en")},
				{"en", pluralSpecsStore.get("en")},
				{"es-es", pluralSpecsStore.get("es")},
				{"de-xx", pluralSpecsStore.get("de")},
			},
		},
		{
			"it-qq,xx,xx-zz,xx-ZZ,zh,en-gb;q=0.8,en;q=0.7,es-ES;q=0.6,de-xx",
			[]*Language{
				{"it-qq", pluralSpecsStore.get("it")},
				{"zh", pluralSpecsStore.get("zh")},
				{"en-gb", pluralSpecsStore.get("en")},
				{"en", pluralSpecsStore.get("en")},
				{"es-es", pluralSpecsStore.get("es")},
				{"de-xx", pluralSpecsStore.get("de")},
			},
		},
		{"en.json", []*Language{{"en", pluralSpecsStore.get("en")}}},
		{"en-US.json", []*Language{{"en-us", pluralSpecsStore.get("en")}}},
		{"en-us.json", []*Language{{"en-us", pluralSpecsStore.get("en")}}},
		{"en-xx.json", []*Language{{"en-xx", pluralSpecsStore.get("en")}}},
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
		{&Language{"zh-hans-cn", nil}, []string{"zh", "zh-hans", "zh-hans-cn"}},
		{&Language{"foo", nil}, []string{"foo"}},
	}
	for _, test := range tests {
		if actual := test.lang.MatchingTags(); !reflect.DeepEqual(test.matches, actual) {
			t.Errorf("matchingTags(%q) = %q expected %q", test.lang.Tag, actual, test.matches)
		}
	}
}
