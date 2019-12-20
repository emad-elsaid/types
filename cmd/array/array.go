//go:generate cp ../../array.go template.go
//go:generate sed -e "s/^package\stypes/package\smain\ntemp:=`/g" -e s/type\sElement\sinterface{}// -e s/ElementArray/{.array}/g -e s/Element/{.item}/g -i template.go
//go:generate sed "$ a `" -i template.go
//go:generate

package main

func main() {

}
