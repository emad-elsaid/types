package types

import _ "embed"

//go:embed array.go
var ArrayTmpl string

//go:embed array_test.go
var ArrayTestTmpl string
