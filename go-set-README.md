# go-set

[![Go Reference](https://pkg.go.dev/badge/github.com/emad-elsaid/go-set.svg)](https://pkg.go.dev/github.com/emad-elsaid/go-set)
[![Go Report Card](https://goreportcard.com/badge/github.com/emad-elsaid/go-set)](https://goreportcard.com/report/github.com/emad-elsaid/go-set)

An **order-preserving** generic Set data structure for Go.

## Why Another Set Library?

Most Go set implementations are either:
- Pre-generics (using `interface{}` or code generation)
- Unordered (iteration order is random)

This library provides a **generic, order-preserving set** that maintains insertion order, making it unique in the Go ecosystem. This is particularly useful when you need both set semantics (uniqueness) and predictable iteration order.

## Features

- ✅ **Generic**: Works with any `comparable` type
- ✅ **Order-preserving**: Maintains insertion order during iteration
- ✅ **Full set operations**: Union, Intersection, Difference, Symmetric Difference
- ✅ **Functional methods**: Map, Filter, Reduce, Partition
- ✅ **Rich API**: 30+ methods for comprehensive set manipulation
- ✅ **Zero dependencies**: Only uses Go standard library
- ✅ **Well-tested**: Extensive test coverage

## Installation

```bash
go get github.com/emad-elsaid/go-set
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/emad-elsaid/go-set"
)

func main() {
    // Create a new set
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

    // Functional operations
    filtered := s.Filter(func(x int) bool {
        return x > 3
    })

    doubled := set.Map(s, func(x int) int {
        return x * 2
    })

    fmt.Println(filtered) // Set{4, 5, 6}
}
```

## API Overview

### Construction

```go
s := set.New(1, 2, 3)              // Create from elements
s := set.New[int]()                // Create empty set
```

### Basic Operations

```go
s.Add(item)                         // Add element, returns true if added
s.Remove(item)                      // Remove element, returns true if removed
s.Contains(item)                    // Check if element exists
s.Size()                            // Get number of elements
s.IsEmpty()                         // Check if set is empty
s.Clear()                           // Remove all elements
s.ToSlice()                         // Convert to slice (maintains order)
s.Clone()                           // Create a shallow copy
```

### Set Operations

```go
s1.Union(s2)                        // Elements in s1 OR s2
s1.Intersection(s2)                 // Elements in s1 AND s2
s1.Difference(s2)                   // Elements in s1 but not in s2
s1.SymmetricDifference(s2)          // Elements in s1 XOR s2
s1.IsSubset(s2)                     // All elements of s1 are in s2
s1.IsSuperset(s2)                   // All elements of s2 are in s1
s1.IsDisjoint(s2)                   // No common elements
s1.Equal(s2)                        // Same elements (order-independent)
```

### Functional Operations

```go
s.Each(func(x int) { ... })         // Iterate in insertion order
s.Filter(func(x int) bool { ... })  // Keep elements matching predicate
s.Reject(func(x int) bool { ... })  // Remove elements matching predicate
s.Find(func(x int) bool { ... })    // Find first matching element
s.All(func(x int) bool { ... })     // Check if all elements match
s.Any(func(x int) bool { ... })     // Check if any element matches
s.None(func(x int) bool { ... })    // Check if no elements match
s.Count(func(x int) bool { ... })   // Count matching elements
s.Partition(func(x int) bool { ... }) // Split into two sets
s.Take(n)                           // Take first n elements
s.Drop(n)                           // Drop first n elements

// Generic functions
set.Map(s, func(x T) U { ... })     // Transform elements to new set
set.Reduce(s, init, func(acc U, x T) U { ... }) // Reduce to single value
```

## Order Preservation Example

```go
s := set.New(5, 2, 8, 1, 9, 3)

var result []int
s.Each(func(item int) {
    result = append(result, item)
})

fmt.Println(result) // [5, 2, 8, 1, 9, 3] - insertion order preserved!
```

## Type Safety with Generics

```go
// Works with any comparable type
intSet := set.New(1, 2, 3)
stringSet := set.New("apple", "banana", "cherry")

type Point struct{ X, Y int }
pointSet := set.New(Point{1, 2}, Point{3, 4})

// Transform between types
stringSet := set.Map(intSet, func(n int) string {
    return fmt.Sprintf("num_%d", n)
})
```

## Performance Characteristics

| Operation | Time Complexity | Space Complexity |
|-----------|----------------|------------------|
| Add       | O(1) average   | O(1)            |
| Remove    | O(n)           | O(1)            |
| Contains  | O(1)           | O(1)            |
| Size      | O(1)           | O(1)            |
| Union     | O(n+m)         | O(n+m)          |
| Intersection | O(min(n,m)) | O(min(n,m))     |
| Difference | O(n)          | O(n)            |

Note: Remove is O(n) due to maintaining insertion order. If you don't need order preservation, consider using a standard map-based set for O(1) removal.

## Comparison with Other Libraries

| Feature | go-set | samber/lo | deckarep/golang-set |
|---------|--------|-----------|---------------------|
| Generics | ✅ | ✅ | ❌ (pre-generics) |
| Order-preserving | ✅ | N/A | ❌ |
| Set operations | ✅ | ❌ | ✅ |
| Functional methods | ✅ | ✅ (on slices) | Limited |
| Thread-safe option | ❌ | ❌ | ✅ |

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see LICENSE file for details.

## Acknowledgments

Inspired by Ruby's Set class and the need for order-preserving sets in Go.
