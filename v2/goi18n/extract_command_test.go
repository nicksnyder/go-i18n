package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func TestExtract(t *testing.T) {

	tests := []struct {
		name     string
		file     string
		messages []*i18n.Message
		err      error
	}{
		{
			name:     "no translations",
			file:     `package main`,
			messages: nil,
			err:      nil,
		},
		{
			name: "global declaration",
			file: `package main

			import "github.com/nicksnyder/go-i18n/v2/i18n"

			var m = &i18n.Message{
				ID: "Plural ID",
			}
			`,
			messages: []*i18n.Message{
				&i18n.Message{
					ID: "Plural ID",
				},
			},
			err: nil,
		},
		{
			name: "short form id only",
			file: `package main

			import "github.com/nicksnyder/go-i18n/v2/i18n"

			func main() {
				bundle := i18n.NewBundle()
				l := i18n.NewLocalizer(bundle, "en")
				l.Localize(&i18n.LocalizeConfig{MessageID: "Plural ID"})
			}
			`,
			messages: []*i18n.Message{
				&i18n.Message{
					ID: "Plural ID",
				},
			},
			err: nil,
		},
		{
			name: "must short form id only",
			file: `package main

			import "github.com/nicksnyder/go-i18n/v2/i18n"

			func main() {
				bundle := i18n.NewBundle()
				l := i18n.NewLocalizer(bundle, "en")
				l.MustLocalize(&i18n.LocalizeConfig{MessageID: "Plural ID"})
			}
			`,
			messages: []*i18n.Message{
				&i18n.Message{
					ID: "Plural ID",
				},
			},
			err: nil,
		},
		{
			name: "custom package name",
			file: `package main

			import bar "github.com/nicksnyder/go-i18n/v2/i18n"

			func main() {
				_ := &bar.Message{
					ID:          "Plural ID",
				}
			}
			`,
			messages: []*i18n.Message{
				&i18n.Message{
					ID: "Plural ID",
				},
			},
			err: nil,
		},
		{
			name: "exhaustive plural translation",
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
			messages: []*i18n.Message{
				&i18n.Message{
					ID:          "Plural ID",
					Description: "Plural description",
					Zero:        "Zero translation",
					One:         "One translation",
					Two:         "Two translation",
					Few:         "Few translation",
					Many:        "Many translation",
					Other:       "Other translation",
				},
			},
			err: nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actualMessages, err := extractMessages([]byte(test.file))
			if err != test.err {
				t.Errorf("expected error: %q\n     got error: %q", test.err, err)
			}
			if !reflect.DeepEqual(actualMessages, test.messages) {
				t.Errorf("file:\n%s\nexpected: %s\n     got: %s", test.file, marshalTest(test.messages), marshalTest(actualMessages))
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
	if err := testableMain([]string{"goi18n", "extract", "-outdir", outdir, "../example/"}); err != nil {
		t.Fatal(err)
	}
	actual, err := ioutil.ReadFile(filepath.Join(outdir, "active.en.toml"))
	if err != nil {
		t.Fatal(err)
	}
	expected := []byte(`HelloPerson = "Hello {{.Name}}"

[ColoredCats]
description = "The number of cats a person has"
one = "I have {{.Count}} {{.Color}} cat."
other = "I have {{.Count}} {{.Color}} cats."

[UnreadEmails]
description = "The number of unread emails a person has"
one = "I have {{.PluralCount}} unread email."
other = "I have {{.PluralCount}} unread emails."
`)
	if !bytes.Equal(actual, expected) {
		t.Fatalf("files not equal\nactual:\n%s\nexpected:\n%s", actual, expected)
	}
}

func marshalTest(value interface{}) string {
	buf, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	return string(buf)
}
