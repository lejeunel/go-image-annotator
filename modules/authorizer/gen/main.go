package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"slices"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"
)

var SkipMethods = []string{"SetAuthRules"}

func isFirstLetterCapitalized(s string) bool {
	if s == "" {
		return false
	}
	r, _ := utf8.DecodeRuneInString(s)
	return unicode.IsUpper(r)
}

func main() {
	structName := flag.String("struct", "", "name of the struct to extract methods from")
	in := flag.String("in", "", "source file containing the struct's methods")
	out := flag.String("out", "", "output file")
	pkg := flag.String("pkg", "", "package name for generated file")
	flag.Parse()

	if *structName == "" || *in == "" || *out == "" || *pkg == "" {
		fmt.Fprintln(os.Stderr, "struct, in, out, and pkg flags are all required")
		os.Exit(1)
	}

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, *in, nil, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse error: %v\n", err)
		os.Exit(1)
	}

	var methods []string
	for _, decl := range node.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Recv == nil || len(fn.Recv.List) == 0 {
			continue // a plain function, not a method
		}
		if !isFirstLetterCapitalized(fn.Name.Name) {
			continue
		}
		if slices.Contains(SkipMethods, fn.Name.Name) {
			continue
		}
		if recvTypeName(fn.Recv.List[0].Type) == *structName {
			methods = append(methods, fn.Name.Name)
		}
	}

	if len(methods) == 0 {
		fmt.Fprintf(os.Stderr, "no methods found on struct %q in %s\n", *structName, *in)
		os.Exit(1)
	}
	sort.Strings(methods)

	var b strings.Builder
	fmt.Fprintf(&b, "// Code generated automatically; DO NOT EDIT.\n\npackage %s\n\nvar validMethods = []string{\n", *pkg)
	for _, m := range methods {
		fmt.Fprintf(&b, "\t%q,\n", m)
	}
	b.WriteString("}\n")

	if err := os.WriteFile(*out, []byte(b.String()), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "write error: %v\n", err)
		os.Exit(1)
	}
}

// recvTypeName extracts the base type name from a receiver expression,
// handling value receivers (T), pointer receivers (*T), and generic
// receivers (T[X]) / (*T[X]).
func recvTypeName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return recvTypeName(t.X)
	case *ast.IndexExpr:
		return recvTypeName(t.X)
	case *ast.IndexListExpr:
		return recvTypeName(t.X)
	default:
		return ""
	}
}
