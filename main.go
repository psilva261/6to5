// Conversion of ES6+ into ES5.1 (wip)
//
// TODO: turn into a script that uses gojafs
package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"github.com/jvatic/goja-babel"
	"github.com/tdewolff/parse/v2"
	"github.com/tdewolff/parse/v2/js"
	"io"
	"log"
	"os"
	"strings"
)

type Features struct {
	ArrowFunctions bool
	BlockScoping   bool
	Classes        bool
	Spread         bool
	Async          bool
	Generators     bool
}

func FindFeatures(src string) (f Features, err error) {
	ast, err := js.Parse(parse.NewInputString(src), js.Options{})
	if err != nil {
		return
	}
	w := walker{
		Features: &f,
	}
	js.Walk(w, ast)
	if strings.Contains(src, "...") {
		w.Spread = true
	}
	return
}

func (f Features) Unsupported() bool {
	return f.ArrowFunctions || f.BlockScoping || f.Classes || f.Spread || f.Async || f.Generators
}

type walker struct{
	*Features
}

func (w walker) Enter(n js.INode) js.IVisitor {
	switch n := n.(type) {
	case *js.ArrowFunc:
		w.ArrowFunctions = true
		if n.Async {
			w.Async = true
		}
	case *js.BlockStmt:
		w.BlockScoping = true
	case *js.ClassDecl:
		w.Classes = true
	case *js.FuncDecl:
		if n.Async {
			w.Async = true
		}
		if n.Generator {
			w.Generators = true
		}
	case *js.MethodDecl:
		if n.Async {
			w.Async = true
		}
		if n.Generator {
			w.Generators = true
		}
	}
	return w
}

func (w walker) Exit(n js.INode) {
	
}

var Cache = os.Getenv("ES6TO5_CACHE")

func cached(hs string) (f *os.File, ok bool) {
	if Cache == "" {
		return
	}
	f, err := os.Open(Cache+"/"+hs)
	ok = err == nil
	return
}

func cache(hs string, bs []byte) {
	if Cache == "" {
		return
	}
	f, err := os.Create(Cache+"/"+hs)
	if err != nil {
		log.Printf("os create: %v", err)
		return
	}
	f.Write(bs)
	f.Close()
}

func Main() (err error) {
	buf := bytes.NewBufferString("")
	if _, err = io.Copy(buf, os.Stdin); err != nil {
		return
	}
	f, err := FindFeatures(buf.String())
	if err == nil && !f.Unsupported() {
		fmt.Println(buf.String())
		return
	}
	h := sha256.New()
	h.Write(buf.Bytes())
	hs := fmt.Sprintf("%x", h.Sum(nil))
	if f, ok := cached(hs); ok {
		defer f.Close()
		_, err = io.Copy(os.Stdout, f)
		return
	}
	babel.Init(1) // Setup 1 transformer (can be any number > 0)
	r, err := babel.Transform(
		strings.NewReader(buf.String()),
		map[string]interface{}{
		"plugins": []string{
			"transform-arrow-functions",
			"transform-block-scoping",
			"transform-classes",
			"transform-spread",
			"transform-parameters",
			"transform-async-to-generator",
			"transform-regenerator",
		},
	})
	if err != nil {
		return fmt.Errorf("transform: %v", err)
	}
	bs, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("read all: %v", err)
	}
	cache(hs, bs)
	fmt.Println(string(bs))

	return
}

func main() {
	if err := Main(); err != nil {
		log.Fatalf("%v",err)
	}
}
