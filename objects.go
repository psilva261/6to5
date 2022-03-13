package main

import (
	"github.com/tdewolff/parse/v2/js"
)

type convertObjects struct {
	parent js.INode
	scope  *js.Scope
	block  *js.BlockStmt
}

func ConvertObjects(ast *js.AST) {
	js.Walk(convertObjects{}, ast)
}

func (c convertObjects) Enter(n js.INode) js.IVisitor {
	switch v := n.(type) {
	case *js.Property:
		if (v.Name == nil || v.Name.IsIdent([]byte(v.Value.String()))) && !v.Spread {
			v.Name = &js.PropertyName{
				Literal: js.LiteralExpr{
					Data: []byte(v.Value.String()),
				},
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

func (c convertObjects) Exit(n js.INode) {}
