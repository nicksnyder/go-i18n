package i18n

import (
	"fmt"
	"sort"
	"strings"
)

// Message is a string that can be localized.
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

	// RightDelim is the right Go template delimiter.
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
	strdata, err := stringMap(v)
	if err != nil {
		return err
	}
	for k, v := range strdata {
		switch strings.ToLower(k) {
		case "id":
			m.ID = v
		case "description":
			m.Description = v
		case "hash":
			m.Hash = v
		case "leftdelim":
			m.LeftDelim = v
		case "rightdelim":
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

type keyTypeErr struct {
	key interface{}
}

func (err *keyTypeErr) Error() string {
	return fmt.Sprintf("expected key to be a string but got %#v", err.key)
}

type valueTypeErr struct {
	value interface{}
}

func (err *valueTypeErr) Error() string {
	return fmt.Sprintf("unsupported type %#v", err.value)
}

func stringMap(v interface{}) (map[string]string, error) {
	switch value := v.(type) {
	case string:
		return map[string]string{
			"other": value,
		}, nil
	case map[string]string:
		return value, nil
	case map[string]interface{}:
		strdata := make(map[string]string, len(value))
		for k, v := range value {
			err := stringSubmap(k, v, strdata)
			if err != nil {
				return nil, err
			}
		}
		return strdata, nil
	case map[interface{}]interface{}:
		strdata := make(map[string]string, len(value))
		for k, v := range value {
			kstr, ok := k.(string)
			if !ok {
				return nil, &keyTypeErr{key: k}
			}
			err := stringSubmap(kstr, v, strdata)
			if err != nil {
				return nil, err
			}
		}
		return strdata, nil
	default:
		return nil, &valueTypeErr{value: value}
	}
}

func stringSubmap(k string, v interface{}, strdata map[string]string) error {
	if k == "translation" {
		switch vt := v.(type) {
		case string:
			strdata["other"] = vt
		default:
			v1Message, err := stringMap(v)
			if err != nil {
				return err
			}
			for kk, vv := range v1Message {
				strdata[kk] = vv
			}
		}
		return nil
	}

	switch vt := v.(type) {
	case string:
		strdata[k] = vt
		return nil
	case nil:
		return nil
	default:
		return fmt.Errorf("expected value for key %q be a string but got %#v", k, v)
	}
}

// messageKeys are the keys allowed in a message.
var messageKeys = map[string]struct{}{
	"id":          {},
	"description": {},
	"hash":        {},
	"leftdelim":   {},
	"rightdelim":  {},
	"zero":        {},
	"one":         {},
	"two":         {},
	"few":         {},
	"many":        {},
	"other":       {},
	"translation": {},
}

func isMessageKey(key string, val any) bool {
	lk := strings.ToLower(key)
	if _, ok := messageKeys[lk]; ok {
		if key == "translation" {
			return true
		}
		if _, ok := val.(string); ok {
			return true
		}
	}
	return false
}

// isMessage returns true if v contains an "id" key combined with only message keys and false if v contains only non-message keys.
// It returns an error if v contains both message and non-message keys.
func isMessage(v interface{}) (bool, error) {
	switch data := v.(type) {
	case nil, string:
		return true, nil
	case map[string]interface{}:
		if _, ok := data["other"]; !ok {
			return false, nil
		}
		messageKeys := make([]string, 0, len(messageKeys))
		nonMessageKeys := make([]string, 0, len(data))
		for k, v := range data {
			if isMessageKey(k, v) {
				messageKeys = append(messageKeys, k)
			} else {
				nonMessageKeys = append(nonMessageKeys, k)
			}
		}
		if len(nonMessageKeys) > 0 {
			return false, &mixedKeysError{
				messageKeys:    messageKeys,
				nonMessageKeys: nonMessageKeys,
			}
		}
		return true, nil
	case map[interface{}]interface{}:
		if _, ok := data["other"]; !ok {
			return false, nil
		}
		messageKeys := make([]string, 0, len(messageKeys))
		nonMessageKeys := make([]string, 0, len(data))
		for key, v := range data {
			k, ok := key.(string)
			if !ok {
				nonMessageKeys = append(nonMessageKeys, fmt.Sprintf("%+v", key))
			} else if isMessageKey(k, v) {
				messageKeys = append(messageKeys, k)
			} else {
				nonMessageKeys = append(nonMessageKeys, k)
			}
		}
		if len(nonMessageKeys) > 0 {
			return false, &mixedKeysError{
				messageKeys:    messageKeys,
				nonMessageKeys: nonMessageKeys,
			}
		}
		return true, nil
	}
	return false, nil
}

type mixedKeysError struct {
	messageKeys    []string
	nonMessageKeys []string
}

func (e *mixedKeysError) Error() string {
	sort.Strings(e.messageKeys)
	sort.Strings(e.nonMessageKeys)
	return fmt.Sprintf("message keys %v mixed with non-message keys %v", e.messageKeys, e.nonMessageKeys)
}
