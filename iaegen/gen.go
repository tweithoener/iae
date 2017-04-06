package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

var src = `
	package aaa

	// aa squares a
	// but only if a is smaller then 10
	//
	// iae: a != 123
	func aa (int a) {
		// iae: a<10

		// return squared a
		return a*a
	}
`

var X = ""

func findPackage(n ast.Node) bool {
	if n == nil {
		return false
	}
	se, ok := n.(*ast.SelectorExpr)
	if !ok {
		return true
	}
	if x, ok := se.X.(*ast.Ident); ok {
		X = x.Name
		println("find: X is identifier", X)
	}
	return true
}

func main() {

	/*
		file := "gen.go"
		buf, err := ioutil.ReadFile(file)
		if err != nil {
			panic(err)
		}
	*/

	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, "gen.go", src, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	ast.Print(fset, f)

	cmap := ast.NewCommentMap(fset, f, f.Comments)
	for _, currcg := range f.Comments {
		println("comment group:", currcg.Text())
		for _, co := range currcg.List {
			text := co.Text
			text = strings.TrimSpace(text)
			text = strings.TrimLeft(text, "/")
			text = strings.TrimSpace(text)
			if !strings.HasPrefix(text, "iae:") {
				println("not relevant: ", co.Text)
				continue
			}
			text = strings.TrimLeft(text, "iae:")
			text = strings.TrimSpace(text)
			println("relevant comment: ", co.Text)

			for n, cgs := range cmap {
				for _, cg := range cgs {
					if currcg == cg {
						println("comment belongs to: ")
						ast.Print(fset, n)
					}
				}
			}
		}
	}

	ast.Inspect(f, func(n ast.Node) bool {
		if n == nil {
			return false
		}

		co, ok := n.(*ast.Comment)
		if ok {
			println("comment: ", co.Text)
		}

		ce, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		se, ok := ce.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		ast.Inspect(n, findPackage)
		println("selector expression")
		if X != "iae" {
			return true
		}
		s := se.Sel
		if s.Name != "Arg" {
			return true
		}
		println("found a relevant call:", s.Name)
		args := ce.Args
		if len(args) != 1 {
			println("already converted", len(args))
			return false
		}
		println("needs conversion")
		start := ce.Lparen
		end := ce.Rparen
		println(string(src[start:end]))

		return false
	})

}
