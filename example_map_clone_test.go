package types

import "fmt"

func ExampleMap_Clone() {
	// Create original map with some entries
	original := &Map[string, int]{}
	original.Store("apples", 5)
	original.Store("bananas", 3)
	original.Store("oranges", 7)

	// Clone the map
	clone := original.Clone()

	// Modify the clone
	clone.Store("apples", 10)     // Update existing
	clone.Store("grapes", 15)     // Add new entry
	clone.Delete("bananas")       // Remove entry

	// Original remains unchanged
	apples, _ := original.Load("apples")
	fmt.Println("Original apples:", apples)
	fmt.Println("Original has bananas:", original.Has("bananas"))
	fmt.Println("Original has grapes:", original.Has("grapes"))

	// Clone has the modifications
	cApples, _ := clone.Load("apples")
	fmt.Println("Clone apples:", cApples)
	fmt.Println("Clone has bananas:", clone.Has("bananas"))
	fmt.Println("Clone has grapes:", clone.Has("grapes"))

	// Output:
	// Original apples: 5
	// Original has bananas: true
	// Original has grapes: false
	// Clone apples: 10
	// Clone has bananas: false
	// Clone has grapes: true
}
