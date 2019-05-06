package translation

import (
	"bytes"
	"fmt"

	"testing"
	gotemplate "text/template"
)

func TestNilTemplate(t *testing.T) {
	expected := "hello"
	tmpl := &template{
		tmpl: nil,
		src:  expected,
	}
	if actual := tmpl.Execute(nil); actual != expected {
		t.Errorf("Execute(nil) returned %s; expected %s", actual, expected)
	}
}

func TestMarshalText(t *testing.T) {
	tmpl := &template{
		tmpl: gotemplate.Must(gotemplate.New("id").Parse("this is a {{.foo}} template")),
		src:  "boom",
	}
	expectedBuf := []byte(tmpl.src)
	if buf, err := tmpl.MarshalText(); !bytes.Equal(buf, expectedBuf) || err != nil {
		t.Errorf("MarshalText() returned %#v, %#v; expected %#v, nil", buf, err, expectedBuf)
	}
}

func TestUnmarshalText(t *testing.T) {
	tmpl := &template{}
	if err := tmpl.UnmarshalText([]byte("hello {{.World}}")); err != nil {
		t.Fatal(err)
	}
	result := tmpl.Execute(map[string]string{
		"World": "world!",
	})
	expected := "hello world!"
	if result != expected {
		t.Errorf("expected %#v; got %#v", expected, result)
	}
}

func BenchmarkExecuteNilTemplate(b *testing.B) {
	template := &template{src: "hello world"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		template.Execute(nil)
	}
}

func BenchmarkExecuteHelloWorldTemplate(b *testing.B) {
	template, err := newTemplate("hello world")
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		template.Execute(nil)
	}
}

// Executing a simple template like this is ~6x slower than Sprintf
// but it is still only a few microseconds which should be sufficiently fast.
// The benefit is that we have nice semantic tags in the translation.
func BenchmarkExecuteHelloNameTemplate(b *testing.B) {
	template, err := newTemplate("hello {{.Name}}")
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		template.Execute(map[string]string{
			"Name": "Nick",
		})
	}
}

var sprintfResult string

func BenchmarkSprintf(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sprintfResult = fmt.Sprintf("hello %s", "nick")
	}
}
