//go:generate cp ../../array.go template.go
//go:generate sed s/type\sElement\sinterface{}// -i template.go
//go:generate sed s/ElementArray/{{.Array}}/g  -i template.go
//go:generate sed s/Element/{{.Element}}/g -i template.go
//go:generate sed 2s/types/{{.Package}}/ -i template.go
//go:generate sed "1 a const tmpl string = `" -i template.go
//go:generate sed "1 a package main" -i template.go
//go:generate sed "$ a `" -i template.go

package main

import (
	"flag"
	"log"
	"os"
	"text/template"
)

type params struct {
	Package string
	Element string
	Array   string
}

func main() {
	var pkg = flag.String("package", "main", "package name the new file will belong to")
	var element = flag.String("element", "string", "the single element of your array")
	var array = flag.String("array", "stringArray", "the name of the slice of your element")
	var output = flag.String("output", "/dev/stdout", "where to write the output")
	flag.Parse()

	p := params{
		Package: *pkg,
		Element: *element,
		Array:   *array,
	}

	f, err := os.OpenFile(*output, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	t := template.Must(template.New("template").Parse(tmpl))
	t.ExecuteTemplate(f, "template", p)
}
