package set_test

import (
	"fmt"

	"github.com/emad-elsaid/go-set"
)

func Example() {
	s := set.New(1, 2, 3, 3, 4, 5)

	// Set automatically deduplicates
	fmt.Println(s.Size()) // 5

	// Add elements
	s.Add(6)
	s.Add(3) // returns false, already exists

	// Set operations
	other := set.New(4, 5, 6, 7)
	union := s.Union(other)
	intersection := s.Intersection(other)

	fmt.Printf("Union size: %d\n", union.Size())
	fmt.Printf("Intersection size: %d\n", intersection.Size())

	// Functional operations
	filtered := s.Filter(func(x int) bool {
		return x > 3
	})

	fmt.Print(filtered)

	// Output:
	// 5
	// Union size: 7
	// Intersection size: 3
	// Set{4, 5, 6}
}

func ExampleNew() {
	// Create a set from integers
	s := set.New(1, 2, 3, 4, 5)
	fmt.Println(s.Size())

	// Create an empty set
	empty := set.New[string]()
	fmt.Println(empty.IsEmpty())

	// Output:
	// 5
	// true
}

func ExampleSet_Filter() {
	s := set.New(1, 2, 3, 4, 5, 6)

	// Filter even numbers
	evens := s.Filter(func(x int) bool {
		return x%2 == 0
	})

	fmt.Print(evens)

	// Output:
	// Set{2, 4, 6}
}

func ExampleSet_Union() {
	s1 := set.New(1, 2, 3)
	s2 := set.New(3, 4, 5)

	union := s1.Union(s2)
	fmt.Print(union)

	// Output:
	// Set{1, 2, 3, 4, 5}
}

func ExampleSet_Intersection() {
	s1 := set.New(1, 2, 3, 4)
	s2 := set.New(3, 4, 5, 6)

	intersection := s1.Intersection(s2)
	fmt.Print(intersection)

	// Output:
	// Set{3, 4}
}

func ExampleMap() {
	s := set.New(1, 2, 3)

	// Transform int set to string set
	stringSet := set.Map(s, func(x int) string {
		return fmt.Sprintf("num_%d", x)
	})

	fmt.Print(stringSet)

	// Output:
	// Set{num_1, num_2, num_3}
}

func ExampleSet_Each() {
	s := set.New(5, 2, 8, 1, 9, 3)

	// Iterate in insertion order
	fmt.Print("Elements in insertion order: ")
	s.Each(func(item int) {
		fmt.Printf("%d ", item)
	})

	// Output:
	// Elements in insertion order: 5 2 8 1 9 3
}

func ExampleReduce() {
	s := set.New(1, 2, 3, 4, 5)

	// Sum all elements
	sum := set.Reduce(s, 0, func(acc, x int) int {
		return acc + x
	})

	fmt.Println(sum)

	// Output:
	// 15
}

func ExampleSet_Partition() {
	s := set.New(1, 2, 3, 4, 5, 6)

	// Partition into even and odd
	evens, odds := s.Partition(func(x int) bool {
		return x%2 == 0
	})

	fmt.Printf("Evens: %v\n", evens)
	fmt.Printf("Odds: %v\n", odds)

	// Output:
	// Evens: Set{2, 4, 6}
	// Odds: Set{1, 3, 5}
}
