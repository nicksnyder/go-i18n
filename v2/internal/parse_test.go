package internal

import (
	"errors"
	"testing"

	"golang.org/x/text/language"
	yaml "gopkg.in/yaml.v2"
)

func TestParseMessageFileBytes(t *testing.T) {
	testCases := []struct {
		name           string
		file           string
		path           string
		unmarshalFuncs map[string]UnmarshalFunc
		messageFile    *MessageFile
		err            error
	}{
		{
			name: "basic test",
			file: `{"hello": "world"}`,
			path: "en.json",
			messageFile: &MessageFile{
				Path:   "en.json",
				Tag:    language.English,
				Format: "json",
				Messages: []*Message{{
					ID:    "hello",
					Other: "world",
				}},
			},
		},
		{
			name: "basic test with dot separator in key",
			file: `{"prepended.hello": "world"}`,
			path: "en.json",
			messageFile: &MessageFile{
				Path:   "en.json",
				Tag:    language.English,
				Format: "json",
				Messages: []*Message{{
					ID:    "prepended.hello",
					Other: "world",
				}},
			},
		},
		{
			name: "invalid test (no key)",
			file: `"hello"`,
			path: "en.json",
			err:  errors.New("invalid translation file, expected key-values, got a single value"),
		},
		{
			name: "nested test",
			file: `{"nested": {"hello": "world"}}`,
			path: "en.json",
			messageFile: &MessageFile{
				Path:   "en.json",
				Tag:    language.English,
				Format: "json",
				Messages: []*Message{{
					ID:    "nested.hello",
					Other: "world",
				}},
			},
		},
		{
			name: "basic test with description",
			file: `{"notnested": {"description": "world"}}`,
			path: "en.json",
			messageFile: &MessageFile{
				Path:   "en.json",
				Tag:    language.English,
				Format: "json",
				Messages: []*Message{{
					ID:          "notnested",
					Description: "world",
				}},
			},
		},
		{
			name: "basic test with id",
			file: `{"key": {"id": "forced.id"}}`,
			path: "en.json",
			messageFile: &MessageFile{
				Path:   "en.json",
				Tag:    language.English,
				Format: "json",
				Messages: []*Message{{
					ID: "forced.id",
				}},
			},
		},
		{
			name: "basic test with description and dummy",
			file: `{"notnested": {"description": "world", "dummy": "nothing"}}`,
			path: "en.json",
			messageFile: &MessageFile{
				Path:   "en.json",
				Tag:    language.English,
				Format: "json",
				Messages: []*Message{{
					ID:          "notnested",
					Description: "world",
				}},
			},
		},
		{
			name: "deeply nested test",
			file: `{"outer": {"nested": {"inner": "value"}}}`,
			path: "en.json",
			messageFile: &MessageFile{
				Path:   "en.json",
				Tag:    language.English,
				Format: "json",
				Messages: []*Message{{
					ID:    "outer.nested.inner",
					Other: "value",
				}},
			},
		},
		{
			name: "multiple nested test",
			file: `{"nested": {"hello": "world", "bye": "all"}}`,
			path: "en.json",
			messageFile: &MessageFile{
				Path:   "en.json",
				Tag:    language.English,
				Format: "json",
				Messages: []*Message{{
					ID:    "nested.hello",
					Other: "world",
				}, {
					ID:    "nested.bye",
					Other: "all",
				}},
			},
		},
		{
			name: "YAML nested test",
			file: `
outer:
    nested:
        inner: "value"`,
			path:           "en.yaml",
			unmarshalFuncs: map[string]UnmarshalFunc{"yaml": yaml.Unmarshal},
			messageFile: &MessageFile{
				Path:   "en.yaml",
				Tag:    language.English,
				Format: "yaml",
				Messages: []*Message{{
					ID:    "outer.nested.inner",
					Other: "value",
				}},
			},
		},
	}
	for _, testCase := range testCases {
		actual, err := ParseMessageFileBytes([]byte(testCase.file), testCase.path, testCase.unmarshalFuncs)
		if (err == nil && testCase.err != nil) ||
			(err != nil && testCase.err == nil) ||
			(err != nil && testCase.err != nil && err.Error() != testCase.err.Error()) {
			t.Errorf("%s failed: expected error %#v; got %#v", testCase.name, testCase.err, err)
			continue
		}
		if actual == nil {
			continue
		}
		if actual.Path != testCase.messageFile.Path {
			t.Errorf("%s failed: expected path %q; got %q", testCase.name, testCase.messageFile.Path, actual.Path)
			continue
		}
		if actual.Tag != testCase.messageFile.Tag {
			t.Errorf("%s failed: expected tag %q; got %q", testCase.name, testCase.messageFile.Tag, actual.Tag)
			continue
		}
		if actual.Format != testCase.messageFile.Format {
			t.Errorf("%s failed: expected format %q; got %q", testCase.name, testCase.messageFile.Format, actual.Format)
			continue
		}
		if !EqualMessages(actual.Messages, testCase.messageFile.Messages) {
			t.Errorf("%s failed: expected %#v; got %#v", testCase.name, testCase.messageFile.Messages, actual.Messages)
			continue
		}
	}
}
