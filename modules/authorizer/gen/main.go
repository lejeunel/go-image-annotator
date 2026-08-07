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
	ifaceName := flag.String("iface", "", "name of the interface to extract methods from")
	in := flag.String("in", "", "source file containing the interface declaration")
	out := flag.String("out", "", "output file")
	pkg := flag.String("pkg", "", "package name for generated file")
	flag.Parse()

	if *ifaceName == "" || *in == "" || *out == "" || *pkg == "" {
		fmt.Fprintln(os.Stderr, "iface, in, out, and pkg flags are all required")
		os.Exit(1)
	}

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, *in, nil, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse error: %v\n", err)
		os.Exit(1)
	}

	var target *ast.InterfaceType
	for _, decl := range node.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok || gen.Tok != token.TYPE {
			continue
		}
		for _, spec := range gen.Specs {
			ts, ok := spec.(*ast.TypeSpec)
			if !ok || ts.Name.Name != *ifaceName {
				continue
			}
			it, ok := ts.Type.(*ast.InterfaceType)
			if !ok {
				fmt.Fprintf(os.Stderr, "%q is not an interface type in %s\n", *ifaceName, *in)
				os.Exit(1)
			}
			target = it
		}
	}
	if target == nil {
		fmt.Fprintf(os.Stderr, "interface %q not found in %s\n", *ifaceName, *in)
		os.Exit(1)
	}

	var methods []string
	if target.Methods != nil {
		for _, field := range target.Methods.List {
			if len(field.Names) == 0 {
				continue // embedded interface; deliberately not followed
			}
			for _, n := range field.Names {
				methodName := n.Name
				if !isFirstLetterCapitalized(methodName) {
					continue
				}
				if slices.Contains(SkipMethods, methodName) {
					continue
				}
				methods = append(methods, methodName)
			}
		}
	}

	if len(methods) == 0 {
		fmt.Fprintf(os.Stderr, "no methods found on interface %q in %s\n", *ifaceName, *in)
		os.Exit(1)
	}

	sort.Strings(methods)

	var b strings.Builder
	fmt.Fprintf(
		&b,
		"// Code generated automatically; DO NOT EDIT.\n\npackage %s\n\nvar ValidMethods = []string{\n",
		*pkg,
	)
	for _, m := range methods {
		fmt.Fprintf(&b, "\t%q,\n", m)
	}
	fmt.Fprintf(&b, "\t\"*\",\n")
	b.WriteString("}\n")

	if err := os.WriteFile(*out, []byte(b.String()), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "write error: %v\n", err)
		os.Exit(1)
	}
}
