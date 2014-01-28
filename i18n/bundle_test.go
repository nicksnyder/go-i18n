package i18n

import (
	"testing"
)

func TestMustLoadTranslationFile(t *testing.T) {
	t.Skipf("not implemented")
}

func TestLoadTranslationFile(t *testing.T) {
	t.Skipf("not implemented")
}

func TestParseTranslationFile(t *testing.T) {
	t.Skipf("not implemented")
}

func TestAdd(t *testing.T) {
	t.Skipf("not implemented")
}

/*

func bundleFixture(t *testing.T) *Bundle {
	l, err := NewLocaleFromString("ar-EG")
	if err != nil {
		t.Errorf(err.Error())
	}
	return &Bundle{
		Locale: l,
		localizedStrings: map[string]*LocalizedString{
			"a": &LocalizedString{
				Id: "a",
			},
			"b": &LocalizedString{
				Id:          "b",
				Translation: "translation(b)",
			},
			"c": &LocalizedString{
				Id: "c",
				Translations: map[PluralCategory]*PluralTranslation{
					Zero:  NewPluralTranslation("zero(c)"),
					One:   NewPluralTranslation("one(c)"),
					Two:   NewPluralTranslation("two(c)"),
					Few:   NewPluralTranslation("few(c)"),
					Many:  NewPluralTranslation("many(c)"),
					Other: NewPluralTranslation("other(c)"),
				},
			},
			"d": &LocalizedString{
				Id: "d",
				Translations: map[PluralCategory]*PluralTranslation{
					Zero: NewPluralTranslation("zero(d)"),
					One:  NewPluralTranslation("one(d)"),
				},
			},
		},
	}
}
*/
