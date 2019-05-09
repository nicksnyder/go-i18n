package i18n

import (
	"reflect"
	"testing"

	"github.com/nicksnyder/go-i18n/v2/internal/plural"
)

func TestMessageTemplate(t *testing.T) {
	mt := NewMessageTemplate(&Message{ID: "HelloWorld", Other: "Hello World"})
	if mt.PluralTemplates[plural.Other].Src != "Hello World" {
		panic(mt.PluralTemplates)
	}
}

func TestNilMessageTemplate(t *testing.T) {
	if mt := NewMessageTemplate(&Message{ID: "HelloWorld"}); mt != nil {
		panic(mt)
	}
}

func TestMessageTemplatePluralFormMissing(t *testing.T) {
	mt := NewMessageTemplate(&Message{ID: "HelloWorld", Other: "Hello World"})
	s, err := mt.Execute(plural.Few, nil, nil)
	if s != "" {
		t.Errorf("expected %q; got %q", "", s)
	}
	expectedErr := pluralFormNotFoundError{pluralForm: plural.Few, messageID: "HelloWorld"}
	if !reflect.DeepEqual(err, expectedErr) {
		t.Errorf("expected error %#v; got %#v", expectedErr, err)
	}
}
