package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
)

type Extractor struct {
	messages []Message
}

type Message struct {
	Content     string
	Context     string
	Translation string
}

func NewExtractor() *Extractor {
	return &Extractor{messages: make([]Message, 0)}
}

func (e *Extractor) ExtractFiles(filenames []string) {
	for _, filename := range filenames {
		e.ExtractFile(filename)
	}
}

func (e *Extractor) ExtractFile(filename string) {
	fmt.Fprintln(os.Stderr, filename)
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filename, nil, 0)
	if err != nil {
		fmt.Fprint(os.Stderr, err.String())
	}
	for _, decl := range file.Decls {
		ast.Walk(e, decl)
	}
}

// Visit satisifies the ast.Visitor interface.
// See http://golang.org/pkg/go/ast/#Visitor
func (e *Extractor) Visit(n ast.Node) (v ast.Visitor) {
	if n == nil {
		return nil
	}
	switch i := n.(type) {
	case *ast.CallExpr:
		if isNewMessageCall(i.Fun) {
			content, context, ok := getNewMessageArgs(i.Args)
			if ok {
				e.messages = append(e.messages, Message{Content: content, Context: context})
			}
		}
	}
	return e
}

// Messages returns the messages extracted.
func (e *Extractor) Messages() []Message {
	return e.messages
}

func isNewMessageCall(e ast.Expr) bool {
	switch i := e.(type) {
	case *ast.SelectorExpr:
		return i.Sel.Name == "NewMessage" && isI18nIdent(i.X)
	}
	return false
}

func isI18nIdent(e ast.Expr) bool {
	switch i := e.(type) {
	case *ast.Ident:
		return i.Name == "i18n"
	}
	return false
}

func getNewMessageArgs(exprs []ast.Expr) (content, context string, ok bool) {
	if len(exprs) != 2 {
		return
	}
	if content, ok = getString(exprs[0]); !ok {
		return
	}
	context, ok = getString(exprs[1])
	return
}

func getString(e ast.Expr) (s string, ok bool) {
	switch i := e.(type) {
	case *ast.BasicLit:
		if i.Kind == token.STRING {
			return i.Value[1 : len(i.Value)-1], true
		}
	}
	return
}
