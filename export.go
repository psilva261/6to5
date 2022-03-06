package main

import (
	"github.com/tdewolff/parse/v2/js"
)

type convertExport struct {
	parent js.INode
	scope  *js.Scope
	block  *js.BlockStmt
}

func ConvertExport(ast *js.AST) {
	js.Walk(convertExport{}, ast)
}

func (c convertExport) Enter(n js.INode) js.IVisitor {
	switch n.(type) {
	case *js.ExportStmt:
		found := -1
		for i, is := range c.block.List {
			if is == n {
				_ = is
				c.block.List[i] = &js.EmptyStmt{}
				found = i
				break
			}
		}
		if found < 0 {
			panic("unreachable")
		}
	}
	if blk, ok := n.(*js.BlockStmt); ok {
		c.block = blk
		c.scope = &blk.Scope
	}
	c.parent = n
	return c
}

func (c convertExport) Exit(n js.INode) {}
