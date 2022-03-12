package main

import (
	"github.com/robertkrimen/otto"
	"github.com/tdewolff/parse/v2"
	"github.com/tdewolff/parse/v2/js"
	"testing"
)

func TestRegexp(t *testing.T) {
	src := `
	/* https://babeljs.io/docs/en/babel-plugin-transform-sticky-regex */
	/* https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/RegExp/sticky */
	var str = 'table football';
	var regex = /foo/y;
	regex.test(str);
	`
	ast, err := js.Parse(parse.NewInputString(src), js.Options{})
	if err != nil {
		t.Fatalf("%v", err)
	}
	js.Walk(convertRegexp{}, ast)
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
	if res.(bool) != true {
		t.Fatalf("%v", res)
	}
}
