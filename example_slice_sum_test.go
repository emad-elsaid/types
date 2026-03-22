package types

import "fmt"

func ExampleSliceSum() {
	numbers := Slice[int]{1, 2, 3, 4, 5}
	total := SliceSum(numbers)
	fmt.Println(total)
	// Output: 15
}

func ExampleSliceSum_floats() {
	prices := Slice[float64]{19.99, 29.99, 9.99}
	total := SliceSum(prices)
	fmt.Printf("%.2f\n", total)
	// Output: 59.97
}
