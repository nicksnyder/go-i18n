package i18n

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/BurntSushi/toml"
	"golang.org/x/text/language"
	yaml "gopkg.in/yaml.v2"
)

var simpleMessage = MustNewMessage(map[string]string{
	"id":    "simple",
	"other": "simple translation",
})

var detailMessage = MustNewMessage(map[string]string{
	"id":          "detail",
	"description": "detail description",
	"other":       "detail translation",
})

var everythingMessage = MustNewMessage(map[string]string{
	"id":          "everything",
	"description": "everything description",
	"zero":        "zero translation",
	"one":         "one translation",
	"two":         "two translation",
	"few":         "few translation",
	"many":        "many translation",
	"other":       "other translation",
})

func TestConcurrentAccess(t *testing.T) {
	bundle := NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.MustParseMessageFileBytes([]byte(`
# Comment
hello = "world"
`), "en.toml")

	count := 10
	errch := make(chan error, count)
	for i := 0; i < count; i++ {
		go func() {
			localized := NewLocalizer(bundle, "en").MustLocalize(&LocalizeConfig{MessageID: "hello"})
			if localized != "world" {
				errch <- fmt.Errorf(`expected "world"; got %q`, localized)
			} else {
				errch <- nil
			}
		}()
	}

	for i := 0; i < count; i++ {
		if err := <-errch; err != nil {
			t.Fatal(err)
		}
	}
}

func TestPseudoLanguage(t *testing.T) {
	bundle := NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	expected := "nuqneH"
	bundle.MustParseMessageFileBytes([]byte(`
# Comment
hello = "`+expected+`"
`), "art-x-klingon.toml")
	{
		localized, err := NewLocalizer(bundle, "art-x-klingon").Localize(&LocalizeConfig{MessageID: "hello"})
		if err != nil {
			t.Fatal(err)
		}
		if localized != expected {
			t.Fatalf("expected %q\ngot %q", expected, localized)
		}
	}
	{
		localized, err := NewLocalizer(bundle, "art").Localize(&LocalizeConfig{MessageID: "hello"})
		if err != nil {
			t.Fatal(err)
		}
		if localized != expected {
			t.Fatalf("expected %q\ngot %q", expected, localized)
		}
	}
	{
		expected := ""
		localized, err := NewLocalizer(bundle, "en").Localize(&LocalizeConfig{MessageID: "hello"})
		if err == nil {
			t.Fatal(err)
		}
		if localized != expected {
			t.Fatalf("expected %q\ngot %q", expected, localized)
		}
	}
}

func TestJSON(t *testing.T) {
	bundle := NewBundle(language.English)
	bundle.MustParseMessageFileBytes([]byte(`{
	"simple": "simple translation",
	"detail": {
		"description": "detail description",
		"other": "detail translation"
	},
	"everything": {
		"description": "everything description",
		"zero": "zero translation",
		"one": "one translation",
		"two": "two translation",
		"few": "few translation",
		"many": "many translation",
		"other": "other translation"
	}
}`), "en-US.json")

	expectMessage(t, bundle, language.AmericanEnglish, "simple", simpleMessage)
	expectMessage(t, bundle, language.AmericanEnglish, "detail", detailMessage)
	expectMessage(t, bundle, language.AmericanEnglish, "everything", everythingMessage)
}

func TestYAML(t *testing.T) {
	bundle := NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)
	bundle.MustParseMessageFileBytes([]byte(`
# Comment
simple: simple translation

# Comment
detail:
  description: detail description 
  other: detail translation

# Comment
everything:
  description: everything description
  zero: zero translation
  one: one translation
  two: two translation
  few: few translation
  many: many translation
  other: other translation
`), "en-US.yaml")

	expectMessage(t, bundle, language.AmericanEnglish, "simple", simpleMessage)
	expectMessage(t, bundle, language.AmericanEnglish, "detail", detailMessage)
	expectMessage(t, bundle, language.AmericanEnglish, "everything", everythingMessage)
}

func TestTOML(t *testing.T) {
	bundle := NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.MustParseMessageFileBytes([]byte(`
# Comment
simple = "simple translation"

# Comment
[detail]
description = "detail description"
other = "detail translation"

# Comment
[everything]
description = "everything description"
zero = "zero translation"
one = "one translation"
two = "two translation"
few = "few translation"
many = "many translation"
other = "other translation"
`), "en-US.toml")

	expectMessage(t, bundle, language.AmericanEnglish, "simple", simpleMessage)
	expectMessage(t, bundle, language.AmericanEnglish, "detail", detailMessage)
	expectMessage(t, bundle, language.AmericanEnglish, "everything", everythingMessage)
}

func TestV1Format(t *testing.T) {
	bundle := NewBundle(language.English)
	bundle.MustParseMessageFileBytes([]byte(`[
	{
		"id": "simple",
		"translation": "simple translation"
	},
	{
		"id": "everything",
		"translation": {
			"zero": "zero translation",
			"one": "one translation",
			"two": "two translation",
			"few": "few translation",
			"many": "many translation",
			"other": "other translation"
		}
	}
]
`), "en-US.json")

	expectMessage(t, bundle, language.AmericanEnglish, "simple", simpleMessage)
	e := *everythingMessage
	e.Description = ""
	expectMessage(t, bundle, language.AmericanEnglish, "everything", &e)
}

func TestV1FlatFormat(t *testing.T) {
	bundle := NewBundle(language.English)
	bundle.MustParseMessageFileBytes([]byte(`{
	"simple": {
		"other": "simple translation"
	},
	"everything": {
		"zero": "zero translation",
		"one": "one translation",
		"two": "two translation",
		"few": "few translation",
		"many": "many translation",
		"other": "other translation"
	}
}
`), "en-US.json")

	expectMessage(t, bundle, language.AmericanEnglish, "simple", simpleMessage)
	e := *everythingMessage
	e.Description = ""
	expectMessage(t, bundle, language.AmericanEnglish, "everything", &e)
}

func expectMessage(t *testing.T, bundle *Bundle, tag language.Tag, messageID string, message *Message) {
	expected := NewMessageTemplate(message)
	actual := bundle.messageTemplates[tag][messageID]
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("bundle.MessageTemplates[%q][%q] = %#v; want %#v", tag, messageID, actual, expected)
	}
}
