package internal

import (
	"sync"

	"github.com/nicksnyder/go-i18n/v2/i18n/template"
)

// Template stores the template for a string.
type Template struct {
	Src        string
	LeftDelim  string
	RightDelim string

	parseOnce      sync.Once
	parsedTemplate template.ParsedTemplate
	parseError     error
}

func (t *Template) Execute(engine template.Engine, data interface{}) (string, error) {
	var pt template.ParsedTemplate
	var err error
	if engine.Cacheable() {
		t.parseOnce.Do(func() {
			t.parsedTemplate, t.parseError = engine.ParseTemplate(t.Src, t.LeftDelim, t.RightDelim)
		})
		pt, err = t.parsedTemplate, t.parseError
	} else {
		pt, err = engine.ParseTemplate(t.Src, t.LeftDelim, t.RightDelim)
	}

	if err != nil {
		return "", err
	}
	return pt.Execute(data)
}
