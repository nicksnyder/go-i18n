package bundle

import (
	"fmt"
	"testing"

	"github.com/nicksnyder/go-i18n/i18n/language"
	"github.com/nicksnyder/go-i18n/i18n/translation"
)

func TestMustLoadTranslationFile(t *testing.T) {
	t.Skipf("not implemented")
}

func TestLoadTranslationFile(t *testing.T) {
	t.Skipf("not implemented")
}

func TestParseTranslationFileBytes(t *testing.T) {
	t.Skipf("not implemented")
}

func TestAddTranslation(t *testing.T) {
	t.Skipf("not implemented")
}

func TestMustTfunc(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected MustTfunc to panic")
		}
	}()
	New().MustTfunc("invalid")
}

func TestTfuncAndLanguage(t *testing.T) {
	b := New()
	translationID := "translation_id"
	englishLanguage := languageWithTag("en-US")
	frenchLanguage := languageWithTag("fr-FR")
	spanishLanguage := languageWithTag("es")
	chineseLanguage := languageWithTag("zh-hans-cn")
	englishTranslation := addFakeTranslation(t, b, englishLanguage, translationID)
	frenchTranslation := addFakeTranslation(t, b, frenchLanguage, translationID)
	spanishTranslation := addFakeTranslation(t, b, spanishLanguage, translationID)
	chineseTranslation := addFakeTranslation(t, b, chineseLanguage, translationID)

	tests := []struct {
		languageIDs      []string
		result           string
		expectedLanguage *language.Language
	}{
		{
			[]string{"invalid"},
			translationID,
			nil,
		},
		{
			[]string{"invalid", "invalid2"},
			translationID,
			nil,
		},
		{
			[]string{"invalid", "en-US"},
			englishTranslation,
			englishLanguage,
		},
		{
			[]string{"en-US", "invalid"},
			englishTranslation,
			englishLanguage,
		},
		{
			[]string{"en-US", "fr-FR"},
			englishTranslation,
			englishLanguage,
		},
		{
			[]string{"invalid", "es"},
			spanishTranslation,
			spanishLanguage,
		},
		{
			[]string{"zh-CN,fr-XX,es"},
			spanishTranslation,
			spanishLanguage,
		},
		{
			[]string{"fr"},
			frenchTranslation,

			// The language is still "fr" even though the translation is provided by "fr-FR"
			languageWithTag("fr"),
		},
		{
			[]string{"zh"},
			chineseTranslation,

			// The language is still "zh" even though the translation is provided by "zh-hans-cn"
			languageWithTag("zh"),
		},
		{
			[]string{"zh-hans"},
			chineseTranslation,

			// The language is still "zh-hans" even though the translation is provided by "zh-hans-cn"
			languageWithTag("zh-hans"),
		},
		{
			[]string{"zh-hans-cn"},
			chineseTranslation,
			languageWithTag("zh-hans-cn"),
		},
	}

	for i, test := range tests {
		tf, lang, err := b.TfuncAndLanguage(test.languageIDs[0], test.languageIDs[1:]...)
		if err != nil && test.expectedLanguage != nil {
			t.Errorf("Tfunc(%v) = error{%q}; expected no error", test.languageIDs, err)
		}
		if err == nil && test.expectedLanguage == nil {
			t.Errorf("Tfunc(%v) = nil error; expected error", test.languageIDs)
		}
		if result := tf(translationID); result != test.result {
			t.Errorf("translation %d was %s; expected %s", i, result, test.result)
		}
		if (lang == nil && test.expectedLanguage != nil) ||
			(lang != nil && test.expectedLanguage == nil) ||
			(lang != nil && test.expectedLanguage != nil && lang.String() != test.expectedLanguage.String()) {
			t.Errorf("lang %d was %s; expected %s", i, lang, test.expectedLanguage)
		}
	}
}

func TestToMap(t *testing.T) {
	testToMatchesMap(t, "with map", testMap, testMap)
	testToMatchesMap(t, "with struct", testMap, testStruct)
	testToMatchesMap(t, "with pointer", testMap, &testStruct)
}

func TestTranslate(t *testing.T) {
	b := New()
	lang := "en-US"
	translationID := "translation_id"
	b.AddTranslation(languageWithTag(lang), testNewTranslation(t, map[string]interface{}{
		"id":          translationID,
		"translation": "{{.Person}} is {{.Age}} years old.",
	}))
	expected := "Bob is 26 years old."

	tf, err := b.Tfunc(lang)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if result := tf(translationID, testStruct); result != expected {
		t.Errorf("expected '%s', got: '%s'", expected, result)
	}
}

func addFakeTranslation(t *testing.T, b *Bundle, lang *language.Language, translationID string) string {
	translation := fakeTranslation(lang, translationID)
	b.AddTranslation(lang, testNewTranslation(t, map[string]interface{}{
		"id":          translationID,
		"translation": translation,
	}))
	return translation
}

func fakeTranslation(lang *language.Language, translationID string) string {
	return fmt.Sprintf("%s(%s)", lang.Tag, translationID)
}

func testNewTranslation(t *testing.T, data map[string]interface{}) translation.Translation {
	translation, err := translation.NewTranslation(data)
	if err != nil {
		t.Fatal(err)
	}
	return translation
}

func languageWithTag(tag string) *language.Language {
	return language.MustParse(tag)[0]
}

func testToMatchesMap(t *testing.T, key string, expected map[string]interface{}, input interface{}) {
	actual, err := toMap(input)
	if err != nil {
		t.Errorf("toMap failed with (%s)")
	}
	testMatchesMap(t, key, expected, actual)
}

func testMatchesMap(t *testing.T, key string, expected map[string]interface{}, actual map[string]interface{}) {
	for k, v := range expected {
		if actual[k] != v {
			t.Errorf("(%s) expected %v, got: %v", key, v, actual[k])
		}
	}
}
