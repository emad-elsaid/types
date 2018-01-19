# Types

[![Go Report Card](https://goreportcard.com/badge/github.com/emad-elsaid/types)](https://goreportcard.com/report/github.com/emad-elsaid/types)
[![GoDoc](https://godoc.org/github.com/emad-elsaid/types?status.svg)](https://godoc.org/github.com/emad-elsaid/types)
[![CircleCI](https://circleci.com/gh/emad-elsaid/types.svg?style=shield)](https://circleci.com/gh/emad-elsaid/types)
[![codecov](https://codecov.io/gh/emad-elsaid/types/branch/master/graph/badge.svg)](https://codecov.io/gh/emad-elsaid/types)

Go implementation for generic types, imitating Ruby types

# Data structures

## Array

A slice of generic type `Element` that can hold any value, even of different types.

### Example

```go
a := Array{1, 2, 3, 4, 5, 6}

// multiply every element by 100
a = a.Map(func(e Element) Element {
    return e.(int) * 100
})

// select any element <= 300
a = a.KeepIf(func(e Element) bool {
    return e.(int) <= 300
})

fmt.Print(a)
// Output: [100 200 300]
```

### Array Methods available:

```
All(block func(Element) bool) bool
Any(block func(Element) bool) bool
At(index int) Element
Collect(block func(Element) Element) Array
Compact() Array
Count() int
CountBy(block func(Element) bool) (count int)
CountElement(element Element) (count int)
Cycle(count int, block func(Element))
Delete(element Element) Array
DeleteAt(index int) Array
DeleteIf(block func(Element) bool) Array
Drop(count int) Array
Each(block func(Element))
EachIndex(block func(int))
Fetch(index int, defaultValue Element) Element
Fill(element Element, start int, length int) Array
FillWith(start int, length int, block func(int) Element) Array
First() Element
Firsts(count int) Array
Flatten() Array
Include(element Element) bool
Index(element Element) int
IndexBy(block func(Element) bool) int
Insert(index int, elements ...Element) Array
IsEmpty() bool
IsEq(other Array) bool
KeepIf(block func(Element) bool) Array
Last() Element
Lasts(count int) Array
Len() int
Map(block func(Element) Element) Array
Max(block func(Element) int) Element
Min(block func(Element) int) Element
Pop() (Array, Element)
Push(element Element) Array
Reverse() Array
Shift() (Element, Array)
Shuffle() Array
Unshift(element Element) Array
```
