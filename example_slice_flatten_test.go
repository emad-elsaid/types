package types

import "fmt"

func ExampleFlatten() {
	nested := [][]int{
		{1, 2, 3},
		{4, 5},
		{6, 7, 8, 9},
	}

	flat := Flatten(nested)
	fmt.Println(flat)

	// Output: [1 2 3 4 5 6 7 8 9]
}

func ExampleFlatten_strings() {
	words := [][]string{
		{"hello", "world"},
		{"foo", "bar"},
		{"baz"},
	}

	result := Flatten(words)
	fmt.Println(result)

	// Output: [hello world foo bar baz]
}

func ExampleFlatten_withEmptySlices() {
	nested := [][]int{
		{1, 2},
		{},
		{3, 4},
		{},
		{5},
	}

	flat := Flatten(nested)
	fmt.Println(flat)

	// Output: [1 2 3 4 5]
}
