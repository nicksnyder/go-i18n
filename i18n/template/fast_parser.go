package template

import (
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/valyala/fasttemplate"
)

// FastParser is a Parser that uses valyala/fasttemplate for simple placeholder substitution.
//
// This parser was introduced primarily for TinyGo compatibility. TinyGo has limited support
// for the standard library's text/template and html/template packages due to their heavy use
// of reflection and other features that are difficult to compile to WebAssembly or embedded
// targets. FastParser provides a lightweight alternative that works reliably in TinyGo
// environments while still supporting the common use case of simple placeholder substitution.
//
// FastParser is also faster than TextParser for templates that only use simple {{.Key}}
// placeholders without any functions or complex template logic.
//
// Supported syntax:
//   - {{.Key}} - dot-prefixed placeholder (standard go-i18n style)
//   - {{Key}} - direct key reference
//
// Limitations (use TextParser for these features):
//   - Template functions (e.g., {{printf "%s" .Name}})
//   - Conditionals (e.g., {{if .Cond}}...{{end}})
//   - Range loops (e.g., {{range .Items}}...{{end}})
//   - Nested field access (e.g., {{.User.Name}})
//   - Method calls (e.g., {{.Name.String}})
type FastParser struct {
	LeftDelim  string
	RightDelim string
}

func (fp *FastParser) Cacheable() bool {
	return true
}

func (fp *FastParser) Parse(src, leftDelim, rightDelim string) (ParsedTemplate, error) {
	if leftDelim == "" {
		leftDelim = fp.LeftDelim
	}
	if leftDelim == "" {
		leftDelim = "{{"
	}

	if !strings.Contains(src, leftDelim) {
		return &identityParsedTemplate{src: src}, nil
	}

	if rightDelim == "" {
		rightDelim = fp.RightDelim
	}
	if rightDelim == "" {
		rightDelim = "}}"
	}

	tmpl, err := fasttemplate.NewTemplate(src, leftDelim, rightDelim)
	if err != nil {
		return nil, err
	}

	return &parsedFastTemplate{
		tmpl:       tmpl,
		leftDelim:  leftDelim,
		rightDelim: rightDelim,
	}, nil
}

type parsedFastTemplate struct {
	tmpl       *fasttemplate.Template
	leftDelim  string
	rightDelim string
}

func (t *parsedFastTemplate) Execute(data any) (string, error) {
	m := toTagFuncMap(data)
	return t.tmpl.ExecuteFuncString(func(w io.Writer, tag string) (int, error) {
		if fn, ok := m[tag]; ok {
			return fn(w, tag)
		}
		return w.Write([]byte(t.leftDelim + tag + t.rightDelim))
	}), nil
}

type tagFunc = func(w io.Writer, tag string) (int, error)

func toTagFuncMap(data any) map[string]tagFunc {
	if data == nil {
		return map[string]tagFunc{}
	}

	switch d := data.(type) {
	case map[string]interface{}:
		return mapToTagFuncs(d)
	case map[string]string:
		m := make(map[string]tagFunc, len(d)*2)
		for k, v := range d {
			val := v
			fn := func(w io.Writer, _ string) (int, error) {
				return w.Write([]byte(val))
			}
			m[k] = fn
			m["."+k] = fn
		}
		return m
	default:
		return structToTagFuncs(data)
	}
}

func mapToTagFuncs(m map[string]interface{}) map[string]tagFunc {
	result := make(map[string]tagFunc, len(m)*2)
	for k, v := range m {
		fn := valueToTagFunc(v)
		result[k] = fn
		if !strings.HasPrefix(k, ".") {
			result["."+k] = fn
		}
	}
	return result
}

func valueToTagFunc(v interface{}) tagFunc {
	switch val := v.(type) {
	case string:
		return func(w io.Writer, _ string) (int, error) {
			return w.Write([]byte(val))
		}
	case []byte:
		return func(w io.Writer, _ string) (int, error) {
			return w.Write(val)
		}
	case fmt.Stringer:
		s := val.String()
		return func(w io.Writer, _ string) (int, error) {
			return w.Write([]byte(s))
		}
	default:
		s := fmt.Sprint(v)
		return func(w io.Writer, _ string) (int, error) {
			return w.Write([]byte(s))
		}
	}
}

func structToTagFuncs(data any) map[string]tagFunc {
	if data == nil {
		return map[string]tagFunc{}
	}

	if s, ok := data.(fmt.Stringer); ok {
		str := s.String()
		fn := func(w io.Writer, _ string) (int, error) {
			return w.Write([]byte(str))
		}
		return map[string]tagFunc{"": fn, ".": fn}
	}

	return reflectStructToTagFuncs(data)
}

func reflectStructToTagFuncs(data any) map[string]tagFunc {
	if data == nil {
		return map[string]tagFunc{}
	}

	val := reflect.ValueOf(data)
	for val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return map[string]tagFunc{}
		}
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return map[string]tagFunc{}
	}

	typ := val.Type()
	m := make(map[string]tagFunc, val.NumField()*2)

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		if field.PkgPath != "" {
			continue
		}
		fieldVal := val.Field(i)
		if fieldVal.CanInterface() {
			fn := valueToTagFunc(fieldVal.Interface())
			m[field.Name] = fn
			m["."+field.Name] = fn
		}
	}

	return m
}
