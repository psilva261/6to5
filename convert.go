package main

import (
	"github.com/tdewolff/parse/v2/js"
)

func Convert(ast *js.AST) {
	js.Walk(convertCls{}, ast)
}
