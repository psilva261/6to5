package main

import (
	"github.com/tdewolff/parse/v2"
	"github.com/tdewolff/parse/v2/js"
	"testing"
)

func TestFindFeatures(t *testing.T) {
	run := func(src string) Features {
		ast, err := js.Parse(parse.NewInputString(src), js.Options{})
		if err != nil {
			t.Fatalf("%v", err)
		}
		f, _ := FindFeatures(ast, src)
		return f
	}
	f := run(`var f = (x) => x*x;`)
	if !f.ArrowFunctions {
		t.Fail()
	}
	f = run(`class A {}`)
	if !f.Classes {
		t.Fail()
	}
	f = run(`console.log(...[1, 2]);`)
	if !f.Spread {
		t.Fail()
	}
	f = run(`async function f() {}`)
	if !f.Async {
		t.Fail()
	}
	f = run(`class A { async f() {} }`)
	if !f.Async {
		t.Fail()
	}
	f = run(`function *f() {}`)
	if !f.Generators {
		t.Fail()
	}
	f = run(`class A { *f() {} }`)
	if !f.Generators {
		t.Fail()
	}
}
