package types

import "fmt"

func ExampleSlice() {
	a := Slice[int]{1, 2, 3, 4, 5, 6}

	// multiply every element by 100
	a = a.Map(func(e int) int {
		return e * 100
	})

	// select any element <= 300
	a = a.KeepIf(func(e int) bool {
		return e <= 300
	})

	fmt.Print(a)

	// Output: [100 200 300]
}
