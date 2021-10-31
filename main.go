// Conversion of ES6+ into ES5.1 (wip)
//
// TODO: turn into a script that uses gojafs
package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"github.com/jvatic/goja-babel"
	"io"
	"log"
	"os"
	"strings"
)

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
			"transform-destructuring",
			"transform-spread",
			"transform-parameters",
			"transform-async-to-generator",
			"transform-regenerator",
			"transform-for-of",
			"proposal-optional-chaining",
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
