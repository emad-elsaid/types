package types

import "fmt"

func ExampleSlice_Zip() {
	// Zip two slices of different types
	numbers := Slice[int]{1, 2, 3, 4}
	letters := Slice[string]{"a", "b", "c"}

	// Zip combines them element-wise
	// Result length is determined by the shorter slice
	pairs := numbers.Zip(letters)

	for _, pair := range pairs {
		fmt.Printf("(%v, %v) ", pair[0], pair[1])
	}

	// Output: (1, a) (2, b) (3, c)
}

func ExampleSlice_Zip_sameTypes() {
	// Zip works with slices of the same type too
	first := Slice[int]{1, 2, 3}
	second := Slice[int]{10, 20, 30}

	pairs := first.Zip(second)

	for _, pair := range pairs {
		sum := pair[0].(int) + pair[1].(int)
		fmt.Printf("%v+%v=%v ", pair[0], pair[1], sum)
	}

	// Output: 1+10=11 2+20=22 3+30=33
}
