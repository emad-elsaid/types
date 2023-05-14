# Types

[![Go Report Card](https://goreportcard.com/badge/github.com/emad-elsaid/types)](https://goreportcard.com/report/github.com/emad-elsaid/types)
[![GoDoc](https://godoc.org/github.com/emad-elsaid/types?status.svg)](https://godoc.org/github.com/emad-elsaid/types)
[![codecov](https://codecov.io/gh/emad-elsaid/types/branch/master/graph/badge.svg)](https://codecov.io/gh/emad-elsaid/types)

Go implementation for generic types, imitating Ruby types

# Data structures

## Slice

A slice of `comparable` type that can hold any value

### Example

```go
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
```

### Slice Methods available:

```go
func (a Slice[T]) All(block func(T) bool) bool
func (a Slice[T]) Any(block func(T) bool) bool
func (a Slice[T]) At(index int) *T
func (a Slice[T]) CountBy(block func(T) bool) (count int)
func (a Slice[T]) CountElement(element T) (count int)
func (a Slice[T]) Cycle(count int, block func(T))
func (a Slice[T]) Delete(element T) Slice[T]
func (a Slice[T]) DeleteAt(index int) Slice[T]
func (a Slice[T]) DeleteIf(block func(T) bool) Slice[T]
func (a Slice[T]) Drop(count int) Slice[T]
func (a Slice[T]) Each(block func(T))
func (a Slice[T]) EachIndex(block func(int))
func (a Slice[T]) Fetch(index int, defaultValue T) T
func (a Slice[T]) Fill(element T, start int, length int) Slice[T]
func (a Slice[T]) FillWith(start int, length int, block func(int) T) Slice[T]
func (a Slice[T]) First() *T
func (a Slice[T]) Firsts(count int) Slice[T]
func (a Slice[T]) Include(element T) bool
func (a Slice[T]) Index(element T) int
func (a Slice[T]) IndexBy(block func(T) bool) int
func (a Slice[T]) Insert(index int, elements ...T) Slice[T]
func (a Slice[T]) IsEmpty() bool
func (a Slice[T]) IsEq(other Slice[T]) bool
func (a Slice[T]) KeepIf(block func(T) bool) Slice[T]
func (a Slice[T]) Last() *T
func (a Slice[T]) Lasts(count int) Slice[T]
func (a Slice[T]) Len() int
func (a Slice[T]) Map(block func(T) T) Slice[T]
func (a Slice[T]) Max(block func(T) int) T
func (a Slice[T]) Min(block func(T) int) T
func (a Slice[T]) Pop() (Slice[T], T)
func (a Slice[T]) Push(element T) Slice[T]
func (a Slice[T]) Reduce(block func(T) bool) Slice[T]
func (a Slice[T]) Reverse() Slice[T]
func (a Slice[T]) Select(block func(T) bool) Slice[T]
func (a Slice[T]) SelectUntil(block func(T) bool) Slice[T]
func (a Slice[T]) Shift() (T, Slice[T])
func (a Slice[T]) Shuffle() Slice[T]
func (a Slice[T]) Unshift(element T) Slice[T]
```
