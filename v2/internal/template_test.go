package internal

import (
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {
	tmpl := &Template{Src: "hello"}
	if err := tmpl.parse("", ""); err != nil {
		t.Fatal(err)
	}
	if tmpl.ParseErr == nil {
		t.Fatal("expected non-nil parse error")
	}
	if tmpl.Template == nil {
		t.Fatal("expected non-nil template")
	}
}

func TestParseError(t *testing.T) {
	expectedErr := fmt.Errorf("expected error")
	tmpl := &Template{ParseErr: &expectedErr}
	if err := tmpl.parse("", ""); err != expectedErr {
		t.Fatalf("expected %#v; got %#v", expectedErr, err)
	}
}
