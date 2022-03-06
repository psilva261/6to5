package main

import (
	"github.com/robertkrimen/otto"
	"github.com/tdewolff/parse/v2"
	"github.com/tdewolff/parse/v2/js"
	"testing"
)

func TestLetConst(t *testing.T) {
	src := `
	const a = 2;
	let b = 3;
	(function() {
		let c = 5;
		const d = 7;
		b += a + c + d;
	})()
	b;
	`
	ast, err := js.Parse(parse.NewInputString(src), js.Options{})
	if err != nil {
		t.Fatalf("%v", err)
	}
	js.Walk(convertLetConst{}, ast)
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
	if res.(float64) != 17 {
		t.Fatalf("%v", res)
	}
}