# CLAUDE.md

This file was created and committed by Claude AI assistant.

## About This Repository

**types** is a Go library that provides generic data structures with a Ruby-inspired API. The library leverages Go's generics (introduced in Go 1.18+) to provide type-safe, reusable implementations of common data structures with functional programming methods.

[![Go Report Card](https://goreportcard.com/badge/github.com/emad-elsaid/types)](https://goreportcard.com/report/github.com/emad-elsaid/types)
[![GoDoc](https://godoc.org/github.com/emad-elsaid/types?status.svg)](https://godoc.org/github.com/emad-elsaid/types)
[![codecov](https://codecov.io/gh/emad-elsaid/types/branch/master/graph/badge.svg)](https://codecov.io/gh/emad-elsaid/types)

## Project Structure

```
types/
├── .github/workflows/go.yml  # CI/CD pipeline
├── set.go                     # Generic Set implementation
├── set_test.go                # Set tests
├── slice.go                   # Generic Slice implementation
├── slice_test.go              # Slice tests
├── map.go                     # Generic Map implementation (thread-safe)
├── map_test.go                # Map tests
├── example_slice_test.go      # Example usage
├── go.mod                     # Go module definition
└── README.md                  # Documentation
```

## Data Structures

### 1. **Set** (`set.go`)
A generic set data structure that stores unique elements while maintaining insertion order.

**Key Features:**
- Order-preserving (elements maintain insertion order)
- Standard set operations: Union, Intersection, Difference, SymmetricDifference
- Set comparisons: IsSubset, IsSuperset, IsDisjoint, Equal
- Functional methods: Map, Filter, Reject, Find, All, Any, None, Count
- Advanced operations: Partition, Take, Drop, Reduce
- Thread-safe operations available through proper usage

**Methods:** Add, Remove, Contains, Size, IsEmpty, Clear, ToSlice, Clone, Each, and more

### 2. **Slice** (`slice.go`)
A Ruby-inspired generic slice with extensive functional programming methods.

**Key Features:**
- Rich set of manipulation methods (50+ methods)
- Negative index support (Ruby-style)
- Functional operations: Map, Filter (KeepIf), Reduce, Select
- Array operations: Push, Pop, Shift, Unshift, Insert, Delete
- Query operations: All, Any, Find, Count, Include
- Transformation: Reverse, Shuffle, Unique, Partition
- Utilities: Cycle, Fill, FillWith, At, Fetch

**Example:**
```go
a := Slice[int]{1, 2, 3, 4, 5, 6}
a = a.Map(func(e int) int { return e * 100 })
a = a.KeepIf(func(e int) bool { return e <= 300 })
// Result: [100 200 300]
```

### 3. **Map** (`map.go`)
A generic wrapper around `sync.Map` providing thread-safe key-value storage with type safety.

**Key Features:**
- Thread-safe concurrent access
- Generic type parameters for keys and values
- Methods: Store, Load, Delete, Range, LoadAndDelete, LoadOrStore, Swap
- Based on Go's standard `sync.Map` implementation

## Technical Details

- **Go Version:** 1.22+
- **Module Path:** `github.com/emad-elsaid/types`
- **License:** MIT (see LICENSE file)
- **Testing:** Comprehensive test coverage with examples
- **CI/CD:** GitHub Actions workflow for automated testing on push/PR

## Design Philosophy

This library follows Ruby's design philosophy of providing expressive, chainable methods that make code more readable and concise. All types use Go's comparable constraint where appropriate, ensuring type safety while maintaining flexibility.

The library emphasizes:
- **Immutability patterns:** Many methods return new instances rather than modifying in place
- **Functional programming:** Rich set of higher-order functions (Map, Filter, Reduce, etc.)
- **Ruby compatibility:** Method names and behaviors inspired by Ruby's Array, Set, and Hash
- **Type safety:** Full use of Go generics for compile-time type checking

## Testing & Quality

The project includes:
- Extensive unit tests for all data structures
- Example-based tests demonstrating usage
- Code coverage reporting via codecov
- Continuous integration via GitHub Actions
- Go Report Card monitoring for code quality

## Initialized with Claude Code

This repository has been initialized for use with Claude Code, Anthropic's official CLI for Claude.
