package types

import "fmt"

func ExampleArray() {
	a := Array{1, 2, 3, 4, 5, 6}

	// multiply every element by 100
	a = a.Map(func(e Element) Element {
		return e.(int) * 100
	})

	// select any element <= 300
	a = a.KeepIf(func(e Element) bool {
		return e.(int) <= 300
	})

	fmt.Print(a)

	// Output: [100 200 300]
}
