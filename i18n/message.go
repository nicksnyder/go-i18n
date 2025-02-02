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

var reservedKeys = map[string]struct{}{
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

func isReserved(key string, val any) bool {
	lk := strings.ToLower(key)
	if _, ok := reservedKeys[lk]; ok {
		if key == "translation" {
			return true
		}
		if _, ok := val.(string); ok {
			return true
		}
	}
	return false
}

// isMessage returns true if v contains only message keys and false if it contains no message keys.
// It returns an error if v contains both message and non-message keys.
// - {"message": {"description": "world"}} is a message
// - {"error": {"description": "world", "foo": "bar"}} is an error
// - {"notmessage": {"description": {"hello": "world"}}} is not a message
// - {"notmessage": {"foo": "bar"}} is not a message
func isMessage(v interface{}) (bool, error) {
	switch data := v.(type) {
	case nil, string:
		return true, nil
	case map[string]interface{}:
		reservedKeys := make([]string, 0, len(reservedKeys))
		unreservedKeys := make([]string, 0, len(data))
		for k, v := range data {
			if isReserved(k, v) {
				reservedKeys = append(reservedKeys, k)
			} else {
				unreservedKeys = append(unreservedKeys, k)
			}
		}
		hasReservedKeys := len(reservedKeys) > 0
		if hasReservedKeys && len(unreservedKeys) > 0 {
			return false, &mixedKeysError{
				reservedKeys:   reservedKeys,
				unreservedKeys: unreservedKeys,
			}
		}
		return hasReservedKeys, nil
	case map[interface{}]interface{}:
		reservedKeys := make([]string, 0, len(reservedKeys))
		unreservedKeys := make([]string, 0, len(data))
		for key, v := range data {
			k, ok := key.(string)
			if !ok {
				unreservedKeys = append(unreservedKeys, fmt.Sprintf("%+v", key))
			} else if isReserved(k, v) {
				reservedKeys = append(reservedKeys, k)
			} else {
				unreservedKeys = append(unreservedKeys, k)
			}
		}
		hasReservedKeys := len(reservedKeys) > 0
		if hasReservedKeys && len(unreservedKeys) > 0 {
			return false, &mixedKeysError{
				reservedKeys:   reservedKeys,
				unreservedKeys: unreservedKeys,
			}
		}
		return hasReservedKeys, nil
	}
	return false, nil
}

type mixedKeysError struct {
	reservedKeys   []string
	unreservedKeys []string
}

func (e *mixedKeysError) Error() string {
	sort.Strings(e.reservedKeys)
	sort.Strings(e.unreservedKeys)
	return fmt.Sprintf("reserved keys %v mixed with unreserved keys %v", e.reservedKeys, e.unreservedKeys)
}
