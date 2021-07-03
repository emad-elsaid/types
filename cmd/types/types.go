package main

import (
	"bytes"
	"flag"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"os"

	"github.com/emad-elsaid/types"
	"golang.org/x/tools/go/ast/astutil"
)

func main() {
	var pkg = flag.String("package", "types", "package name the new file will belong to")
	var element = flag.String("element", "string", "the single element of your array")
	var array = flag.String("array", "stringArray", "the name of the slice of your element")
	var output = flag.String("output", "/dev/stdout", "file name of the output")
	flag.Parse()

	replacements := map[string]string{
		"types":   *pkg,
		"Element": *element,
		"Array":   *array,
	}

	arrayOut := parseAndReplace("array.go", types.ArrayTmpl, replacements)
	os.WriteFile(*output, arrayOut, 0755)
}

func parseAndReplace(inputFileName, inputContent string, replacements map[string]string) []byte {
	fset := token.NewFileSet()
	parsed, err := parser.ParseFile(fset, inputFileName, inputContent, parser.ParseComments)
	if err != nil {
		log.Fatalf("error parsing file: %s", err)
	}

	astutil.Apply(parsed, func(cr *astutil.Cursor) bool {
		t, ok := cr.Node().(*ast.Ident)
		if !ok {
			return true
		}

		if v, ok := replacements[t.Name]; ok {
			t.Name = v
		}

		return true
	}, nil)

	out := bytes.NewBuffer([]byte{})
	err = printer.Fprint(out, fset, parsed)
	if err != nil {
		log.Fatalf("error serializing: %s", err)
	}

	return out.Bytes()
}
