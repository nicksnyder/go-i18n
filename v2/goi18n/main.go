// Command goi18n manages message files used by the i18n package.
//
//     go get -u github.com/nicksnyder/go-i18n/v2/goi18n
//     goi18n -help
//
// Use `goi18n extract` to create a message file that contains the messages defined in your Go source files.
//     # en.toml
//     [PersonCats]
//     description = "The number of cats a person has"
//     one = "{{.Name}} has {{.Count}} cat."
//     other = "{{.Name}} has {{.Count}} cats."
//
// Use `goi18n merge` to create message files for translation.
//     # translate.es.toml
//     [PersonCats]
//     description = "The number of cats a person has"
//     hash = "sha1-f937a0e05e19bfe6cd70937c980eaf1f9832f091"
//     one = "{{.Name}} has {{.Count}} cat."
//     other = "{{.Name}} has {{.Count}} cats."
//
// Use `goi18n merge` to merge translated message files with your existing message files.
//     # active.es.toml
//     [PersonCats]
//     description = "The number of cats a person has"
//     hash = "sha1-f937a0e05e19bfe6cd70937c980eaf1f9832f091"
//     one = "{{.Name}} tiene {{.Count}} gato."
//     other = "{{.Name}} tiene {{.Count}} gatos."
//
// Load the active messages into your bundle.
//     bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
//     bundle.MustLoadMessageFile("active.es.toml")
package main

import (
	"flag"
	"fmt"
	"os"

	"golang.org/x/text/language"
)

func mainUsage() {
	fmt.Fprintf(os.Stderr, `goi18n (v2) is a tool for managing message translations.

Usage:

	goi18n command [arguments]

The commands are:

	merge		merge message files
	extract		extract messages from Go files

Workflow:

	Use 'goi18n extract' to create a message file that contains the messages defined in your Go source files.

		# en.toml
		[PersonCats]
		description = "The number of cats a person has"
		one = "{{.Name}} has {{.Count}} cat."
		other = "{{.Name}} has {{.Count}} cats."

	Use 'goi18n merge' to create message files for translation.

		# translate.es.toml
		[PersonCats]
		description = "The number of cats a person has"
		hash = "sha1-f937a0e05e19bfe6cd70937c980eaf1f9832f091"
		one = "{{.Name}} has {{.Count}} cat."
		other = "{{.Name}} has {{.Count}} cats."

	Use 'goi18n merge' to merge translated message files with your existing message files.

		# active.es.toml
		[PersonCats]
		description = "The number of cats a person has"
		hash = "sha1-f937a0e05e19bfe6cd70937c980eaf1f9832f091"
		one = "{{.Name}} tiene {{.Count}} gato."
		other = "{{.Name}} tiene {{.Count}} gatos."

	Load the active messages into your bundle.

		bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
		bundle.MustLoadMessageFile("active.es.toml")
`)
	os.Exit(2)
}

type command interface {
	name() string
	parse(arguments []string)
	execute() error
}

func main() {
	if err := testableMain(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func testableMain(args []string) error {
	flags := flag.NewFlagSet("goi18n", flag.ExitOnError)
	flags.Usage = mainUsage
	flags.Parse(args)
	if len(args) == 1 {
		mainUsage()
	}
	commands := []command{
		&mergeCommand{},
		&extractCommand{},
	}
	cmdName := args[1]
	for _, cmd := range commands {
		if cmd.name() == cmdName {
			cmd.parse(args[2:])
			return cmd.execute()
		}
	}
	return fmt.Errorf("goi18n: unknown subcommand %s", cmdName)
}

type languageTag language.Tag

func (lt languageTag) String() string {
	return lt.Tag().String()
}

func (lt *languageTag) Set(value string) error {
	t, err := language.Parse(value)
	if err != nil {
		return err
	}
	*lt = languageTag(t)
	return nil
}

func (lt languageTag) Tag() language.Tag {
	tag := language.Tag(lt)
	if tag.IsRoot() {
		return language.English
	}
	return tag
}
