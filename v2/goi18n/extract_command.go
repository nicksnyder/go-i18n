package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func usageExtract() {
	fmt.Fprintf(os.Stderr, `usage: goi18n extract [options] [paths]

Extract walks the files and directories in paths and extracts all messages to a single file.
If no files or paths are provided, it walks the current working directory.

	xx-yy.active.format
		This file contains messages that should be loaded at runtime.

Flags:

	-sourceLanguage tag
		The language tag of the extracted messages (e.g. en, en-US, zh-Hant-CN).
		Default: en

	-outdir directory
		Write message files to this directory.
		Default: .

	-format format
		Output message files in this format.
		Supported formats: json, toml, yaml
		Default: toml
`)
}

type extractCommand struct {
	paths          []string
	sourceLanguage languageTag
	outdir         string
	format         string
}

func (ec *extractCommand) name() string {
	return "extract"
}

func (ec *extractCommand) parse(args []string) error {
	flags := flag.NewFlagSet("extract", flag.ExitOnError)
	flags.Usage = usageExtract

	flags.Var(&ec.sourceLanguage, "sourceLanguage", "en")
	flags.StringVar(&ec.outdir, "outdir", ".", "")
	flags.StringVar(&ec.format, "format", "toml", "")
	if err := flags.Parse(args); err != nil {
		return err
	}

	ec.paths = flags.Args()
	return nil
}

func (ec *extractCommand) execute() error {
	if len(ec.paths) == 0 {
		ec.paths = []string{"."}
	}
	messages := []*i18n.Message{}
	for _, path := range ec.paths {
		if err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			if filepath.Ext(path) != ".go" {
				return nil
			}

			// Don't extract from test files.
			if strings.HasSuffix(path, "_test.go") {
				return nil
			}

			buf, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			msgs, err := extractMessages(buf)
			if err != nil {
				return err
			}
			messages = append(messages, msgs...)
			return nil
		}); err != nil {
			return err
		}
	}
	messageTemplates := map[string]*i18n.MessageTemplate{}
	for _, m := range messages {
		if mt := i18n.NewMessageTemplate(m); mt != nil {
			messageTemplates[m.ID] = mt
		}
	}
	path, content, err := writeFile(ec.outdir, "active", ec.sourceLanguage.Tag(), ec.format, messageTemplates, true)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, content, 0666)
}

// extractMessages extracts messages from the bytes of a Go source file.
func extractMessages(buf []byte) ([]*i18n.Message, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", buf, parser.AllErrors)
	if err != nil {
		return nil, err
	}
	extractor := newExtractor(file)
	ast.Walk(extractor, file)
	return extractor.messages, nil
}

func newExtractor(file *ast.File) *extractor {
	return &extractor{i18nPackageName: i18nPackageName(file)}
}

type extractor struct {
	i18nPackageName string
	messages        []*i18n.Message
}

func (e *extractor) Visit(node ast.Node) ast.Visitor {
	e.extractMessages(node)
	return e
}

func (e *extractor) extractMessages(node ast.Node) {
	cl, ok := node.(*ast.CompositeLit)
	if !ok {
		return
	}
	switch t := cl.Type.(type) {
	case *ast.SelectorExpr:
		if !e.isMessageType(t) {
			return
		}
		e.extractMessage(cl)
	case *ast.ArrayType:
		if !e.isMessageType(t.Elt) {
			return
		}
		for _, el := range cl.Elts {
			ecl, ok := el.(*ast.CompositeLit)
			if !ok {
				continue
			}
			e.extractMessage(ecl)
		}
	case *ast.MapType:
		if !e.isMessageType(t.Value) {
			return
		}
		for _, el := range cl.Elts {
			kve, ok := el.(*ast.KeyValueExpr)
			if !ok {
				continue
			}
			vcl, ok := kve.Value.(*ast.CompositeLit)
			if !ok {
				continue
			}
			e.extractMessage(vcl)
		}
	}
}

func (e *extractor) isMessageType(expr ast.Expr) bool {
	se := unwrapSelectorExpr(expr)
	if se == nil {
		return false
	}
	if se.Sel.Name != "Message" && se.Sel.Name != "LocalizeConfig" {
		return false
	}
	x, ok := se.X.(*ast.Ident)
	if !ok {
		return false
	}
	return x.Name == e.i18nPackageName
}

func unwrapSelectorExpr(e ast.Expr) *ast.SelectorExpr {
	switch et := e.(type) {
	case *ast.SelectorExpr:
		return et
	case *ast.StarExpr:
		se, _ := et.X.(*ast.SelectorExpr)
		return se
	default:
		return nil
	}
}

func (e *extractor) extractMessage(cl *ast.CompositeLit) {
	data := make(map[string]string)
	for _, elt := range cl.Elts {
		kve, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		key, ok := kve.Key.(*ast.Ident)
		if !ok {
			continue
		}
		v, ok := extractStringLiteral(kve.Value)
		if !ok {
			continue
		}
		data[key.Name] = v
	}
	if len(data) == 0 {
		return
	}
	if messageID := data["MessageID"]; messageID != "" {
		data["ID"] = messageID
	}
	e.messages = append(e.messages, i18n.MustNewMessage(data))
}

func extractStringLiteral(expr ast.Expr) (string, bool) {
	switch v := expr.(type) {
	case *ast.BasicLit:
		if v.Kind != token.STRING {
			return "", false
		}
		s := v.Value[1 : len(v.Value)-1]
		if v.Value[0] == '"' {
			s = strings.Replace(s, `\"`, `"`, -1)
		}
		return s, true
	case *ast.BinaryExpr:
		if v.Op != token.ADD {
			return "", false
		}
		x, ok := extractStringLiteral(v.X)
		if !ok {
			return "", false
		}
		y, ok := extractStringLiteral(v.Y)
		if !ok {
			return "", false
		}
		return x + y, true
	case *ast.Ident:
		if v.Obj == nil {
			return "", false
		}
		switch z := v.Obj.Decl.(type) {
		case *ast.ValueSpec:
			if len(z.Values) == 0 {
				return "", false
			}
			s, ok := extractStringLiteral(z.Values[0])
			if !ok {
				return "", false
			}
			return s, true
		}
		return "", false
	default:
		return "", false
	}
}

func i18nPackageName(file *ast.File) string {
	for _, i := range file.Imports {
		if i.Path.Kind == token.STRING && i.Path.Value == `"github.com/nicksnyder/go-i18n/v2/i18n"` {
			if i.Name == nil {
				return "i18n"
			}
			return i.Name.Name
		}
	}
	return ""
}
