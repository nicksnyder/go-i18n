package goio

import (
	"github.com/nicksnyder/go-i18n/src/pkg/msg"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
)

const defaultFilename = "unknown"

type Reader struct {

}

func NewReader() msg.Reader {
	return &Reader{}
}

func (r *Reader) ReadMessages(ir io.Reader) ([]msg.Message, os.Error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, defaultFilename, ir, 0)
	if err != nil {
		return nil, err
	}

	v := &visitor{msgs: make([]msg.Message, 0)}
	for _, decl := range file.Decls {
		ast.Walk(v, decl)
	}

	return v.msgs, nil
}

type visitor struct {
	msgs []msg.Message
}

// Visit satisifies the ast.Visitor interface.
// See http://golang.org/pkg/go/ast/#Visitor
func (v *visitor) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}
	switch i := n.(type) {
	case *ast.CallExpr:
		if isNewMessageCall(i.Fun) {
			content, context, ok := getNewMessageArgs(i.Args)
			if ok {
				id := msg.Id(context, content)
				v.msgs = append(v.msgs, msg.Message{Id: id, Context: context, Content: content})
			}
		}
	}
	return v
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
