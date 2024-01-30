package template

import (
	"bytes"
	"strings"
	texttemplate "text/template"
)

type ParsedTemplate interface {
	Execute(data any) (string, error)
}

type Engine interface {
	ParseTemplate(src, leftDelim, rightDelim string) (ParsedTemplate, error)

	// Cacheable returns true if the ParsedTemplate returned by ParseTemplate is safe to cache.
	Cacheable() bool
}

// IdentityEngine is an Engine that does no parsing and returns template strings unchanged.
type IdentityEngine struct{}

func (IdentityEngine) Cacheable() bool {
	// Caching is not necessary because ParseTemplate is cheap.
	return false
}

func (IdentityEngine) ParseTemplate(src, leftDelim, rightDelim string) (ParsedTemplate, error) {
	return &identityParsedTemplate{src: src}, nil
}

type identityParsedTemplate struct {
	src string
}

func (t *identityParsedTemplate) Execute(data any) (string, error) {
	return t.src, nil
}

// TextEngine is an Engine for text/template.
type TextEngine struct {
	LeftDelim  string
	RightDelim string
	Funcs      texttemplate.FuncMap
	Option     string
}

func (te *TextEngine) Cacheable() bool {
	return te.Funcs == nil
}

func (te *TextEngine) ParseTemplate(src, leftDelim, rightDelim string) (ParsedTemplate, error) {
	if leftDelim == "" {
		leftDelim = te.LeftDelim
	}
	if leftDelim == "" {
		leftDelim = "{{"
	}
	if !strings.Contains(src, leftDelim) {
		// Fast path to avoid parsing a template that has no actions.
		return &identityParsedTemplate{src: src}, nil
	}

	if rightDelim == "" {
		rightDelim = te.RightDelim
	}
	if rightDelim == "" {
		rightDelim = "}}"
	}

	tmpl, err := texttemplate.New("").Delims(leftDelim, rightDelim).Funcs(te.Funcs).Parse(src)
	if err != nil {
		return nil, err
	}
	return &parsedTextTemplate{tmpl: tmpl}, nil
}

type parsedTextTemplate struct {
	tmpl *texttemplate.Template
}

func (t *parsedTextTemplate) Execute(data any) (string, error) {
	var buf bytes.Buffer
	if err := t.tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
