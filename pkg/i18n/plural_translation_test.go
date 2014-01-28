package i18n

import (
	"reflect"
	"testing"
)

func mustTemplate(t *testing.T, src string) *template {
	tmpl, err := newTemplate(src)
	if err != nil {
		t.Fatal(err)
	}
	return tmpl
}

func pluralTranslationFixture(t *testing.T, id string, pluralCategories ...PluralCategory) *pluralTranslation {
	templates := make(map[PluralCategory]*template, len(pluralCategories))
	for _, pc := range pluralCategories {
		templates[pc] = mustTemplate(t, string(pc))
	}
	return &pluralTranslation{id, templates}
}

func verifyDeepEqual(t *testing.T, actual, expected interface{}) {
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("\n%#v\nnot equal to expected value\n%#v", actual, expected)
	}
}

func TestPluralTranslationMerge(t *testing.T) {
	pt := pluralTranslationFixture(t, "id", One, Other)
	oneTemplate, otherTemplate := pt.templates[One], pt.templates[Other]

	pt.Merge(pluralTranslationFixture(t, "id"))
	verifyDeepEqual(t, pt.templates, map[PluralCategory]*template{
		One:   oneTemplate,
		Other: otherTemplate,
	})

	pt2 := pluralTranslationFixture(t, "id", One, Two)
	pt.Merge(pt2)
	verifyDeepEqual(t, pt.templates, map[PluralCategory]*template{
		One:   pt2.templates[One],
		Two:   pt2.templates[Two],
		Other: otherTemplate,
	})
}

/* Test implementations from old idea

func TestCopy(t *testing.T) {
	ls := &LocalizedString{
		Id:          "id",
		Translation: testingTemplate(t, "translation {{.Hello}}"),
		Translations: map[PluralCategory]*template{
			One:   testingTemplate(t, "plural {{.One}}"),
			Other: testingTemplate(t, "plural {{.Other}}"),
		},
	}

	c := ls.Copy()
	delete(c.Translations, One)
	if _, ok := ls.Translations[One]; !ok {
		t.Errorf("deleting plural translation from copy deleted it from the original")
	}
	c.Translations[Two] = testingTemplate(t, "plural {{.Two}}")
	if _, ok := ls.Translations[Two]; ok {
		t.Errorf("adding plural translation to copy added it to the original")
	}
}

func TestNormalize(t *testing.T) {
	oneTemplate := testingTemplate(t, "one {{.One}}")
	ls := &LocalizedString{
		Translation: testingTemplate(t, "single {{.Single}}"),
		Translations: map[PluralCategory]*template{
			One: oneTemplate,
			Two: testingTemplate(t, "two {{.Two}}"),
		},
	}
	ls.Normalize(LanguageWithCode("en"))
	if ls.Translation != nil {
		t.Errorf("ls.Translation is %#v; expected nil", ls.Translation)
	}
	if actual := ls.Translations[Two]; actual != nil {
		t.Errorf("ls.Translation[Two] is %#v; expected nil", actual)
	}
	if actual := ls.Translations[One]; actual != oneTemplate {
		t.Errorf("ls.Translations[One] is %#v; expected %#v", actual, oneTemplate)
	}
	if _, ok := ls.Translations[Other]; !ok {
		t.Errorf("ls.Translations[Other] shouldn't be empty")
	}
}

func TestMergeTranslation(t *testing.T) {
	ls := &LocalizedString{}

	translation := testingTemplate(t, "one {{.Hello}}")
	ls.Merge(&LocalizedString{
		Translation: translation,
	})
	if ls.Translation != translation {
		t.Errorf("expected %#v; got %#v", translation, ls.Translation)
	}

	ls.Merge(&LocalizedString{})
	if ls.Translation != translation {
		t.Errorf("expected %#v; got %#v", translation, ls.Translation)
	}

	translation = testingTemplate(t, "two {{.Hello}}")
	ls.Merge(&LocalizedString{
		Translation: translation,
	})
	if ls.Translation != translation {
		t.Errorf("expected %#v; got %#v", translation, ls.Translation)
	}
}

func TestMergeTranslations(t *testing.T) {
	ls := &LocalizedString{}

	oneTemplate := testingTemplate(t, "one {{.One}}")
	otherTemplate := testingTemplate(t, "other {{.Other}}")
	ls.Merge(&LocalizedString{
		Translations: map[PluralCategory]*template{
			One:   oneTemplate,
			Other: otherTemplate,
		},
	})
	if actual := ls.Translations[One]; actual != oneTemplate {
		t.Errorf("ls.Translations[One] expected %#v; got %#v", oneTemplate, actual)
	}
	if actual := ls.Translations[Other]; actual != otherTemplate {
		t.Errorf("ls.Translations[Other] expected %#v; got %#v", otherTemplate, actual)
	}

	ls.Merge(&LocalizedString{
		Translations: map[PluralCategory]*template{},
	})
	if actual := ls.Translations[One]; actual != oneTemplate {
		t.Errorf("ls.Translations[One] expected %#v; got %#v", oneTemplate, actual)
	}
	if actual := ls.Translations[Other]; actual != otherTemplate {
		t.Errorf("ls.Translations[Other] expected %#v; got %#v", otherTemplate, actual)
	}

	twoTemplate := testingTemplate(t, "two {{.Two}}")
	otherTemplate = testingTemplate(t, "second other {{.Other}}")
	ls.Merge(&LocalizedString{
		Translations: map[PluralCategory]*template{
			Two:   twoTemplate,
			Other: otherTemplate,
		},
	})
	if actual := ls.Translations[One]; actual != oneTemplate {
		t.Errorf("ls.Translations[One] expected %#v; got %#v", oneTemplate, actual)
	}
	if actual := ls.Translations[Two]; actual != twoTemplate {
		t.Errorf("ls.Translations[Two] expected %#v; got %#v", twoTemplate, actual)
	}
	if actual := ls.Translations[Other]; actual != otherTemplate {
		t.Errorf("ls.Translations[Other] expected %#v; got %#v", otherTemplate, actual)
	}
}

func TestMissingTranslations(t *testing.T) {
	en := LanguageWithCode("en")

	tests := []struct {
		localizedString *LocalizedString
		language        *Language
		expected        bool
	}{
		{
			&LocalizedString{},
			en,
			true,
		},
		{
			&LocalizedString{Translation: testingTemplate(t, "single {{.Single}}")},
			en,
			false,
		},
		{
			&LocalizedString{
				Translation: testingTemplate(t, "single {{.Single}}"),
				Translations: map[PluralCategory]*template{
					One: testingTemplate(t, "one {{.One}}"),
				}},
			en,
			true,
		},
		{
			&LocalizedString{Translations: map[PluralCategory]*template{
				One: testingTemplate(t, "one {{.One}}"),
			}},
			en,
			true,
		},
		{
			&LocalizedString{Translations: map[PluralCategory]*template{
				One:   nil,
				Other: nil,
			}},
			en,
			true,
		},
		{
			&LocalizedString{Translations: map[PluralCategory]*template{
				One:   testingTemplate(t, ""),
				Other: testingTemplate(t, ""),
			}},
			en,
			true,
		},
		{
			&LocalizedString{Translations: map[PluralCategory]*template{
				One:   testingTemplate(t, "one {{.One}}"),
				Other: testingTemplate(t, "other {{.Other}}"),
			}},
			en,
			false,
		},
	}

	for _, tt := range tests {
		if actual := tt.localizedString.MissingTranslations(tt.language); actual != tt.expected {
			t.Errorf("expected %t got %t for %s, %#v",
				tt.expected, actual, tt.language.code, tt.localizedString)
		}
	}
}

func TestHasTranslations(t *testing.T) {
	en := LanguageWithCode("en")

	tests := []struct {
		localizedString *LocalizedString
		language        *Language
		expected        bool
	}{
		{
			&LocalizedString{},
			en,
			false,
		},
		{
			&LocalizedString{Translation: testingTemplate(t, "single {{.Single}}")},
			en,
			true,
		},
		{
			&LocalizedString{
				Translation:  testingTemplate(t, "single {{.Single}}"),
				Translations: map[PluralCategory]*template{}},
			en,
			false,
		},
		{
			&LocalizedString{Translations: map[PluralCategory]*template{
				One: testingTemplate(t, "one {{.One}}"),
			}},
			en,
			true,
		},
		{
			&LocalizedString{Translations: map[PluralCategory]*template{
				Two: testingTemplate(t, "two {{.Two}}"),
			}},
			en,
			false,
		},
		{
			&LocalizedString{Translations: map[PluralCategory]*template{
				One: nil,
			}},
			en,
			false,
		},
		{
			&LocalizedString{Translations: map[PluralCategory]*template{
				One: testingTemplate(t, ""),
			}},
			en,
			false,
		},
	}

	for _, tt := range tests {
		if actual := tt.localizedString.HasTranslations(tt.language); actual != tt.expected {
			t.Errorf("expected %t got %t for %s, %#v",
				tt.expected, actual, tt.language.code, tt.localizedString)
		}
	}
}

func testingTemplate(t *testing.T, src string) *template {
	tmpl, err := newTemplate(src)
	if err != nil {
		t.Fatal(err)
	}
	return tmpl
}
*/
