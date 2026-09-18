package i18n

import (
	"reflect"
	"testing"

	"github.com/nicksnyder/go-i18n/v2/internal/plural"
)

func TestMessageTemplate(t *testing.T) {
	mt := NewMessageTemplate(&Message{ID: "HelloWorld", Other: "Hello World"})
	if mt.PluralTemplates[plural.Other].Src != "Hello World" {
		t.Fatal(mt.PluralTemplates)
	}
}

func TestNilMessageTemplate(t *testing.T) {
	if mt := NewMessageTemplate(&Message{ID: "HelloWorld"}); mt != nil {
		t.Fatal(mt)
	}
}

func TestEmptyOtherMessageTemplateSkippedByDefault(t *testing.T) {
	m, err := NewMessage("")
	if err != nil {
		t.Fatal(err)
	}
	if mt := NewMessageTemplate(m); mt != nil {
		t.Fatal("NewMessageTemplate still skips empty other")
	}
	mt := newMessageTemplate(m, true)
	if mt == nil {
		t.Fatal("expected template for explicit empty other")
	}
	if src := mt.PluralTemplates[plural.Other].Src; src != "" {
		t.Fatalf("expected empty src; got %q", src)
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
