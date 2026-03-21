package types_test

import (
	"fmt"

	"github.com/emad-elsaid/types"
)

func ExampleSlice_ChunkWhile() {
	// Group consecutive numbers
	numbers := types.Slice[int]{1, 2, 4, 5, 7, 9}
	chunks := numbers.ChunkWhile(func(a, b int) bool {
		return b-a == 1
	})
	fmt.Println(chunks)

	// Output:
	// [[1 2] [4 5] [7] [9]]
}

func ExampleSlice_ChunkWhile_strings() {
	// Group strings by same length
	words := types.Slice[string]{"a", "b", "cc", "dd", "eee", "f"}
	chunks := words.ChunkWhile(func(a, b string) bool {
		return len(a) == len(b)
	})
	fmt.Println(chunks)

	// Output:
	// [[a b] [cc dd] [eee] [f]]
}
