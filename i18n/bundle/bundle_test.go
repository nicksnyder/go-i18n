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

func createBenchmarkTranslateFunc(b *testing.B) func(data interface{}) {
	bundle := New()
	lang := "en-US"
	translationID := "translation_id"
	translation, err := translation.NewTranslation(map[string]interface{}{
		"id": translationID,
		"translation": map[string]interface{}{
			"one":   "{{.Person}} is {{.Count}} year old.",
			"other": "{{.Person}} is {{.Count}} years old.",
		},
	})
	if err != nil {
		b.Fatal(err)
	}
	bundle.AddTranslation(languageWithTag(lang), translation)
	expected := "Bob is 26 years old."

	tf, err := bundle.Tfunc(lang)
	if err != nil {
		b.Fatal(err)
	}

	return func(data interface{}) {
		if result := tf(translationID, 26, data); result != expected {
			b.Fatalf("expected %q, got %q", expected, result)
		}
	}
}

func BenchmarkTranslateWithMap(b *testing.B) {
	data := map[string]interface{}{
		"Person": "Bob",
	}
	tf := createBenchmarkTranslateFunc(b)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tf(data)
	}
}

func BenchmarkTranslateWithStruct(b *testing.B) {
	data := struct{ Person string }{Person: "Bob"}
	tf := createBenchmarkTranslateFunc(b)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tf(data)
	}
}

func BenchmarkTranslateWithStructPointer(b *testing.B) {
	data := &struct{ Person string }{Person: "Bob"}
	tf := createBenchmarkTranslateFunc(b)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tf(data)
	}
}
