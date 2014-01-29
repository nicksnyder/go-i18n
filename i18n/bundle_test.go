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

func TestAdd(t *testing.T) {
	t.Skipf("not implemented")
}

func TestMustTfunc(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected MustTfunc to panic")
		}
	}()
	NewBundle().MustTfunc("invalid")
}

func TestTfunc(t *testing.T) {
	b := NewBundle()
	translationId := "translation_id"
	englishTranslation := "en-US(translation_id)"
	b.Add(mustNewLocale("en-US"), testNewTranslation(t, map[string]interface{}{
		"id":          translationId,
		"translation": englishTranslation,
	}))
	frenchTranslation := "fr-FR(translation_id)"
	b.Add(mustNewLocale("fr-FR"), testNewTranslation(t, map[string]interface{}{
		"id":          translationId,
		"translation": frenchTranslation,
	}))

	tests := []struct {
		localeIds []string
		valid     bool
		result    string
	}{
		{
			[]string{"invalid"},
			false,
			translationId,
		},
		{
			[]string{"invalid", "invalid2"},
			false,
			translationId,
		},
		{
			[]string{"invalid", "en-US"},
			true,
			englishTranslation,
		},
		{
			[]string{"en-US", "invalid"},
			true,
			englishTranslation,
		},
		{
			[]string{"en-US", "fr-FR"},
			true,
			englishTranslation,
		},
	}

	for _, test := range tests {
		tf, err := b.Tfunc(test.localeIds[0], test.localeIds[1:]...)
		if err != nil && test.valid {
			t.Errorf("Tfunc for %v returned error %s", test.localeIds, err)
		}
		if err == nil && !test.valid {
			t.Errorf("Tfunc for %v returned nil error", test.localeIds)
		}
		if result := tf(translationId); result != test.result {
			t.Errorf("translation was %s; expected %s", result, test.result)
		}
	}
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
