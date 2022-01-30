package main

import (
	"testing"
)

func TestFindFeatures(t *testing.T) {
	f, _ := FindFeatures(`var f = (x) => x*x;`)
	if !f.ArrowFunctions {
		t.Fail()
	}
	f, _ = FindFeatures(`let a = 1; { let a = 2; }`)
	if !f.BlockScoping {
		t.Fail()
	}
	f, _ = FindFeatures(`class A {}`)
	if !f.Classes {
		t.Fail()
	}
	f, _ = FindFeatures(`console.log(...[1, 2]);`)
	if !f.Spread {
		t.Fail()
	}
	f, _ = FindFeatures(`async function f() {}`)
	if !f.Async {
		t.Fail()
	}
	f, _ = FindFeatures(`class A { async f() {} }`)
	if !f.Async {
		t.Fail()
	}
	f, _ = FindFeatures(`function *f() {}`)
	if !f.Generators {
		t.Fail()
	}
	f, _ = FindFeatures(`class A { *f() {} }`)
	if !f.Generators {
		t.Fail()
	}
}
