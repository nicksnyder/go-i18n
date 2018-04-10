package i18n

import (
	"encoding/json"
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

func TestJSON(t *testing.T) {
	var bundle Bundle
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
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
	var bundle Bundle
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
	var bundle Bundle
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

func expectMessage(t *testing.T, bundle Bundle, tag language.Tag, messageID string, message *Message) {
	expected := NewMessageTemplate(message)
	actual := bundle.MessageTemplates[tag][messageID]
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("bundle.MessageTemplates[%q][%q] = %#v; want %#v", tag, messageID, actual, expected)
	}
}
