package main

import (
	"github.com/robertkrimen/otto"
	"github.com/tdewolff/parse/v2"
	"github.com/tdewolff/parse/v2/js"
	"testing"
)

func TestObjects(t *testing.T) {
	src := `
	var a=2, b=5;
	var obj = {a:a, b: b};
	obj.a+obj.b;
	`
	ast, err := js.Parse(parse.NewInputString(src), js.Options{})
	if err != nil {
		t.Fatalf("%v", err)
	}
	js.Walk(convertObjects{}, ast)
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
	if res.(float64) != 7 {
		t.Fatalf("%v", res)
	}
}
