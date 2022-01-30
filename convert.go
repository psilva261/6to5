package main

import (
	"errors"
	"github.com/tdewolff/parse/v2/js"
)

func classToProto(cls *js.ClassDecl, parent *js.Scope) (p js.BlockStmt, err error) {
	if cls.Name == nil {
		return p, errors.New("no class name")
	}
	f := js.FuncDecl{
		Name: &(*cls.Name),
	}
	f.Body.Scope.Parent = parent
	for _, def := range cls.Definitions {
		_ = def
	}
	for _, m := range cls.Methods {
		name := &js.Var{
			Data: []byte(m.Name.String()),
		}
		mt := js.FuncDecl{
			Name: name,
		}
		mt.Body.Scope.Parent = &f.Body.Scope
		for _, stmt := range m.Body.List {
			mt.Body.List = append(mt.Body.List, stmt)
		}
		f.Body.List = append(f.Body.List, mt)
		if m.Name.String() == "constructor" {
			f.Params = m.Params
			cl := &js.CallExpr{
				X: &js.DotExpr{
					X: &js.Var{
						Data: []byte("constructor"),
					},
					Y: js.LiteralExpr{
						Data: []byte("call"),
					},
				},
			}
			cl.Args.List = append(cl.Args.List, js.Arg{Value: &js.Var{Data: []byte("this")}})
			for _, par := range m.Params.List {
				arg := js.Arg{
					Value: par.Default,
				}
				_ = arg
				//cl.Args.List = append(cl.Args.List, arg)
			}
			clSt := js.ExprStmt{
				Value: cl,
			}
			f.Body.List = append(f.Body.List, clSt)
			_ = clSt
		} else {
			mt.Params = m.Params
			f.Body.List = append(f.Body.List, &js.ExprStmt{
				Value: js.BinaryExpr{
					Op: js.EqToken,
					X: js.DotExpr{
						X: &js.Var{
							Data: []byte("this"),
						},
						Y: js.LiteralExpr{
							Data: []byte(m.Name.String()),
						},
					},
					Y: &js.Var{
						Data: []byte(m.Name.String()),
					},
				},
			})
		}
	}
	p.List = append(p.List, f)
	return
}
