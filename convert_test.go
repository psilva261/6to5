package main

import (
	"github.com/robertkrimen/otto"
	"github.com/tdewolff/parse/v2"
	"github.com/tdewolff/parse/v2/js"
	"testing"
)

func TestConvert(t *testing.T) {
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
	var a = new A(10, 20);
	a.info();
	`
	ast, err := js.Parse(parse.NewInputString(src), js.Options{})
	if err != nil {
		t.Fatalf("%v", err)
	}
	Convert(ast)
	src = ast.JS()
	t.Logf("%v", src)
	vm := otto.New()
	v, err := vm.Run(src)
	if err != nil {
		t.Fatalf("%v", err)
	}
	res, err := v.Export()
	if err != nil {
		t.Fatalf("%v", err)
	}
	if res.(string) != "10 20" {
		t.Fatalf("%v", res)
	}
}

func TestConvert2(t *testing.T) {
	srcs := []string{
		`class A {}`,
		`new class {}`,
		`var cls = class {}`,
	}
	for _, src := range srcs {
		ast, err := js.Parse(parse.NewInputString(src), js.Options{})
		if err != nil {
			t.Fatalf("%v", err)
		}
		Convert(ast)
		src = ast.JS()
		t.Logf("%v", src)
		vm := otto.New()
		_, err = vm.Run(src)
		if err != nil {
			panic(err)
		}
	}
}

func TestConvertExport(t *testing.T) {
	src := `
let a = class {
    constructor(e) {
    }

    foo() {}
};

export { a as A };
	`
	ast, err := js.Parse(parse.NewInputString(src), js.Options{})
	if err != nil {
		t.Fatalf("%v", err)
	}
	Convert(ast)
	src = ast.JS()
	t.Logf("%v", src)
	vm := otto.New()
	_, err = vm.Run(src)
	if err != nil {
		panic(err)
	}
}

func TestConvertLetConst(t *testing.T) {
	src := `
const a = 10;
let r = 0;
for (let i = a; i <= a; i++) {
	r += i;
}
(function() {
	let a = 1000;
	r += a;
})()
r += 10*a;
r;
	`
	ast, err := js.Parse(parse.NewInputString(src), js.Options{})
	if err != nil {
		t.Fatalf("%v", err)
	}
	Convert(ast)
	src = ast.JS()
	t.Logf("%v", src)
	vm := otto.New()
	v, err := vm.Run(src)
	if err != nil {
		panic(err)
	}
	res, err := v.Export()
	if err != nil {
		t.Fatalf("%v", err)
	}
	if res.(float64) != 1110 {
		t.Fatalf("%v", res)
	}
}
