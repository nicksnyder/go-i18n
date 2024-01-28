package i18n

import (
	"strings"
	"testing"
	"text/template"
)

func TestExecute(t *testing.T) {
	tests := []struct {
		template *Template
		engine   TemplateEngine
		data     interface{}
		result   string
		err      string
		noallocs bool
	}{
		{
			template: &Template{
				Src: "hello",
			},
			result:   "hello",
			noallocs: true,
		},
		{
			template: &Template{
				Src: "hello {{.Noun}}",
			},
			data: map[string]string{
				"Noun": "world",
			},
			result: "hello world",
		},
		{
			template: &Template{
				Src: "hello {{world}}",
			},
			engine: &TextTemplateEngine{
				Funcs: template.FuncMap{
					"world": func() string {
						return "world"
					},
				},
			},
			result: "hello world",
		},
		{
			template: &Template{
				Src: "hello {{",
			},
			err:      "unclosed action",
			noallocs: true,
		},
	}

	for _, test := range tests {
		t.Run(test.template.Src, func(t *testing.T) {
			if test.engine == nil {
				test.engine = &TextTemplateEngine{}
			}
			result, err := test.template.execute(test.engine, test.data)
			if actual := str(err); !strings.Contains(str(err), test.err) {
				t.Errorf("expected err %q to contain %q", actual, test.err)
			}
			if result != test.result {
				t.Errorf("expected result %q; got %q", test.result, result)
			}
			allocs := testing.AllocsPerRun(10, func() {
				_, _ = test.template.execute(test.engine, test.data)
			})
			if test.noallocs && allocs > 0 {
				t.Errorf("expected no allocations; got %f", allocs)
			}
		})
	}
}

func str(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
