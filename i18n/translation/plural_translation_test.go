package translation

import (
	"github.com/nicksnyder/go-i18n/i18n/plural"
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

func pluralTranslationFixture(t *testing.T, id string, pluralCategories ...plural.Category) *pluralTranslation {
	templates := make(map[plural.Category]*template, len(pluralCategories))
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
	pt := pluralTranslationFixture(t, "id", plural.One, plural.Other)
	oneTemplate, otherTemplate := pt.templates[plural.One], pt.templates[plural.Other]

	pt.Merge(pluralTranslationFixture(t, "id"))
	verifyDeepEqual(t, pt.templates, map[plural.Category]*template{
		plural.One:   oneTemplate,
		plural.Other: otherTemplate,
	})

	pt2 := pluralTranslationFixture(t, "id", plural.One, plural.Two)
	pt.Merge(pt2)
	verifyDeepEqual(t, pt.templates, map[plural.Category]*template{
		plural.One:   pt2.templates[plural.One],
		plural.Two:   pt2.templates[plural.Two],
		plural.Other: otherTemplate,
	})
}

/* Test implementations from old idea

func TestCopy(t *testing.T) {
	ls := &LocalizedString{
		ID:          "id",
		Translation: testingTemplate(t, "translation {{.Hello}}"),
		Translations: map[plural.Category]*template{
			plural.One:   testingTemplate(t, "plural {{.One}}"),
			plural.Other: testingTemplate(t, "plural {{.Other}}"),
		},
	}

	c := ls.Copy()
	delete(c.Translations, plural.One)
	if _, ok := ls.Translations[plural.One]; !ok {
		t.Errorf("deleting plural translation from copy deleted it from the original")
	}
	c.Translations[plural.Two] = testingTemplate(t, "plural {{.Two}}")
	if _, ok := ls.Translations[plural.Two]; ok {
		t.Errorf("adding plural translation to copy added it to the original")
	}
}

func TestNormalize(t *testing.T) {
	oneTemplate := testingTemplate(t, "one {{.One}}")
	ls := &LocalizedString{
		Translation: testingTemplate(t, "single {{.Single}}"),
		Translations: map[plural.Category]*template{
			plural.One: oneTemplate,
			plural.Two: testingTemplate(t, "two {{.Two}}"),
		},
	}
	ls.Normalize(LanguageWithCode("en"))
	if ls.Translation != nil {
		t.Errorf("ls.Translation is %#v; expected nil", ls.Translation)
	}
	if actual := ls.Translations[plural.Two]; actual != nil {
		t.Errorf("ls.Translation[plural.Two] is %#v; expected nil", actual)
	}
	if actual := ls.Translations[plural.One]; actual != oneTemplate {
		t.Errorf("ls.Translations[plural.One] is %#v; expected %#v", actual, oneTemplate)
	}
	if _, ok := ls.Translations[plural.Other]; !ok {
		t.Errorf("ls.Translations[plural.Other] shouldn't be empty")
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
		Translations: map[plural.Category]*template{
			plural.One:   oneTemplate,
			plural.Other: otherTemplate,
		},
	})
	if actual := ls.Translations[plural.One]; actual != oneTemplate {
		t.Errorf("ls.Translations[plural.One] expected %#v; got %#v", oneTemplate, actual)
	}
	if actual := ls.Translations[plural.Other]; actual != otherTemplate {
		t.Errorf("ls.Translations[plural.Other] expected %#v; got %#v", otherTemplate, actual)
	}

	ls.Merge(&LocalizedString{
		Translations: map[plural.Category]*template{},
	})
	if actual := ls.Translations[plural.One]; actual != oneTemplate {
		t.Errorf("ls.Translations[plural.One] expected %#v; got %#v", oneTemplate, actual)
	}
	if actual := ls.Translations[plural.Other]; actual != otherTemplate {
		t.Errorf("ls.Translations[plural.Other] expected %#v; got %#v", otherTemplate, actual)
	}

	twoTemplate := testingTemplate(t, "two {{.Two}}")
	otherTemplate = testingTemplate(t, "second other {{.Other}}")
	ls.Merge(&LocalizedString{
		Translations: map[plural.Category]*template{
			plural.Two:   twoTemplate,
			plural.Other: otherTemplate,
		},
	})
	if actual := ls.Translations[plural.One]; actual != oneTemplate {
		t.Errorf("ls.Translations[plural.One] expected %#v; got %#v", oneTemplate, actual)
	}
	if actual := ls.Translations[plural.Two]; actual != twoTemplate {
		t.Errorf("ls.Translations[plural.Two] expected %#v; got %#v", twoTemplate, actual)
	}
	if actual := ls.Translations[plural.Other]; actual != otherTemplate {
		t.Errorf("ls.Translations[plural.Other] expected %#v; got %#v", otherTemplate, actual)
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
				Translations: map[plural.Category]*template{
					plural.One: testingTemplate(t, "one {{.One}}"),
				}},
			en,
			true,
		},
		{
			&LocalizedString{Translations: map[plural.Category]*template{
				plural.One: testingTemplate(t, "one {{.One}}"),
			}},
			en,
			true,
		},
		{
			&LocalizedString{Translations: map[plural.Category]*template{
				plural.One:   nil,
				plural.Other: nil,
			}},
			en,
			true,
		},
		{
			&LocalizedString{Translations: map[plural.Category]*template{
				plural.One:   testingTemplate(t, ""),
				plural.Other: testingTemplate(t, ""),
			}},
			en,
			true,
		},
		{
			&LocalizedString{Translations: map[plural.Category]*template{
				plural.One:   testingTemplate(t, "one {{.One}}"),
				plural.Other: testingTemplate(t, "other {{.Other}}"),
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
				Translations: map[plural.Category]*template{}},
			en,
			false,
		},
		{
			&LocalizedString{Translations: map[plural.Category]*template{
				plural.One: testingTemplate(t, "one {{.One}}"),
			}},
			en,
			true,
		},
		{
			&LocalizedString{Translations: map[plural.Category]*template{
				plural.Two: testingTemplate(t, "two {{.Two}}"),
			}},
			en,
			false,
		},
		{
			&LocalizedString{Translations: map[plural.Category]*template{
				plural.One: nil,
			}},
			en,
			false,
		},
		{
			&LocalizedString{Translations: map[plural.Category]*template{
				plural.One: testingTemplate(t, ""),
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
