package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestExtract(t *testing.T) {

	tests := []struct {
		name       string
		fileName   string
		file       string
		activeFile []byte
	}{
		{
			name:     "no translations",
			fileName: "file.go",
			file:     `package main`,
		},
		{
			name:     "global declaration",
			fileName: "file.go",
			file: `package main

			import "github.com/nicksnyder/go-i18n/v2/i18n"

			var m = &i18n.Message{
				ID: "Plural ID",
			}
			`,
		},
		{
			name:     "escape",
			fileName: "file.go",
			file: `package main

			import "github.com/nicksnyder/go-i18n/v2/i18n"

			var a = &i18n.Message{
				ID:    "a",
				Other: "a \" b",
			}
			var b = &i18n.Message{
				ID:    "b",
				Other: ` + "`" + `a " b` + "`" + `,
			}
			`,
			activeFile: []byte(`a = "a \" b"
b = "a \" b"
`),
		},
		{
			name:     "array",
			fileName: "file.go",
			file: `package main

			import "github.com/nicksnyder/go-i18n/v2/i18n"

			var a = []*i18n.Message{
				{
					ID:    "a",
					Other: "a",
				},
				{
					ID:    "b",
					Other: "b",
				},
			}
			`,
			activeFile: []byte(`a = "a"
b = "b"
`),
		},
		{
			name:     "map",
			fileName: "file.go",
			file: `package main

			import "github.com/nicksnyder/go-i18n/v2/i18n"

			var a = map[string]*i18n.Message{
				"a": {
					ID:    "a",
					Other: "a",
				},
				"b": {
					ID:    "b",
					Other: "b",
				},
			}
			`,
			activeFile: []byte(`a = "a"
b = "b"
`),
		},
		{
			name:     "no extract from test",
			fileName: "file_test.go",
			file: `package main

			import "github.com/nicksnyder/go-i18n/v2/i18n"

			func main() {
				bundle := i18n.NewBundle(language.English)
				l := i18n.NewLocalizer(bundle, "en")
				l.Localize(&i18n.LocalizeConfig{MessageID: "Plural ID"})
			}
			`,
		},
		{
			name:     "must short form id only",
			fileName: "file.go",
			file: `package main

			import "github.com/nicksnyder/go-i18n/v2/i18n"

			func main() {
				bundle := i18n.NewBundle(language.English)
				l := i18n.NewLocalizer(bundle, "en")
				l.MustLocalize(&i18n.LocalizeConfig{MessageID: "Plural ID"})
			}
			`,
		},
		{
			name:     "custom package name",
			fileName: "file.go",
			file: `package main

			import bar "github.com/nicksnyder/go-i18n/v2/i18n"

			func main() {
				_ := &bar.Message{
					ID:          "Plural ID",
				}
			}
			`,
		},
		{
			name:     "exhaustive plural translation",
			fileName: "file.go",
			file: `package main

			import "github.com/nicksnyder/go-i18n/v2/i18n"

			func main() {
				_ := &i18n.Message{
					ID:          "Plural ID",
					Description: "Plural description",
					Zero:        "Zero translation",
					One:         "One translation",
					Two:         "Two translation",
					Few:         "Few translation",
					Many:        "Many translation",
					Other:       "Other translation",
				}
			}
			`,
			activeFile: []byte(`["Plural ID"]
description = "Plural description"
few = "Few translation"
many = "Many translation"
one = "One translation"
other = "Other translation"
two = "Two translation"
zero = "Zero translation"
`),
		},
		{
			name:     "concat id",
			fileName: "file.go",
			file: `package main

			import "github.com/nicksnyder/go-i18n/v2/i18n"

			func main() {
				_ := &i18n.Message{
					ID: "Plural" +
						" " +
						"ID",
				}
			}
			`,
		},
		{
			name:     "global declaration",
			fileName: "file.go",
			file: `package main

			import "github.com/nicksnyder/go-i18n/v2/i18n"

			const constID = "ConstantID"
			
			var m = &i18n.Message{
				ID: constID,
				Other: "ID is a constant",
			}
			`,
			activeFile: []byte(`ConstantID = "ID is a constant"
`),
		},
		{
			name:     "undefined identifier in composite lit",
			fileName: "file.go",
			file: `package main

			import "github.com/nicksnyder/go-i18n/v2/i18n"

			var m = &i18n.LocalizeConfig{
				Funcs: Funcs,
			}
			`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			indir := mustTempDir("TestExtractCommandIn")
			defer os.RemoveAll(indir)
			outdir := mustTempDir("TestExtractCommandOut")
			defer os.RemoveAll(outdir)

			inpath := filepath.Join(indir, test.fileName)
			if err := ioutil.WriteFile(inpath, []byte(test.file), 0666); err != nil {
				t.Fatal(err)
			}

			if code := testableMain([]string{"extract", "-outdir", outdir, indir}); code != 0 {
				t.Fatalf("expected exit code 0; got %d\n", code)
			}

			files, err := ioutil.ReadDir(outdir)
			if err != nil {
				t.Fatal(err)
			}
			if len(files) != 1 {
				t.Fatalf("expected 1 file; got %#v", files)
			}
			actualFile := files[0]
			expectedName := "active.en.toml"
			if actualFile.Name() != expectedName {
				t.Fatalf("expected %s; got %s", expectedName, actualFile.Name())
			}

			outpath := filepath.Join(outdir, actualFile.Name())
			actual, err := ioutil.ReadFile(outpath)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(actual, test.activeFile) {
				t.Fatalf("\nexpected:\n%s\n\ngot:\n%s", test.activeFile, actual)
			}
		})
	}
}

func TestExtractCommand(t *testing.T) {
	outdir, err := ioutil.TempDir("", "TestExtractCommand")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(outdir)
	if code := testableMain([]string{"extract", "-outdir", outdir, "../example/"}); code != 0 {
		t.Fatalf("expected exit code 0; got %d", code)
	}
	actual, err := ioutil.ReadFile(filepath.Join(outdir, "active.en.toml"))
	if err != nil {
		t.Fatal(err)
	}
	expected := []byte(`HelloPerson = "Hello {{.Name}}"

[MyUnreadEmails]
description = "The number of unread emails I have"
one = "I have {{.PluralCount}} unread email."
other = "I have {{.PluralCount}} unread emails."

[PersonUnreadEmails]
description = "The number of unread emails a person has"
one = "{{.Name}} has {{.UnreadEmailCount}} unread email."
other = "{{.Name}} has {{.UnreadEmailCount}} unread emails."
`)
	if !bytes.Equal(actual, expected) {
		t.Fatalf("files not equal\nactual:\n%s\nexpected:\n%s", actual, expected)
	}
}
