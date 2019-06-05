package internal

import (
	"testing"
	"text/template"
)

func TestExecute(t *testing.T) {
	tests := []struct {
		template *Template
		funcs    template.FuncMap
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
			funcs: template.FuncMap{
				"world": func() string {
					return "world"
				},
			},
			result: "hello world",
		},
		{
			template: &Template{
				Src: "hello {{",
			},
			err:      "template: :1: unexpected unclosed action in command",
			noallocs: true,
		},
	}

	for _, test := range tests {
		t.Run(test.template.Src, func(t *testing.T) {
			result, err := test.template.Execute(test.funcs, test.data)
			if actual := str(err); actual != test.err {
				t.Errorf("expected err %q; got %q", test.err, actual)
			}
			if result != test.result {
				t.Errorf("expected result %q; got %q", test.result, result)
			}
			allocs := testing.AllocsPerRun(10, func() {
				_, _ = test.template.Execute(test.funcs, test.data)
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
