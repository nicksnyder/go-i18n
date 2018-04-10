package i18n

import (
	"fmt"
	"strings"
)

// Message is a message that can be localized.
type Message struct {
	// ID uniquely identifies the message.
	ID string

	// Hash uniquely identifies the content of the message
	// that this message was translated from.
	Hash string

	// Description describes the message to give additional
	// context to translators that may be relevant for translation.
	Description string

	// LeftDelim is the left Go template delimiter.
	LeftDelim string

	// RightDelim is the right Go template delimiter.``
	RightDelim string

	// Zero is the content of the message for the CLDR plural form "zero".
	Zero string

	// One is the content of the message for the CLDR plural form "one".
	One string

	// Two is the content of the message for the CLDR plural form "two".
	Two string

	// Few is the content of the message for the CLDR plural form "few".
	Few string

	// Many is the content of the message for the CLDR plural form "many".
	Many string

	// Other is the content of the message for the CLDR plural form "other".
	Other string
}

// NewMessage parses data and returns a new message.
func NewMessage(data interface{}) (*Message, error) {
	m := &Message{}
	if err := m.unmarshalInterface(data); err != nil {
		return nil, err
	}
	return m, nil
}

// MustNewMessage is similar to NewMessage except it panics if an error happens.
func MustNewMessage(data interface{}) *Message {
	m, err := NewMessage(data)
	if err != nil {
		panic(err)
	}
	return m
}

// unmarshalInterface unmarshals a message from data.
func (m *Message) unmarshalInterface(v interface{}) error {
	strdata := make(map[string]string)
	switch value := v.(type) {
	case string:
		strdata["other"] = value
	case map[string]string:
		strdata = value
	case map[string]interface{}:
		for k, v := range value {
			vstr, ok := v.(string)
			if !ok {
				return fmt.Errorf("expected value for key %q be a string but got %#v", k, v)
			}
			strdata[k] = vstr
		}
	case map[interface{}]interface{}:
		for k, v := range value {
			kstr, ok := k.(string)
			if !ok {
				return fmt.Errorf("expected key to be a string but got %#v", k)
			}
			vstr, ok := v.(string)
			if !ok {
				return fmt.Errorf("expected value for key %q be a string but got %#v", k, v)
			}
			strdata[kstr] = vstr
		}
	default:
		return fmt.Errorf("unsupported type %#v", value)
	}
	for k, v := range strdata {
		switch strings.ToLower(k) {
		case "id":
			m.ID = v
		case "description":
			m.Description = v
		case "hash":
			m.Hash = v
		case "leftDelim":
			m.LeftDelim = v
		case "rightDelim":
			m.RightDelim = v
		case "zero":
			m.Zero = v
		case "one":
			m.One = v
		case "two":
			m.Two = v
		case "few":
			m.Few = v
		case "many":
			m.Many = v
		case "other":
			m.Other = v
		}
	}
	return nil
}
