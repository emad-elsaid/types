package types

import (
	"fmt"
	"sort"
)

func ExampleSlice_Tally() {
	// Count word frequencies in a slice
	words := Slice[string]{"apple", "banana", "apple", "cherry", "banana", "apple", "date"}
	frequencies := words.Tally()

	// Sort keys for consistent output
	keys := make([]string, 0, len(frequencies))
	for k := range frequencies {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, word := range keys {
		fmt.Printf("%s: %d\n", word, frequencies[word])
	}

	// Output:
	// apple: 3
	// banana: 2
	// cherry: 1
	// date: 1
}

func ExampleSlice_Tally_numbers() {
	// Find frequency distribution of dice rolls
	rolls := Slice[int]{1, 6, 3, 6, 2, 6, 5, 1, 6, 4}
	distribution := rolls.Tally()

	// Find the most common roll
	maxCount := 0
	mostCommon := 0
	for roll, count := range distribution {
		if count > maxCount {
			maxCount = count
			mostCommon = roll
		}
	}

	fmt.Printf("Most common roll: %d (appeared %d times)\n", mostCommon, maxCount)

	// Output:
	// Most common roll: 6 (appeared 4 times)
}
