package locale

import (
	"github.com/nicksnyder/go-i18n/i18n/language"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		localeID string
		lang     *language.Language
	}{
		{"en-US", language.LanguageWithID("en")},
		{"en_US", language.LanguageWithID("en")},
		{"zh-CN", language.LanguageWithID("zh")},
		{"zh-TW", language.LanguageWithID("zh")},
		{"pt-BR", language.LanguageWithID("pt-BR")},
		{"pt_BR", language.LanguageWithID("pt-BR")},
		{"pt-PT", language.LanguageWithID("pt")},
		{"pt_PT", language.LanguageWithID("pt")},
		{"zh-Hans-CN", nil},
		{"zh-Hant-TW", nil},
		{"xx-Yyen-US", nil},
		{"en US", nil},
		{"en-US-en-US", nil},
		{".en-US..en-US.", nil},
	}
	for _, test := range tests {
		loc, err := New(test.localeID)
		if loc == nil && test.lang != nil {
			t.Errorf("New(%q) = <nil>, %q; expected %q, <nil>", test.localeID, err, test.lang)
		}
		if loc != nil && loc.Language != test.lang {
			t.Errorf("New(%q) = %q; expected %q", test.localeID, loc.Language, test.lang)
		}
	}
}
