package main

import (
	"github.com/robertkrimen/otto"
	"github.com/tdewolff/parse/v2"
	"github.com/tdewolff/parse/v2/js"
	"testing"
)

func TestArrow(t *testing.T) {
	src := `
	var f = (x, y) => x + y;
	f(2, 3);
	`
	ast, err := js.Parse(parse.NewInputString(src), js.Options{})
	if err != nil {
		t.Fatalf("%v", err)
	}
	js.Walk(convertArrow{}, ast)
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
	if res.(float64) != 5 {
		t.Fatalf("%v", res)
	}
}

func TestArrow2(t *testing.T) {
	src := `
	var a = [1, 2, 3];
	var m = function(ary, f) {
		var res = Array(ary.length);
		for (var i = 0; i < ary.length; i++) {
			res[i] = f(ary[i]);
		}
		return res;
	}
	a = m(a, x => x * x);
	a[0] + a[1] + a[2];
	`
	ast, err := js.Parse(parse.NewInputString(src), js.Options{})
	if err != nil {
		t.Fatalf("%v", err)
	}
	js.Walk(convertArrow{}, ast)
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
	if res.(float64) != 14 {
		t.Fatalf("%v", res)
	}
}

func TestArrow3(t *testing.T) {
	src := `
	(x => x * x)(3);
	`
	ast, err := js.Parse(parse.NewInputString(src), js.Options{})
	if err != nil {
		t.Fatalf("%v", err)
	}
	js.Walk(convertArrow{}, ast)
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
	if res.(float64) != 9 {
		t.Fatalf("%v", res)
	}
}
