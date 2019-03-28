package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"golang.org/x/text/language"
)

// UnmarshalFunc unmarshals data into v.
type UnmarshalFunc func(data []byte, v interface{}) error

// MessageFile represents a parsed message file.
type MessageFile struct {
	Path     string
	Tag      language.Tag
	Format   string
	Messages []*Message
}

// ParseMessageFileBytes returns the messages parsed from file.
func ParseMessageFileBytes(buf []byte, path string, unmarshalFuncs map[string]UnmarshalFunc) (*MessageFile, error) {
	lang, format := parsePath(path)
	tag := language.Make(lang)
	messageFile := &MessageFile{
		Path:   path,
		Tag:    tag,
		Format: format,
	}
	if len(buf) == 0 {
		return messageFile, nil
	}
	unmarshalFunc := unmarshalFuncs[messageFile.Format]
	if unmarshalFunc == nil {
		if messageFile.Format == "json" {
			unmarshalFunc = json.Unmarshal
		} else {
			return nil, fmt.Errorf("no unmarshaler registered for %s", messageFile.Format)
		}
	}
	var err error
	var raw interface{}
	if err = unmarshalFunc(buf, &raw); err != nil {
		return nil, err
	}

	if messageFile.Messages, err = recGetMessages(raw, true); err != nil {
		return nil, err
	}

	return messageFile, nil
}

const nestedSeparator = "."

// recGetMessages looks for translation messages inside "raw" parameter,
// scanning nested maps using recursion.
func recGetMessages(raw interface{}, isInitialCall bool) ([]*Message, error) {
	var messages []*Message

	switch data := raw.(type) {
	case string:
		if isInitialCall {
			return nil, errors.New("invalid translation file, expected key-values, got a single value")
		}
		m, err := NewMessage(data)
		if err != nil {
			return nil, err
		}
		messages = []*Message{m}

	case map[string]interface{}:
		if isMessage(raw) {
			m, err := NewMessage(data)
			if err != nil {
				return nil, err
			}
			messages = []*Message{m}
		} else {
			messages = make([]*Message, 0, len(data))
			for id, data := range data {
				// recursively scan map items
				childMessages, err := recGetMessages(data, false)
				if err != nil {
					return nil, err
				}
				for _, m := range childMessages {
					if m.justCreated {
						if m.ID == "" {
							m.ID = id // start with innermost key
						}
					} else {
						m.ID = id + nestedSeparator + m.ID // update ID with each nested key on the way
					}
					m.justCreated = false
					messages = append(messages, m)
				}
			}
		}

	case map[interface{}]interface{}:
		if isMessage(raw) {
			m, err := NewMessage(data)
			if err != nil {
				return nil, err
			}
			messages = []*Message{m}
		} else {
			messages = make([]*Message, 0, len(data))
			for id, data := range data {
				strid, ok := id.(string)
				if !ok {
					return nil, fmt.Errorf("expected key to be string but got %#v", id)
				}
				// recursively scan map items
				childMessages, err := recGetMessages(data, false)
				if err != nil {
					return nil, err
				}
				for _, m := range childMessages {
					if m.ID == "" {
						m.ID = strid // start with innermost key
					} else {
						m.ID = strid + nestedSeparator + m.ID // update ID with each nested key on the way
					}
					messages = append(messages, m)
				}
			}
		}

	case []interface{}:
		// Backward compatibility for v1 file format.
		messages = make([]*Message, 0, len(data))
		for _, data := range data {
			// recursively scan slice items
			childMessages, err := recGetMessages(data, false)
			if err != nil {
				return nil, err
			}
			for _, m := range childMessages {
				messages = append(messages, m)
			}
		}

	default:
		return nil, fmt.Errorf("unsupported file format %T", raw)
	}

	return messages, nil
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
