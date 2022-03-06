package main

import (
	"github.com/tdewolff/parse/v2/js"
)

func Convert(ast *js.AST) {
	js.Walk(convertLetConst{}, ast)
	js.Walk(convertExport{}, ast)
	js.Walk(convertCls{}, ast)
	js.Walk(convertArrow{}, ast)
}
