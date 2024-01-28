package i18n

import (
	"bytes"
	"strings"
	texttemplate "text/template"
)

type ParsedTemplate interface {
	Execute(data any) (string, error)
}

type TemplateEngine interface {
	// Cacheable returns true if the ParsedTemplate returned by ParseTemplate is safe to cache.
	Cacheable() bool

	ParseTemplate(src, leftDelim, rightDelim string) (ParsedTemplate, error)
}

type NoTemplateEngine struct{}

func (NoTemplateEngine) Cacheable() bool {
	// Caching is not necessary because ParseTemplate is cheap.
	return false
}

func (NoTemplateEngine) ParseTemplate(src, leftDelim, rightDelim string) (ParsedTemplate, error) {
	return &identityParsedTemplate{src: src}, nil
}

type TextTemplateEngine struct {
	LeftDelim  string
	RightDelim string
	Funcs      texttemplate.FuncMap
	Option     string
}

func (te *TextTemplateEngine) Cacheable() bool {
	return te.Funcs == nil
}

func (te *TextTemplateEngine) ParseTemplate(src, leftDelim, rightDelim string) (ParsedTemplate, error) {
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

type identityParsedTemplate struct {
	src string
}

func (t *identityParsedTemplate) Execute(data any) (string, error) {
	return t.src, nil
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

// type TemplateEngine interface {
// 	Execute(src string, data any) (string, error)
// }

// type NoTemplateEngine struct{}

// func (*NoTemplateEngine) Execute(src string, data any) (string, error) {
// 	return src, nil
// }

// type TextTemplateEngine2 struct {
// 	LeftDelim  string
// 	RightDelim string
// 	Funcs      texttemplate.FuncMap
// 	Option     string

// 	cache      map[string]*executeResult
// 	cacheMutex sync.RWMutex
// }

// type executeResult struct {
// 	tmpl *texttemplate.Template
// 	err  error
// }

// func (t *TextTemplateEngine2) Execute(src string, data any) (string, error) {
// 	tmpl, err := t.getTemplate(src)
// 	if err != nil {
// 		return "", err
// 	}
// 	var buf bytes.Buffer
// 	if err := tmpl.Execute(&buf, data); err != nil {
// 		return "", err
// 	}
// 	return buf.String(), nil
// }

// func (t *TextTemplateEngine2) getTemplate(template string) (*texttemplate.Template, error) {
// 	// It is not safe to use the cache if t.Funcs or t.Option is set.
// 	if t.Funcs != nil || t.Option != "" {
// 		return texttemplate.New("").Delims(t.LeftDelim, t.RightDelim).Funcs(t.Funcs).Option(t.Option).Parse(template)
// 	}

// 	// If there is a cached result, return it.
// 	t.cacheMutex.RLock()
// 	result := t.cache[template]
// 	t.cacheMutex.RUnlock()
// 	if result != nil {
// 		return result.tmpl, result.err
// 	}

// 	// Parse the template and save it to the cache
// 	tmpl, err := texttemplate.New("").Delims(t.LeftDelim, t.RightDelim).Parse(template)
// 	r := &executeResult{
// 		tmpl: tmpl,
// 		err:  err,
// 	}
// 	t.cacheMutex.Lock()
// 	t.cache[template] = r
// 	t.cacheMutex.Unlock()

// 	return tmpl, err
// }
