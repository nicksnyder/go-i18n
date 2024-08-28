package i18n_test

import (
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

func TestNotLocalizedInfo(t *testing.T) {
	defer i18n.ResetSingletonContext()

	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.MustParseMessageFileBytes([]byte(`
NotLocalized = "Text is not localized!"
HelloCookie = "Hello Cookie!"
`), "en.toml")

	localizer := i18n.NewLocalizer(bundle, "en-US")
	i18n.SetLocalizerInstance(localizer)

	type test struct {
		name     string
		id       string
		expected string
	}

	// No "not localized" message defined.
	tests := []test{
		{
			name:     "ID - no match found",
			id:       "Hello",
			expected: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := i18n.Localize(test.id)
			if actual != test.expected {
				t.Errorf("expected %q, got %q", test.expected, actual)
			}
		})
	}

	// "not localized" message defined!
	i18n.SetUseNotLocalizedInfo(true)

	tests = []test{
		{
			name:     "ID - no match found",
			id:       "Hello",
			expected: "Text is not localized!",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := i18n.Localize(test.id)
			if actual != test.expected {
				t.Errorf("expected %q, got %q", test.expected, actual)
			}
		})
	}

}

func TestLocalize(t *testing.T) {
	defer i18n.ResetSingletonContext()

	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.MustParseMessageFileBytes([]byte(`
HelloCookie = "Hello Cookie!"
`), "en.toml")

	localizer := i18n.NewLocalizer(bundle, "en-US")
	i18n.SetLocalizerInstance(localizer)

	tests := []struct {
		name     string
		id       string
		expected string
	}{
		{
			name:     "ID - match found",
			id:       "HelloCookie",
			expected: "Hello Cookie!",
		},
		{
			name:     "ID - no match found",
			id:       "Hello",
			expected: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := i18n.Localize(test.id)
			if actual != test.expected {
				t.Errorf("expected %q, got %q", test.expected, actual)
			}
		})
	}
}

func TestLocalizePlural(t *testing.T) {
	defer i18n.ResetSingletonContext()

	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.MustParseMessageFileBytes([]byte(`
HelloCookie = "Hello Cookie!"

[Cookies]
one = "I have {{.PluralCount}} cookie!"
other = "I have {{.PluralCount}} cookies!"
`), "en.toml")

	localizer := i18n.NewLocalizer(bundle, "en-US")
	i18n.SetLocalizerInstance(localizer)

	tests := []struct {
		name     string
		id       string
		count    int
		expected string
	}{
		{
			name:     "hello cookie found",
			id:       "HelloCookie",
			expected: "Hello Cookie!",
		},
		{
			name:     "hello cookie NOT found",
			id:       "HelloCookie",
			count:    1,
			expected: "",
		},
		{
			name:     "1 cookie",
			id:       "Cookies",
			count:    1,
			expected: "I have 1 cookie!",
		},
		{
			name:     "4 cookies",
			id:       "Cookies",
			count:    4,
			expected: "I have 4 cookies!",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := i18n.LocalizePlural(test.id, test.count)
			if actual != test.expected {
				t.Errorf("expected %q, got %q", test.expected, actual)
			}
		})
	}
}

func TestLocalizeTemplate(t *testing.T) {
	defer i18n.ResetSingletonContext()

	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.MustParseMessageFileBytes([]byte(`
HelloCookie = "Hello Cookie!"

[CookiesX]
other = "I have {{.X}} cookies!"
`), "en.toml")

	localizer := i18n.NewLocalizer(bundle, "en-US")
	i18n.SetLocalizerInstance(localizer)

	tests := []struct {
		name     string
		id       string
		key      string
		keyValue any
		expected string
	}{
		{
			name:     "hello cookie found",
			id:       "HelloCookie",
			expected: "Hello Cookie!",
		},
		{
			name:     "found chocolate cookies",
			id:       "CookiesX",
			key:      "X",
			keyValue: "chocolate",
			expected: "I have chocolate cookies!",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := i18n.LocalizeTemplate(test.id, map[string]any{
				test.key: test.keyValue,
			})
			if actual != test.expected {
				t.Errorf("expected %q, got %q", test.expected, actual)
			}
		})
	}
}

func TestLocalizeTemplateSingle(t *testing.T) {
	defer i18n.ResetSingletonContext()

	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.MustParseMessageFileBytes([]byte(`
HelloCookie = "Hello Cookie!"

[CookiesX]
other = "I have {{.X}} cookies!"
`), "en.toml")

	localizer := i18n.NewLocalizer(bundle, "en-US")
	i18n.SetLocalizerInstance(localizer)

	tests := []struct {
		name     string
		id       string
		key      string
		keyValue any
		expected string
	}{
		{
			name:     "hello cookie found",
			id:       "HelloCookie",
			expected: "Hello Cookie!",
		},
		{
			name:     "found chocolate cookies",
			id:       "CookiesX",
			key:      "X",
			keyValue: "chocolate",
			expected: "I have chocolate cookies!",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := i18n.LocalizeTemplateSingle(test.id, test.key, test.keyValue)
			if actual != test.expected {
				t.Errorf("expected %q, got %q", test.expected, actual)
			}
		})
	}
}

func TestLocalizeTemplateSingleWithPlural(t *testing.T) {
	defer i18n.ResetSingletonContext()

	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.MustParseMessageFileBytes([]byte(`
HelloCookie = "Hello Cookie!"

[CookiesX]
zero = "I have no {{.X}} cookie!"
one = "I have {{.PluralCount}} {{.X}} cookie!"
other = "I have {{.PluralCount}} {{.X}} cookies!"
`), "en.toml")

	localizer := i18n.NewLocalizer(bundle, "en-US")
	i18n.SetLocalizerInstance(localizer)

	tests := []struct {
		name     string
		id       string
		count    int
		key      string
		keyValue any
		expected string
	}{
		{
			name:     "hello cookie found",
			id:       "HelloCookie",
			expected: "Hello Cookie!",
		},
		// { // FIXME: Fails because of addPluralRules() in rule_gen.go for language "en". Zero is not defined.
		// 	name:     "found zero chocolate cookies",
		// 	id:       "CookiesX",
		// 	count:    0,
		// 	key:      "X",
		// 	keyValue: "chocolate",
		// 	expected: "I have no chocolate cookie!",
		// },
		{
			name:     "found 1 chocolate cookie",
			id:       "CookiesX",
			count:    1,
			key:      "X",
			keyValue: "chocolate",
			expected: "I have 1 chocolate cookie!",
		},
		{
			name:     "found 4 chocolate cookies",
			id:       "CookiesX",
			count:    4,
			key:      "X",
			keyValue: "chocolate",
			expected: "I have 4 chocolate cookies!",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := i18n.LocalizeTemplateSingleWithPlural(test.id, test.count, test.key, test.keyValue)
			if actual != test.expected {
				t.Errorf("expected %q, got %q", test.expected, actual)
			}
		})
	}
}

func TestLocalizeTemplateX(t *testing.T) {
	defer i18n.ResetSingletonContext()

	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.MustParseMessageFileBytes([]byte(`
HelloCookie = "Hello Cookie!"

[CookiesABC]
other = "I have {{.A}} cookies in my {{.B}} back at {{.C}}!"
`), "en.toml")

	localizer := i18n.NewLocalizer(bundle, "en-US")
	i18n.SetLocalizerInstance(localizer)
	i18n.SetABCParams([]string{"A", "B", "C"})

	tests := []struct {
		name      string
		id        string
		keyValues []any
		expected  string
	}{
		{
			name:     "hello cookie found",
			id:       "HelloCookie",
			expected: "Hello Cookie!",
		},
		{
			name:      "found chocolate cookies",
			id:        "CookiesABC",
			keyValues: []any{"chocolate", "basket", "home"},
			expected:  "I have chocolate cookies in my basket back at home!",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := i18n.LocalizeTemplateX(test.id, test.keyValues...)
			if actual != test.expected {
				t.Errorf("expected %q, got %q", test.expected, actual)
			}
		})
	}
}

func TestLocalizeTemplateXPlural(t *testing.T) {
	defer i18n.ResetSingletonContext()

	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.MustParseMessageFileBytes([]byte(`
HelloCookie = "Hello Cookie!"

[CookiesABC]
other = "I have {{.PluralCount}} {{.A}} cookies in my {{.B}} back at {{.C}}!"
`), "en.toml")

	localizer := i18n.NewLocalizer(bundle, "en-US")
	i18n.SetLocalizerInstance(localizer)
	i18n.SetABCParams([]string{"A", "B", "C"})

	tests := []struct {
		name      string
		id        string
		count     int
		keyValues []any
		expected  string
	}{
		{
			name:     "hello cookie found",
			id:       "HelloCookie",
			expected: "Hello Cookie!",
		},
		{
			name:      "found chocolate cookies",
			id:        "CookiesABC",
			count:     400,
			keyValues: []any{"chocolate", "basket", "home"},
			expected:  "I have 400 chocolate cookies in my basket back at home!",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := i18n.LocalizeTemplateXPlural(test.id, test.count, test.keyValues...)
			if actual != test.expected {
				t.Errorf("expected %q, got %q", test.expected, actual)
			}
		})
	}
}
