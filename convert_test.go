package main

import (
	"github.com/dop251/goja"
	"github.com/tdewolff/parse/v2"
	"github.com/tdewolff/parse/v2/js"
	"testing"
)

func TestClassToProto(t *testing.T) {
	src := `
	class A {
		constructor(w, h) {
			this.width = w;
			this.height = h;
		}

		info() {
			return String(this.width) + ' ' + String(this.height);
		}
	}
	`
	ast, err := js.Parse(parse.NewInputString(src), js.Options{})
	if err != nil {
		t.Fatalf("%v", err)
	}
	cls := ast.BlockStmt.List[0].(*js.ClassDecl)
	p, err := classToProto(cls, &ast.BlockStmt.Scope)
	if err != nil {
		t.Fatalf("%v", err)
	}
	t.Logf("p=%+v", p)
	t.Logf("p'=%+v", p.JS())
	vm := goja.New()
	v, err := vm.RunString(p.JS() + " a = new A(2, 5); a.info();")
	if err != nil {
		panic(err)
	}
	res := v.Export().(string)
	if res != "2 5" {
		t.Fatalf("%v", res)
	}
}

func TestClassToProtoExtends(t *testing.T) {
	src := `
	class A {
		constructor(w, h) {
			this.width = w;
			this.height = h;
		}

		info() {
			return String(this.width) + ' ' + String(this.height);
		}
	}

	class B extends A {
		constructor() {
			super(10, 20);
		}
	}
	`
	ast, err := js.Parse(parse.NewInputString(src), js.Options{})
	if err != nil {
		t.Fatalf("%v", err)
	}
	clsA := ast.BlockStmt.List[0].(*js.ClassDecl)
	clsB := ast.BlockStmt.List[1].(*js.ClassDecl)
	p, err := classToProto(clsA, &ast.BlockStmt.Scope)
	if err != nil {
		t.Fatalf("%v", err)
	}
	q, err := classToProto(clsB, &ast.BlockStmt.Scope)
	if err != nil {
		t.Fatalf("%v", err)
	}
	vm := goja.New()
	src = p.JS() + " " + q.JS() + " b = new B(); b.info();"
	t.Logf("src=%+v", src)
	v, err := vm.RunString(src)
	if err != nil {
		panic(err)
	}
	res := v.Export().(string)
	if res != "10 20" {
		t.Fatalf("%v", res)
	}
}

func TestClassToProtoAnon(t *testing.T) {
	src := `
	new class {
		constructor(w, h) {
			this.width = w;
			this.height = h;
		}

		info() {
			return String(this.width) + ' ' + String(this.height);
		}
	}
	`
	ast, err := js.Parse(parse.NewInputString(src), js.Options{})
	if err != nil {
		t.Fatalf("%v", err)
	}
	cls := ast.BlockStmt.List[0].(*js.ExprStmt).Value.(*js.NewExpr).X.(*js.ClassDecl)
	p, err := classToProto(cls, &ast.BlockStmt.Scope)
	if err != nil {
		t.Fatalf("%v", err)
	}
	t.Logf("p=%+v", p)
	t.Logf("p'=%+v", p.JS())
	vm := goja.New()
	src = "a = new "+p.JS() + "(2, 5); a.info();"
	t.Logf("src=%v", src)
	v, err := vm.RunString(src)
	if err != nil {
		panic(err)
	}
	res := v.Export().(string)
	if res != "2 5" {
		t.Fatalf("%v", res)
	}
}
