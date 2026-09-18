package i18n

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"testing/fstest"

	"github.com/BurntSushi/toml"
	yaml "go.yaml.in/yaml/v3"
	"golang.org/x/text/language"
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
	"leftDelim":   "<<",
	"rightDelim":  ">>",
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

func TestParseMessageFileBytesWithTag(t *testing.T) {
	bundle := NewBundle(language.English)
	mf, err := bundle.ParseMessageFileBytesWithTag([]byte(`{"hello": "world"}`), "locales/en/translation.json", language.French)
	if err != nil {
		t.Fatal(err)
	}
	if mf.Tag != language.French {
		t.Fatalf("expected tag %q; got %q", language.French, mf.Tag)
	}
	expectMessage(t, bundle, language.French, "hello", MustNewMessage(map[string]string{
		"id":    "hello",
		"other": "world",
	}))
}

func TestLoadMessageFileWithTag(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "locales", "en", "translation.json")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(`{"hello": "world"}`), 0o644); err != nil {
		t.Fatal(err)
	}

	inferred := NewBundle(language.English)
	mf, err := inferred.LoadMessageFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if mf.Tag != language.English {
		t.Fatalf("inferred tag: got %q want en", mf.Tag)
	}

	overridden := NewBundle(language.English)
	mf, err = overridden.LoadMessageFileWithTag(path, language.French)
	if err != nil {
		t.Fatal(err)
	}
	if mf.Tag != language.French {
		t.Fatalf("override tag: got %q want fr", mf.Tag)
	}
	expectMessage(t, overridden, language.French, "hello", MustNewMessage(map[string]string{
		"id":    "hello",
		"other": "world",
	}))

	fsys := fstest.MapFS{
		"locales/de/translation.json": {Data: []byte(`{"hello": "hallo"}`)},
	}
	fromFS := NewBundle(language.English)
	mf, err = fromFS.LoadMessageFileFSWithTag(fsys, "locales/de/translation.json", language.German)
	if err != nil {
		t.Fatal(err)
	}
	if mf.Tag != language.German {
		t.Fatalf("fs override tag: got %q want de", mf.Tag)
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
		"other": "other translation",
		"leftDelim": "<<",
		"rightDelim": ">>"
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
  leftDelim: "<<"
  rightDelim: ">>"
`), "en-US.yaml")

	expectMessage(t, bundle, language.AmericanEnglish, "simple", simpleMessage)
	expectMessage(t, bundle, language.AmericanEnglish, "detail", detailMessage)
	expectMessage(t, bundle, language.AmericanEnglish, "everything", everythingMessage)
}

func TestInvalidYAML(t *testing.T) {
	bundle := NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)
	_, err := bundle.ParseMessageFileBytes([]byte(`
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
  leftDelim: "<<"
  rightDelmin: ">>"
  garbage: something

description: translation
`), "en-US.yaml")

	expectedErr := &mixedKeysError{
		reservedKeys:   []string{"description"},
		unreservedKeys: []string{"detail", "everything", "simple"},
	}
	if err == nil {
		t.Fatalf("expected error %#v; got nil", expectedErr)
	}
	if err.Error() != expectedErr.Error() {
		t.Fatalf("expected error %q; got %q", expectedErr, err)
	}
	if c := len(bundle.messageTemplates); c > 0 {
		t.Fatalf("expected no message templates in bundle; got %d", c)
	}
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
leftDelim = "<<"
rightDelim = ">>"
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
	expectMessage(t, bundle, language.AmericanEnglish, "everything", newV1EverythingMessage())
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
	expectMessage(t, bundle, language.AmericanEnglish, "everything", newV1EverythingMessage())
}

func expectMessage(t *testing.T, bundle *Bundle, tag language.Tag, messageID string, message *Message) {
	expected := NewMessageTemplate(message)
	actual := bundle.messageTemplates[tag][messageID]
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("bundle.MessageTemplates[%q][%q]\ngot  %#v\nwant %#v", tag, messageID, actual, expected)
	}
}

func newV1EverythingMessage() *Message {
	e := *everythingMessage
	e.Description = ""
	e.LeftDelim = ""
	e.RightDelim = ""
	return &e
}
