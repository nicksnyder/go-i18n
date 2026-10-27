package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestExtract(t *testing.T) {
	tests := []struct {
		name             string
		fileName         string
		file             string
		activeFile       []byte
		expectedExitCode int
		expectedErr      error
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
			name:     "duplicate ID with same message is not an error",
			fileName: "file.go",
			file: `package main

			import "github.com/nicksnyder/go-i18n/v2/i18n"

			var m1 = &i18n.Message{
				ID: "m",
				Other: "m",
			}
			var m2 = &i18n.Message{
				ID: "m",
				Other: "m",
			}
			`,
			activeFile: []byte(`m = "m"
`),
		},
		{
			name:     "duplicate ID with different message is an error",
			fileName: "file.go",
			file: `package main

			import "github.com/nicksnyder/go-i18n/v2/i18n"

			var m1 = &i18n.Message{
				ID: "m",
				Other: "m1",
			}
			var m2 = &i18n.Message{
				ID: "m",
				Other: "m2",
			}
			`,
			expectedExitCode: 1,
		},
		{
			name:     "escape newline",
			fileName: "file.go",
			file: `package main

			import "github.com/nicksnyder/go-i18n/v2/i18n"

			var hasnewline = &i18n.Message{
				ID:    "hasnewline",
				Other: "\nfoo\nbar\\",
			}
			`,
			activeFile: []byte(`hasnewline = "\nfoo\nbar\\"
`),
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
		{
			name:     "casted const",
			fileName: "file.go",
			file: `package main

			import "github.com/nicksnyder/go-i18n/v2/i18n"

			type ConstType string

			const Const ConstType = "my const"

			var m = &i18n.LocalizeConfig{
				ID: "id",
				Other: string(Const),
			}
			`,
			activeFile: []byte(`id = "my const"
`),
		},
		{
			name:     "const with imported type",
			fileName: "file.go",
			file: `package main

			import (
				"github.com/nicksnyder/go-i18n/v2/i18n"
				"example.com/keys"
			)

			const base keys.ID = "imported"
			const id keys.ID = base + "ID"

			func f() {
				const other keys.Text = "local"
				_ = &i18n.Message{
					ID: string(id),
					Other: string(other),
				}
			}
			`,
			activeFile: []byte(`importedID = "local"
`),
		},
		{
			name:     "var declarations",
			fileName: "file.go",
			file: `package main

			import (
				"github.com/nicksnyder/go-i18n/v2/i18n"
				"example.com/keys"
			)

			var base keys.ID = "var"

			func f() {
				var id = string(base) + "ID"
				_ = &i18n.Message{
					ID: id,
					Other: "from " + id,
				}
			}
			`,
			activeFile: []byte(`varID = "from varID"
`),
		},
		{
			name:     "string conversion without arguments",
			fileName: "file.go",
			file: `package main

			import "github.com/nicksnyder/go-i18n/v2/i18n"

			var m = &i18n.Message{
				ID: string(),
			}
			`,
		},
		{
			name:     "self-referencing declarations",
			fileName: "file.go",
			file: `package main

			import "github.com/nicksnyder/go-i18n/v2/i18n"

			var a = a
			const b, c = c, b

			var m = &i18n.Message{
				ID: "id",
				Description: a + b,
				Other: "hello",
			}
			`,
			activeFile: []byte(`id = "hello"
`),
		},
		{
			name:     "multi-name declarations",
			fileName: "file.go",
			file: `package main

			import "github.com/nicksnyder/go-i18n/v2/i18n"

			const a, b = "x", "y"
			const c, d = "z", c

			var m = &i18n.Message{
				ID: "id",
				Description: d,
				Other: b,
			}
			`,
			activeFile: []byte(`[id]
description = "z"
other = "y"
`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			indir := mustTempDir("TestExtractCommandIn")
			defer mustRemoveAll(t, indir)
			outdir := mustTempDir("TestExtractCommandOut")
			defer mustRemoveAll(t, outdir)

			inpath := filepath.Join(indir, test.fileName)
			if err := os.WriteFile(inpath, []byte(test.file), 0666); err != nil {
				t.Fatal(err)
			}

			code := testableMain([]string{"extract", "-outdir", outdir, indir})
			if code != test.expectedExitCode {
				t.Fatalf("expected exit code %d; got %d\n", test.expectedExitCode, code)
			}

			files, err := os.ReadDir(outdir)
			if err != nil {
				t.Fatal(err)
			}

			if code != 0 {
				if len(files) != 0 {
					t.Fatalf("expected 0 files; got %#v", files)
				}
				return
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
			actual, err := os.ReadFile(outpath)
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
	outdir, err := os.MkdirTemp("", "TestExtractCommand")
	if err != nil {
		t.Fatal(err)
	}
	defer mustRemoveAll(t, outdir)
	if code := testableMain([]string{"extract", "-outdir", outdir, "../example/"}); code != 0 {
		t.Fatalf("expected exit code 0; got %d", code)
	}
	actual, err := os.ReadFile(filepath.Join(outdir, "active.en.toml"))
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

func TestExtractCommandConstantsInOtherFiles(t *testing.T) {
	indir := mustTempDir("TestExtractCommandConstantsIn")
	defer mustRemoveAll(t, indir)
	outdir := mustTempDir("TestExtractCommandConstantsOut")
	defer mustRemoveAll(t, outdir)

	// Two directories declare the same package name and constant name with
	// different values; each must resolve to its own package's constant.
	files := map[string]string{
		"a/ids.go": `package msgs
const id = "a" + suffix
const suffix = "ID"
`,
		"a/msgs.go": `package msgs
import "github.com/nicksnyder/go-i18n/v2/i18n"
var m = &i18n.Message{ID: id, Other: text}
`,
		"a/text.go": `package msgs
const text = "A text"
`,
		"b/ids.go": `package msgs
const id = "bID"
const text = "B text"
`,
		"b/msgs.go": `package msgs
import "github.com/nicksnyder/go-i18n/v2/i18n"
var m = &i18n.Message{ID: id, Other: text}
`,
		// Platform-specific files may declare the same name; each file must
		// use its own declaration.
		"c/c_linux.go": `package msgs
import "github.com/nicksnyder/go-i18n/v2/i18n"
const text = "linux text"
var m = &i18n.Message{ID: "cLinux", Other: text}
`,
		"c/c_windows.go": `package msgs
import "github.com/nicksnyder/go-i18n/v2/i18n"
const text = "windows text"
var m = &i18n.Message{ID: "cWindows", Other: text}
`,
	}
	for name, content := range files {
		path := filepath.Join(indir, name)
		if err := os.MkdirAll(filepath.Dir(path), 0777); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0666); err != nil {
			t.Fatal(err)
		}
	}

	if code := testableMain([]string{"extract", "-outdir", outdir, indir}); code != 0 {
		t.Fatalf("expected exit code 0; got %d", code)
	}
	actual, err := os.ReadFile(filepath.Join(outdir, "active.en.toml"))
	if err != nil {
		t.Fatal(err)
	}
	expected := []byte(`aID = "A text"
bID = "B text"
cLinux = "linux text"
cWindows = "windows text"
`)
	if !bytes.Equal(actual, expected) {
		t.Fatalf("files not equal\nactual:\n%s\nexpected:\n%s", actual, expected)
	}
}
