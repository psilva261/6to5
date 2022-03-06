package main

import (
	"fmt"
	"github.com/tdewolff/parse/v2/js"
)

type convertArrow struct {
	parent js.INode
	scope  *js.Scope
	block  *js.BlockStmt
}

func ConvertArrow(ast *js.AST) {
	js.Walk(convertArrow{}, ast)
}

func (c convertArrow) Enter(n js.INode) js.IVisitor {
	switch v := n.(type) {
	case *js.ArrowFunc:
		f := &js.FuncDecl{
			Async:  v.Async,
			Params: v.Params,
			Body:   v.Body,
		}
		that := AddThat(c.block, c.scope)
		ft := BindThat(f, that)
		switch vv := c.parent.(type) {
		case *js.Arg:
			vv.Value = ft.(js.IExpr)
		case *js.BindingElement:
			vv.Default = ft.(js.IExpr)
		case *js.GroupExpr:
			vv.X = ft.(js.IExpr)
		case *js.Property:
			vv.Value = ft.(js.IExpr)
		default:
			panic(fmt.Sprintf("unknown parent %T", c.parent))
		}
	}
	if blk, ok := n.(*js.BlockStmt); ok {
		c.block = blk
		c.scope = &blk.Scope
	}
	c.parent = n
	return c
}

func (c convertArrow) Exit(n js.INode) {}

func BindThat(f *js.FuncDecl, that *js.Var) js.INode {
	bind := &js.DotExpr{
		X: f,
		Y: js.LiteralExpr{
			Data: []byte("bind"),
		},
	}
	c := &js.CallExpr{
		X: bind,
		Args: js.Args{
			List: []js.Arg{
				js.Arg{
					Value: that,
				},
			},
		},
	}
	return c
}

func getThat(blk *js.BlockStmt) (that *js.Var, ok bool) {
	for _, st := range blk.List {
		d, ok := st.(*js.VarDecl)
		if !ok {
			continue
		}
		if be := d.List[0]; string(be.Default.String()) == "this" {
			v, ok := be.Binding.(*js.Var)
			if ok {
				return v, true
			}
		}
	}
	return
}

func AddThat(blk *js.BlockStmt, scope *js.Scope) (that *js.Var) {
	if that, ok := getThat(blk); ok {
		return that
	}
	that = &js.Var{
		Decl: js.VariableDecl,
		Data: []byte("that$" + unique()),
	}
	d := &js.VarDecl{
		TokenType: js.VarToken,
		List: []js.BindingElement{
			js.BindingElement{
				Default: &js.LiteralExpr{
					Data: []byte("this"),
				},
				Binding: that,
			},
		},
		Scope: scope,
	}
	blk.List = append([]js.IStmt{d}, blk.List...)
	return
}
