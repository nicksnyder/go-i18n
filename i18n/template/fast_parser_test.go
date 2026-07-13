package template

import (
	"testing"
)

func TestFastParser_Cacheable(t *testing.T) {
	fp := &FastParser{}
	if !fp.Cacheable() {
		t.Error("FastParser should always be cacheable")
	}
}

func TestFastParser_Parse(t *testing.T) {
	tests := []struct {
		name       string
		src        string
		leftDelim  string
		rightDelim string
		data       any
		want       string
		wantErr    bool
	}{
		{
			name: "simple substitution with dot prefix",
			src:  "Hello {{.Name}}!",
			data: map[string]interface{}{"Name": "World"},
			want: "Hello World!",
		},
		{
			name: "simple substitution without dot prefix",
			src:  "Hello {{Name}}!",
			data: map[string]interface{}{"Name": "World"},
			want: "Hello World!",
		},
		{
			name: "multiple placeholders",
			src:  "{{.Greeting}}, {{.Name}}!",
			data: map[string]interface{}{"Greeting": "Hello", "Name": "World"},
			want: "Hello, World!",
		},
		{
			name: "no placeholders",
			src:  "Hello World!",
			data: map[string]interface{}{"Name": "Ignored"},
			want: "Hello World!",
		},
		{
			name: "nil data",
			src:  "Hello World!",
			data: nil,
			want: "Hello World!",
		},
		{
			name: "missing key keeps placeholder",
			src:  "Hello {{.Name}}!",
			data: map[string]interface{}{},
			want: "Hello {{.Name}}!",
		},
		{
			name:       "custom delimiters",
			src:        "Hello <<.Name>>!",
			leftDelim:  "<<",
			rightDelim: ">>",
			data:       map[string]interface{}{"Name": "World"},
			want:       "Hello World!",
		},
		{
			name: "map[string]string data",
			src:  "Hello {{.Name}}!",
			data: map[string]string{"Name": "World"},
			want: "Hello World!",
		},
		{
			name: "integer value",
			src:  "Count: {{.Count}}",
			data: map[string]interface{}{"Count": 42},
			want: "Count: 42",
		},
		{
			name: "plural count style",
			src:  "You have {{.PluralCount}} items",
			data: map[string]interface{}{"PluralCount": 5},
			want: "You have 5 items",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fp := &FastParser{}
			parsed, err := fp.Parse(tt.src, tt.leftDelim, tt.rightDelim)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}

			got, err := parsed.Execute(tt.data)
			if err != nil {
				t.Errorf("Execute() error = %v", err)
				return
			}
			if got != tt.want {
				t.Errorf("Execute() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFastParser_ParseError(t *testing.T) {
	fp := &FastParser{}
	_, err := fp.Parse("Hello {{Name", "", "")
	if err == nil {
		t.Error("expected error for unclosed tag")
	}
}

func TestFastParser_StructData(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	fp := &FastParser{}
	parsed, err := fp.Parse("{{.Name}} is {{.Age}} years old", "", "")
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	got, err := parsed.Execute(Person{Name: "Alice", Age: 30})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	want := "Alice is 30 years old"
	if got != want {
		t.Errorf("Execute() = %q, want %q", got, want)
	}
}

func TestFastParser_PointerStructData(t *testing.T) {
	type Person struct {
		Name string
	}

	fp := &FastParser{}
	parsed, err := fp.Parse("Hello {{.Name}}", "", "")
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	got, err := parsed.Execute(&Person{Name: "Bob"})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	want := "Hello Bob"
	if got != want {
		t.Errorf("Execute() = %q, want %q", got, want)
	}
}

func TestFastParser_ParserDelims(t *testing.T) {
	fp := &FastParser{
		LeftDelim:  "[[",
		RightDelim: "]]",
	}

	parsed, err := fp.Parse("Hello [[.Name]]!", "", "")
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	got, err := parsed.Execute(map[string]interface{}{"Name": "World"})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	want := "Hello World!"
	if got != want {
		t.Errorf("Execute() = %q, want %q", got, want)
	}
}

func BenchmarkFastParser_Execute(b *testing.B) {
	fp := &FastParser{}
	parsed, _ := fp.Parse("Hello {{.Name}}, you have {{.Count}} messages", "", "")
	data := map[string]interface{}{"Name": "User", "Count": 42}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = parsed.Execute(data)
	}
}

func BenchmarkTextParser_Execute(b *testing.B) {
	tp := &TextParser{}
	parsed, _ := tp.Parse("Hello {{.Name}}, you have {{.Count}} messages", "", "")
	data := map[string]interface{}{"Name": "User", "Count": 42}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = parsed.Execute(data)
	}
}
