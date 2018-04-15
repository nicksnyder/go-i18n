package i18n

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/nicksnyder/go-i18n/v2/internal"
	"github.com/nicksnyder/go-i18n/v2/internal/plural"

	"golang.org/x/text/language"
)

// UnmarshalFunc unmarshals data into v.
type UnmarshalFunc = internal.UnmarshalFunc

// Bundle stores a set of messages and pluralization rules.
// Most applications only need a single bundle
// that is initialized early in the application's lifecycle.
type Bundle struct {
	messageTemplates map[language.Tag]map[string]*internal.MessageTemplate
	pluralRules      map[language.Base]*plural.Rule
	unmarshalFuncs   map[string]UnmarshalFunc
	defaultTag       language.Tag
	tags             []language.Tag
	matcher          language.Matcher
}

// NewBundle returns a new bundle that contains the
// CLDR plural rules and a json unmarshaler.
func NewBundle(defaultTag language.Tag) *Bundle {
	b := &Bundle{
		defaultTag:  defaultTag,
		pluralRules: plural.DefaultRules(),
		unmarshalFuncs: map[string]UnmarshalFunc{
			"json": json.Unmarshal,
		},
	}
	b.addTag(defaultTag)
	return b
}

// RegisterPluralRule registers a plural rule for a language base.
// func (b *Bundle) RegisterPluralRule(base language.Base, rule *plural.Rule) {
// 	b.pluralRules[base] = rule
// }

// RegisterUnmarshalFunc registers an UnmarshalFunc for format.
func (b *Bundle) RegisterUnmarshalFunc(format string, unmarshalFunc UnmarshalFunc) {
	if b.unmarshalFuncs == nil {
		b.unmarshalFuncs = make(map[string]UnmarshalFunc)
	}
	b.unmarshalFuncs[format] = unmarshalFunc
}

// LoadMessageFile loads the bytes from path
// and then calls ParseMessageFileBytes.
func (b *Bundle) LoadMessageFile(path string) (*MessageFile, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return b.ParseMessageFileBytes(buf, path)
}

// MustLoadMessageFile is similar to LoadTranslationFile
// except it panics if an error happens.
func (b *Bundle) MustLoadMessageFile(path string) {
	if _, err := b.LoadMessageFile(path); err != nil {
		panic(err)
	}
}

// MessageFile represents a parsed message file.
type MessageFile = internal.MessageFile

// ParseMessageFileBytes parses the bytes in buf to add translations to the bundle.
//
// The format of the file is everything after the last ".".
//
// The language tag of the file is everything after the second to last "." or after the last path separator, but before the format.
func (b *Bundle) ParseMessageFileBytes(buf []byte, path string) (*MessageFile, error) {
	messageFile, err := internal.ParseMessageFileBytes(buf, path, b.unmarshalFuncs)
	if err != nil {
		return nil, err
	}
	if err := b.AddMessages(messageFile.Tag, messageFile.Messages...); err != nil {
		return nil, err
	}
	return messageFile, nil
}

// MustParseMessageFileBytes is similar to ParseMessageFileBytes
// except it panics if an error happens.
func (b *Bundle) MustParseMessageFileBytes(buf []byte, path string) {
	if _, err := b.ParseMessageFileBytes(buf, path); err != nil {
		panic(err)
	}
}

// AddMessages adds messages for a language.
// It is useful if your messages are in a format not supported by ParseMessageFileBytes.
func (b *Bundle) AddMessages(tag language.Tag, messages ...*Message) error {
	if b.pluralRules == nil {
		b.pluralRules = plural.DefaultRules()
	}
	base, _ := tag.Base()
	pluralRule := b.pluralRules[base]
	if pluralRule == nil {
		return fmt.Errorf("no plural rule registered for %s", base)
	}
	b.pluralRules[base] = pluralRule
	if b.messageTemplates == nil {
		b.messageTemplates = map[language.Tag]map[string]*internal.MessageTemplate{}
	}
	if b.messageTemplates[tag] == nil {
		b.messageTemplates[tag] = map[string]*internal.MessageTemplate{}
		b.addTag(tag)
	}
	for _, m := range messages {
		b.messageTemplates[tag][m.ID] = internal.NewMessageTemplate(m)
	}
	return nil
}

// MustAddMessages is similar to AddMessages except it panics if an error happens.
func (b *Bundle) MustAddMessages(tag language.Tag, messages ...*Message) {
	if err := b.AddMessages(tag, messages...); err != nil {
		panic(err)
	}
}

func (b *Bundle) addTag(tag language.Tag) {
	for _, t := range b.tags {
		if t == tag {
			// Tag already exists
			return
		}
	}
	b.tags = append(b.tags, tag)
	b.matcher = language.NewMatcher(b.tags)
}
