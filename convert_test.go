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
		constructor(h, w) {
			this.height = h;
			this.width = w;
		}

		info() {
			return String(this.height) + ' ' + String(this.width);
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
