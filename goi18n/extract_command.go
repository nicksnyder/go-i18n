package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/constant"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"reflect"
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
	filesByDir := map[string][][]byte{}
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

			buf, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			dir := filepath.Dir(path)
			filesByDir[dir] = append(filesByDir[dir], buf)
			return nil
		}); err != nil {
			return err
		}
	}
	messages := []*i18n.Message{}
	for _, files := range filesByDir {
		msgs, err := extractMessages(files)
		if err != nil {
			return err
		}
		messages = append(messages, msgs...)
	}
	messageTemplates := map[string]*i18n.MessageTemplate{}
	for _, m := range messages {
		if mt := i18n.NewMessageTemplate(m); mt != nil {
			if duplicateMessage, ok := messageTemplates[m.ID]; ok && !reflect.DeepEqual(mt, duplicateMessage) {
				return &duplicateMessageIDErr{messageID: m.ID}
			}
			messageTemplates[m.ID] = mt
		}
	}
	path, content, err := writeFile(ec.outdir, "active", ec.sourceLanguage.Tag(), ec.format, messageTemplates, true)
	if err != nil {
		return err
	}
	return os.WriteFile(path, content, 0666)
}

type duplicateMessageIDErr struct {
	messageID string
}

func (e *duplicateMessageIDErr) Error() string {
	return fmt.Sprintf("duplicate message ID: %s", e.messageID)
}

// extractMessages extracts messages from the bytes of the Go source files in one directory.
func extractMessages(bufs [][]byte) ([]*i18n.Message, error) {
	fset := token.NewFileSet()
	filesByPackage := map[string][]*ast.File{}
	for _, buf := range bufs {
		file, err := parser.ParseFile(fset, "", buf, parser.AllErrors|parser.SkipObjectResolution)
		if err != nil {
			return nil, err
		}
		// Imports are not loaded, so a constant declared with an imported type
		// would have an invalid type and no value. String values do not depend
		// on the declared type, so drop it.
		ast.Inspect(file, func(n ast.Node) bool {
			if d, ok := n.(*ast.GenDecl); ok && d.Tok == token.CONST {
				for _, spec := range d.Specs {
					spec.(*ast.ValueSpec).Type = nil
				}
			}
			return true
		})
		filesByPackage[file.Name.Name] = append(filesByPackage[file.Name.Name], file)
	}
	var messages []*i18n.Message
	for _, files := range filesByPackage {
		batches := [][]*ast.File{files}
		if redeclaresName(files) {
			// Files that declare the same name are usually built for different
			// platforms, so checking them together would resolve the wrong value.
			batches = nil
			for _, file := range files {
				batches = append(batches, []*ast.File{file})
			}
		}
		for _, batch := range batches {
			// Type-check only to compute constant values. Imports are not loaded,
			// so errors about imported identifiers are expected and ignored.
			info := &types.Info{
				Types: map[ast.Expr]types.TypeAndValue{},
				Defs:  map[*ast.Ident]types.Object{},
				Uses:  map[*ast.Ident]types.Object{},
			}
			conf := types.Config{Error: func(error) {}}
			_, _ = conf.Check("", fset, batch, info)
			varInits := map[*types.Var]ast.Expr{}
			for _, file := range batch {
				ast.Inspect(file, func(n ast.Node) bool {
					if s, ok := n.(*ast.ValueSpec); ok && len(s.Values) == len(s.Names) {
						for i, name := range s.Names {
							if v, ok := info.Defs[name].(*types.Var); ok {
								varInits[v] = s.Values[i]
							}
						}
					}
					return true
				})
			}
			for _, file := range batch {
				extractor := &extractor{i18nPackageName: i18nPackageName(file), info: info, varInits: varInits}
				ast.Walk(extractor, file)
				messages = append(messages, extractor.messages...)
			}
		}
	}
	return messages, nil
}

// redeclaresName reports whether a package-level name is declared in more than one of files.
func redeclaresName(files []*ast.File) bool {
	declaredIn := map[string]*ast.File{}
	declare := func(file *ast.File, name *ast.Ident) bool {
		if name.Name == "_" {
			return false
		}
		if f, ok := declaredIn[name.Name]; ok && f != file {
			return true
		}
		declaredIn[name.Name] = file
		return false
	}
	for _, file := range files {
		for _, decl := range file.Decls {
			switch d := decl.(type) {
			case *ast.FuncDecl:
				if d.Recv == nil && d.Name.Name != "init" && declare(file, d.Name) {
					return true
				}
			case *ast.GenDecl:
				for _, spec := range d.Specs {
					switch s := spec.(type) {
					case *ast.ValueSpec:
						for _, name := range s.Names {
							if declare(file, name) {
								return true
							}
						}
					case *ast.TypeSpec:
						if declare(file, s.Name) {
							return true
						}
					}
				}
			}
		}
	}
	return false
}

type extractor struct {
	i18nPackageName string
	info            *types.Info
	varInits        map[*types.Var]ast.Expr
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
		v, ok := e.stringValue(kve.Value, nil)
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

// stringValue returns the string value of expr if it is a constant string
// expression, or is built from vars whose initializers are.
// visiting holds the vars being resolved so that cycles terminate.
func (e *extractor) stringValue(expr ast.Expr, visiting map[*types.Var]bool) (string, bool) {
	if v := e.info.Types[expr].Value; v != nil {
		if v.Kind() != constant.String {
			return "", false
		}
		return constant.StringVal(v), true
	}
	switch x := expr.(type) {
	case *ast.Ident:
		v, ok := e.info.Uses[x].(*types.Var)
		if !ok || visiting[v] || e.varInits[v] == nil {
			return "", false
		}
		if visiting == nil {
			visiting = map[*types.Var]bool{}
		}
		visiting[v] = true
		defer delete(visiting, v)
		return e.stringValue(e.varInits[v], visiting)
	case *ast.ParenExpr:
		return e.stringValue(x.X, visiting)
	case *ast.BinaryExpr:
		if x.Op != token.ADD {
			return "", false
		}
		l, ok := e.stringValue(x.X, visiting)
		if !ok {
			return "", false
		}
		r, ok := e.stringValue(x.Y, visiting)
		return l + r, ok
	case *ast.CallExpr:
		if fun, ok := x.Fun.(*ast.Ident); ok && fun.Name == "string" && len(x.Args) == 1 {
			return e.stringValue(x.Args[0], visiting)
		}
	}
	return "", false
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
