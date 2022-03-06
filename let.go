package main

import (
	"github.com/tdewolff/parse/v2/js"
	"strconv"
)

var lastUnique int

func unique() string {
	lastUnique++
	return strconv.Itoa(lastUnique)
}

type convertLetConst struct {
	parent js.INode
	scope  *js.Scope
	block  *js.BlockStmt
}

func ConvertLetConst(ast *js.AST) {
	js.Walk(convertLetConst{}, ast)
}

func (c convertLetConst) Enter(n js.INode) js.IVisitor {
	switch v := n.(type) {
	case *js.VarDecl:
		for _, be := range v.List {
			vr := be.Binding.(*js.Var)
			if v.TokenType == js.VarToken && vr.Decl == js.VariableDecl {
				continue
			}
			v.TokenType = js.VarToken
			vr.Decl = js.VariableDecl
			if c.scope.Parent == nil {
				break
			}
			vr.Data = []byte(string(vr.Data) + "$" + unique())
		}
	}
	if blk, ok := n.(*js.BlockStmt); ok {
		c.block = blk
		c.scope = &blk.Scope
	}
	c.parent = n
	return c
}

func (c convertLetConst) Exit(n js.INode) {}
