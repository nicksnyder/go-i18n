package i18n

import (
	"reflect"
	"testing"
)

func TestNewMessage(t *testing.T) {
	tests := []struct {
		name    string
		data    interface{}
		message *Message
		err     error
	}{
		{
			name: "string",
			data: "other",
			message: &Message{
				Other: "other",
			},
		},
		{
			name: "nil value",
			data: map[interface{}]interface{}{
				"ID":    "id",
				"Zero":  nil,
				"Other": "other",
			},
			message: &Message{
				ID:    "id",
				Other: "other",
			},
		},
		{
			name: "map[string]string",
			data: map[string]string{
				"ID":          "id",
				"Hash":        "hash",
				"Description": "description",
				"LeftDelim":   "leftdelim",
				"RightDelim":  "rightdelim",
				"Zero":        "zero",
				"One":         "one",
				"Two":         "two",
				"Few":         "few",
				"Many":        "many",
				"Other":       "other",
			},
			message: &Message{
				ID:          "id",
				Hash:        "hash",
				Description: "description",
				LeftDelim:   "leftdelim",
				RightDelim:  "rightdelim",
				Zero:        "zero",
				One:         "one",
				Two:         "two",
				Few:         "few",
				Many:        "many",
				Other:       "other",
			},
		},
		{
			name: "map[string]interface{}",
			data: map[string]interface{}{
				"ID":          "id",
				"Hash":        "hash",
				"Description": "description",
				"LeftDelim":   "leftdelim",
				"RightDelim":  "rightdelim",
				"Zero":        "zero",
				"One":         "one",
				"Two":         "two",
				"Few":         "few",
				"Many":        "many",
				"Other":       "other",
			},
			message: &Message{
				ID:          "id",
				Hash:        "hash",
				Description: "description",
				LeftDelim:   "leftdelim",
				RightDelim:  "rightdelim",
				Zero:        "zero",
				One:         "one",
				Two:         "two",
				Few:         "few",
				Many:        "many",
				Other:       "other",
			},
		},
		{
			name: "map[interface{}]interface{}",
			data: map[interface{}]interface{}{
				"ID":          "id",
				"Hash":        "hash",
				"Description": "description",
				"LeftDelim":   "leftdelim",
				"RightDelim":  "rightdelim",
				"Zero":        "zero",
				"One":         "one",
				"Two":         "two",
				"Few":         "few",
				"Many":        "many",
				"Other":       "other",
			},
			message: &Message{
				ID:          "id",
				Hash:        "hash",
				Description: "description",
				LeftDelim:   "leftdelim",
				RightDelim:  "rightdelim",
				Zero:        "zero",
				One:         "one",
				Two:         "two",
				Few:         "few",
				Many:        "many",
				Other:       "other",
			},
		},
		{
			name: "map[int]int",
			data: map[interface{}]interface{}{
				1: 2,
			},
			err: &keyTypeErr{key: 1},
		},
		{
			name: "int",
			data: 1,
			err:  &valueTypeErr{value: 1},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := NewMessage(test.data)
			if !reflect.DeepEqual(err, test.err) {
				t.Fatalf("expected %#v; got %#v", test.err, err)
			}
			if !reflect.DeepEqual(actual, test.message) {
				t.Fatalf("\nexpected\n%#v\ngot\n%#v", test.message, actual)
			}
		})
	}
}

func TestKeyTypeErr(t *testing.T) {
	expected := "expected key to be a string but got 1"
	if actual := (&keyTypeErr{key: 1}).Error(); actual != expected {
		t.Fatalf("expected %#v; got %#v", expected, actual)
	}
}

func TestValueTypeErr(t *testing.T) {
	expected := "unsupported type 1"
	if actual := (&valueTypeErr{value: 1}).Error(); actual != expected {
		t.Fatalf("expected %#v; got %#v", expected, actual)
	}
}
