package internal

import (
	"fmt"
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

func TestMessageTemplatePluralFormMissed(t *testing.T) {
	// test special case when appropriate plural form is not defined
	mt := NewMessageTemplate(&Message{ID: "MissPluralForm", One: "Hello World", Other: "Hello World"})
	// this plural form is not defined in NewMessageTemplate call
	var undefined plural.Form = plural.Few
	// it's OK to receive error on Execute exit, otherwise, one way or another, raise panic
	_, err := mt.Execute(undefined, nil, nil)
	if err == nil {
		panic(fmt.Sprintf("Message template %v should return error when search for a %q plural form", mt, undefined))
	}
}
