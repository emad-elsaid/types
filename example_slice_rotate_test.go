package types

import "fmt"

func ExampleSlice_Rotate() {
	nums := Slice[int]{1, 2, 3, 4, 5}

	// Rotate left by 2 positions
	left := nums.Rotate(2)
	fmt.Println(left)

	// Rotate right by 2 positions
	right := nums.Rotate(-2)
	fmt.Println(right)

	// Rotate by more than length
	bigRotate := nums.Rotate(7) // equivalent to rotating by 2
	fmt.Println(bigRotate)

	// Output:
	// [3 4 5 1 2]
	// [4 5 1 2 3]
	// [3 4 5 1 2]
}
