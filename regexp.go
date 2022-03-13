package main

import (
	"fmt"
	"github.com/tdewolff/parse/v2/js"
	"strings"
)

type convertRegexp struct {
	parent js.INode
	scope  *js.Scope
	block  *js.BlockStmt
}

func ConvertRegexp(ast *js.AST) {
	js.Walk(convertRegexp{}, ast)
}

func (c convertRegexp) Enter(n js.INode) js.IVisitor {
	switch v := n.(type) {
	case *js.LiteralExpr:
		if v.TokenType == js.RegExpToken {
			_=v
			re := toRegexp(v)
			switch vv := c.parent.(type) {
			case *js.Arg:
				vv.Value = re.(js.IExpr)
			case *js.BindingElement:
				vv.Default = re.(js.IExpr)
			case *js.DotExpr:
				vv.X = re.(js.IExpr)
			case *js.GroupExpr:
				vv.X = re.(js.IExpr)
			case *js.Property:
				vv.Value = re.(js.IExpr)
			default:
				panic(fmt.Sprintf("unknown parent %T", c.parent))
			}
	}
	}
	if blk, ok := n.(*js.BlockStmt); ok {
		c.block = blk
		c.scope = &blk.Scope
	}
	c.parent = n
	return c
}

func (c convertRegexp) Exit(n js.INode) {}

func toRegexp(l *js.LiteralExpr) (res js.INode) {
	data := string(l.Data)
	data = strings.TrimPrefix(data, "/")
	i := strings.LastIndex(data, "/")
	p := data[:i]
	p = strings.ReplaceAll(p, `\`, `\\`)
	p = strings.ReplaceAll(p, `"`, `\"`)
	fl := data[i+1:]
	args := []js.Arg{
		js.Arg{
			Value: js.LiteralExpr{
				TokenType: js.StringToken,
				Data: []byte(`"` + p + `"`),
			},
		},
	}
	if len(fl) > 0 {
		flags := js.Arg{
			Value: js.LiteralExpr{
				TokenType: js.StringToken,
				Data: []byte(`"` + fl + `"`),
			},
		}
		args = append(args, flags)
	}
	res = js.CallExpr{
		X: js.LiteralExpr{
			Data: []byte("RegExp"),
		},
		Args: js.Args{
			List: args,
		},
	}
	return
}
