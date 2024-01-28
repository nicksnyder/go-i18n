package i18n

import (
	"sync"
)

// Template stores the template for a string.
type Template struct {
	Src        string
	LeftDelim  string
	RightDelim string

	parseOnce      sync.Once
	parsedTemplate ParsedTemplate
	parseError     error
}

func (t *Template) execute(engine TemplateEngine, data interface{}) (string, error) {
	var pt ParsedTemplate
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
