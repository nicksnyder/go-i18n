package i18n

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"golang.org/x/text/language"
)

// UnmarshalFunc unmarshals data into v.
type UnmarshalFunc func(data []byte, v interface{}) error

// Bundle stores all messages and pluralization rules.
// Generally, your application should only need a single bundle
// that is initialized early in your application's lifecycle.
type Bundle struct {
	// MessageTemplates maps language tags to language ids to message templates.
	MessageTemplates map[language.Tag]map[string]*MessageTemplate

	// pluralRules maps language tags to their plural rules.
	pluralRules map[language.Base]*PluralRule

	// UnmarshalFuncs maps file formats to unmarshal functions.
	UnmarshalFuncs map[string]UnmarshalFunc

	defaultTag language.Tag
	tags       []language.Tag
	matcher    language.Matcher
}

// NewBundle returns a new bundle that contains the
// CLDR plural rules and a json unmarshaler.
func NewBundle(defaultTag language.Tag) *Bundle {
	b := &Bundle{
		defaultTag:  defaultTag,
		pluralRules: DefaultPluralRules(),
		UnmarshalFuncs: map[string]UnmarshalFunc{
			"json": json.Unmarshal,
		},
	}
	b.addTag(defaultTag)
	return b
}

// RegisterUnmarshalFunc registers an UnmarshalFunc for format.
func (b *Bundle) RegisterUnmarshalFunc(format string, unmarshalFunc UnmarshalFunc) {
	if b.UnmarshalFuncs == nil {
		b.UnmarshalFuncs = make(map[string]UnmarshalFunc)
	}
	b.UnmarshalFuncs[format] = unmarshalFunc
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
type MessageFile struct {
	Path     string
	Tag      language.Tag
	Format   string
	Messages []*Message
}

// ParseMessageFileBytes parses the bytes in buf to add translations to the bundle.
// It is useful for parsing translation files embedded with go-bindata.
//
// The format of the file is everything after the last ".".
//
// The language tag of path is the last match of LanguageTagRegex, excluding everything after the last ".".
func (b *Bundle) ParseMessageFileBytes(buf []byte, path string) (*MessageFile, error) {
	lang, format := parsePath(path)
	tag, err := language.Parse(lang)
	if err != nil {
		return nil, err
	}
	messageFile := &MessageFile{
		Path:   path,
		Tag:    tag,
		Format: format,
	}
	if len(buf) == 0 {
		return messageFile, nil
	}
	var unmarshalFunc UnmarshalFunc
	if b.UnmarshalFuncs != nil {
		unmarshalFunc = b.UnmarshalFuncs[messageFile.Format]
	}
	if unmarshalFunc == nil {
		return nil, fmt.Errorf("no unmarshaler registered for %s", messageFile.Format)
	}
	var raw interface{}
	if err := unmarshalFunc(buf, &raw); err != nil {
		return nil, err
	}
	switch data := raw.(type) {
	case map[string]interface{}:
		messageFile.Messages = make([]*Message, 0, len(data))
		for id, data := range data {
			m, err := NewMessage(data)
			if err != nil {
				return nil, err
			}
			m.ID = id
			messageFile.Messages = append(messageFile.Messages, m)
		}
	case map[interface{}]interface{}:
		messageFile.Messages = make([]*Message, 0, len(data))
		for id, data := range data {
			strid, ok := id.(string)
			if !ok {
				return nil, fmt.Errorf("expected key to be string but got %#v", id)
			}
			m, err := NewMessage(data)
			if err != nil {
				return nil, err
			}
			m.ID = strid
			messageFile.Messages = append(messageFile.Messages, m)
		}
	case []interface{}:
		// Backward compatability for v1 file format.
		messageFile.Messages = make([]*Message, 0, len(data))
		for _, data := range data {
			m, err := NewMessage(data)
			if err != nil {
				return nil, err
			}
			messageFile.Messages = append(messageFile.Messages, m)
		}
	default:
		return nil, fmt.Errorf("unsupported file format %T", raw)
	}
	if err := b.AddMessages(messageFile.Tag, messageFile.Messages...); err != nil {
		return nil, err
	}
	return messageFile, nil
}

func parsePath(path string) (langTag, format string) {
	formatStartIdx := -1
	for i := len(path) - 1; i >= 0; i-- {
		c := path[i]
		if os.IsPathSeparator(c) {
			if formatStartIdx != -1 {
				langTag = path[i+1 : formatStartIdx]
			}
			return
		}
		if path[i] == '.' {
			if formatStartIdx != -1 {
				langTag = path[i+1 : formatStartIdx]
				return
			}
			if formatStartIdx == -1 {
				format = path[i+1:]
				formatStartIdx = i
			}
		}
	}
	if formatStartIdx != -1 {
		langTag = path[:formatStartIdx]
	}
	return
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
		b.pluralRules = DefaultPluralRules()
	}
	base, _ := tag.Base()
	pluralRule := b.pluralRules[base]
	if pluralRule == nil {
		return fmt.Errorf("no plural rule registered for %s", base)
	}
	b.pluralRules[base] = pluralRule
	if b.MessageTemplates == nil {
		b.MessageTemplates = map[language.Tag]map[string]*MessageTemplate{}
	}
	if b.MessageTemplates[tag] == nil {
		b.MessageTemplates[tag] = map[string]*MessageTemplate{}
		b.addTag(tag)
	}
	for _, m := range messages {
		b.MessageTemplates[tag][m.ID] = NewMessageTemplate(m)
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
