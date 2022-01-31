package main

import (
	"errors"
	"github.com/tdewolff/parse/v2/js"
)

/*
https://medium.com/@dange.laxmikant/simplified-inheritance-in-js-es5-way-60b4ff19b008

function C(a) {
	this.a=a;
}
C.prototype.info=function() {console.log(this.a);}
function D() {
	C.call(this, 1); // "super"
}
D.prototype=Object.create(C);
d=new D();
d.info();
*/

type supr struct {
	extends js.IExpr
}

func (s supr) Enter(n js.INode) js.IVisitor {
	if v, ok := n.(*js.CallExpr); ok && v.X.String() == "super" {
		v.X = dotExpr(s.extends.String(), "call")
		v.Args.List = append([]js.Arg{
			js.Arg{
				Value: &js.Var{
					Data: []byte("this"),
				},
			},
		}, v.Args.List...)
	}
	return s
}

func (s supr) Exit(n js.INode) {}

var ErrAnonClass = errors.New("cannot process anonymous class")

func classToProto(cls *js.ClassDecl, parent *js.Scope) (p js.BlockStmt, err error) {
	if cls.Name == nil {
		return p, ErrAnonClass
	}
	f := js.FuncDecl{
		Name: &(*cls.Name),
	}
	f.Body.Scope.Parent = parent
	construct, other := methods(cls)
	if construct != nil {
		if cls.Extends != nil {
			js.Walk(supr{cls.Extends}, construct)
		}
		mt := js.FuncDecl{}
		mt.Body.Scope.Parent = &f.Body.Scope
		f.Params = construct.Params
		for _, stmt := range construct.Body.List {
			f.Body.List = append(f.Body.List, stmt)
		}
	}
	p.List = append(p.List, f)
	if cls.Extends != nil {
		p.List = append(p.List, extend(cls.Name, cls.Extends))
	}
	for _, m := range other {
		name := &js.Var{
			Data: []byte(m.Name.String()),
		}
		mt := js.FuncDecl{}
		mt.Body.Scope.Parent = &f.Body.Scope
		for _, stmt := range m.Body.List {
			mt.Body.List = append(mt.Body.List, stmt)
		}
		mt.Params = m.Params
		p.List = append(p.List, regMethod(cls.Name, name, mt))
	}
	return
}

func extend(name *js.Var, extends js.IExpr) (s js.ExprStmt) {
	s.Value = js.BinaryExpr{
		Op: js.EqToken,
		X: dotExpr(name.String(), "prototype"),
		Y: js.CallExpr{
			X: dotExpr("Object", "create"),
			Args: js.Args{
				List: []js.Arg{
					js.Arg{
						Value: dotExpr(extends.String(), "prototype"),
					},
				},
			},
		},
	}
	return
}

func regMethod(clsName, name *js.Var, mt js.FuncDecl) (s js.ExprStmt) {
	s.Value = js.BinaryExpr{
		Op: js.EqToken,
		X:  dotExpr(clsName.String(), "prototype", name.String()),
		Y:  mt,
	}
	return
}

func dotExpr(items ...string) (expr js.DotExpr) {
	if len(items) == 2 {
		expr.X = js.LiteralExpr{
			Data: []byte(items[0]),
		}

	} else {
		expr.X = dotExpr(items[:len(items)-1]...)
	}
	expr.Y = js.LiteralExpr{
		Data: []byte(items[len(items)-1]),
	}
	return
}

func methods(cls *js.ClassDecl) (constructor *js.MethodDecl, other []*js.MethodDecl) {
	other = make([]*js.MethodDecl, 0, len(cls.Methods))
	for _, m := range cls.Methods {
		if m.Name.String() == "constructor" {
			constructor = m
		} else {
			other = append(other, m)
		}
	}
	return
}
