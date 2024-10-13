package i18n

import (
	"errors"
	"reflect"
	"sort"
	"testing"

	"golang.org/x/text/language"
	yaml "gopkg.in/yaml.v3"
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
			name: "nested with reserved key",
			file: `{"nested": {"description": {"other": "world"}}}`,
			path: "en.json",
			messageFile: &MessageFile{
				Path:   "en.json",
				Tag:    language.English,
				Format: "json",
				Messages: []*Message{{
					ID:    "nested.description",
					Other: "world",
				}},
			},
		},
		{
			name: "basic test reserved key top level",
			file: `{"other": "world", "foo": "bar"}`,
			path: "en.json",
			err: &mixedKeysError{
				reservedKeys:   []string{"other"},
				unreservedKeys: []string{"foo"},
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
			err:  errInvalidTranslationFile,
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
			err: &mixedKeysError{
				reservedKeys:   []string{"description"},
				unreservedKeys: []string{"dummy"},
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
		{
			name: "YAML empty key test",
			file: `
some-keys:
    non-empty-key: not empty
    empty-key-but-type-specified: ""
    empty-key:
    null-key: null`,
			path:           "en.yaml",
			unmarshalFuncs: map[string]UnmarshalFunc{"yaml": yaml.Unmarshal},
			messageFile: &MessageFile{
				Path:   "en.yaml",
				Tag:    language.English,
				Format: "yaml",
				Messages: []*Message{
					{
						ID:    "some-keys.non-empty-key",
						Other: "not empty",
					},
					{
						ID: "some-keys.empty-key-but-type-specified",
					},
					{
						ID: "some-keys.empty-key",
					},
					{
						ID: "some-keys.null-key",
					},
				},
			},
		},
		{
			name: "YAML number key test",
			file: `
some-keys:
    hello: world
    2: legit`,
			path:           "en.yaml",
			unmarshalFuncs: map[string]UnmarshalFunc{"yaml": yaml.Unmarshal},
			err:            errors.New("expected key to be string but got 2"),
		},
		{
			name: "YAML extra number key test",
			file: `
some-keys:
    other: world
    2: legit`,
			path:           "en.yaml",
			unmarshalFuncs: map[string]UnmarshalFunc{"yaml": yaml.Unmarshal},
			err: &mixedKeysError{
				reservedKeys:   []string{"other"},
				unreservedKeys: []string{"2"},
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual, err := ParseMessageFileBytes([]byte(testCase.file), testCase.path, testCase.unmarshalFuncs)
			if (err == nil && testCase.err != nil) ||
				(err != nil && testCase.err == nil) ||
				(err != nil && testCase.err != nil && err.Error() != testCase.err.Error()) {
				t.Fatalf("expected error %#v; got %#v", testCase.err, err)
			}
			if testCase.messageFile == nil && actual != nil || testCase.messageFile != nil && actual == nil {
				t.Fatalf("expected message file %#v; got %#v", testCase.messageFile, actual)
			}
			if testCase.messageFile != nil {
				if actual.Path != testCase.messageFile.Path {
					t.Errorf("expected path %q; got %q", testCase.messageFile.Path, actual.Path)
				}
				if actual.Tag != testCase.messageFile.Tag {
					t.Errorf("expected tag %q; got %q", testCase.messageFile.Tag, actual.Tag)
				}
				if actual.Format != testCase.messageFile.Format {
					t.Errorf("expected format %q; got %q", testCase.messageFile.Format, actual.Format)
				}
				if !equalMessages(actual.Messages, testCase.messageFile.Messages) {
					t.Errorf("expected %#v; got %#v", deref(testCase.messageFile.Messages), deref(actual.Messages))
				}
			}
		})
	}
}

func deref(mptrs []*Message) []Message {
	messages := make([]Message, len(mptrs))
	for i, m := range mptrs {
		messages[i] = *m
	}
	return messages
}

// equalMessages compares two slices of messages, ignoring private fields and order.
// Sorts both input slices, which are therefore modified by this function.
func equalMessages(m1, m2 []*Message) bool {
	if len(m1) != len(m2) {
		return false
	}

	var less = func(m []*Message) func(int, int) bool {
		return func(i, j int) bool {
			return m[i].ID < m[j].ID
		}
	}
	sort.Slice(m1, less(m1))
	sort.Slice(m2, less(m2))

	for i, m := range m1 {
		if !reflect.DeepEqual(m, m2[i]) {
			return false
		}
	}
	return true
}
