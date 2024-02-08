package internal

import (
	"sync"

	"github.com/KosyanMedia/go-i18n/i18n/template"
)

// Template stores the template for a string and a cached version of the parsed template if they are cacheable.
type Template struct {
	Src        string
	LeftDelim  string
	RightDelim string

	parseOnce      sync.Once
	parsedTemplate template.ParsedTemplate
	parseError     error
}

func (t *Template) Execute(parser template.Parser, data interface{}) (string, error) {
	var pt template.ParsedTemplate
	var err error
	if parser.Cacheable() {
		t.parseOnce.Do(func() {
			t.parsedTemplate, t.parseError = parser.Parse(t.Src, t.LeftDelim, t.RightDelim)
		})
		pt, err = t.parsedTemplate, t.parseError
	} else {
		pt, err = parser.Parse(t.Src, t.LeftDelim, t.RightDelim)
	}

	if err != nil {
		return "", err
	}
	return pt.Execute(data)
}
