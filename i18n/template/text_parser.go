package template

import (
	"bytes"
	"strings"
	"text/template"
)

// TextParser is a Parser that uses text/template.
type TextParser struct {
	LeftDelim  string
	RightDelim string
	Funcs      template.FuncMap
	Option     string
}

func (te *TextParser) Cacheable() bool {
	return te.Funcs == nil
}

func (te *TextParser) Parse(src, leftDelim, rightDelim string) (ParsedTemplate, error) {
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

	option := "missingkey=default"
	if te.Option != "" {
		option = te.Option
	}

	tmpl, err := template.New("").Delims(leftDelim, rightDelim).Option(option).Funcs(te.Funcs).Parse(src)
	if err != nil {
		return nil, err
	}
	return &parsedTextTemplate{tmpl: tmpl}, nil
}

type parsedTextTemplate struct {
	tmpl *template.Template
}

func (t *parsedTextTemplate) Execute(data any) (string, error) {
	var buf bytes.Buffer
	if err := t.tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
