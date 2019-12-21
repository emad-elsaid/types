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

```go
func (a Array) All(block func(Element) bool) bool
func (a Array) Any(block func(Element) bool) bool
func (a Array) At(index int) Element
func (a Array) Compact() Array
func (a Array) CountBy(block func(Element) bool) (count int)
func (a Array) CountElement(element Element) (count int)
func (a Array) Cycle(count int, block func(Element))
func (a Array) Delete(element Element) Array
func (a Array) DeleteAt(index int) Array
func (a Array) DeleteIf(block func(Element) bool) Array
func (a Array) Drop(count int) Array
func (a Array) Each(block func(Element))
func (a Array) EachIndex(block func(int))
func (a Array) Fetch(index int, defaultValue Element) Element
func (a Array) Fill(element Element, start int, length int) Array
func (a Array) FillWith(start int, length int, block func(int) Element) Array
func (a Array) First() Element
func (a Array) Firsts(count int) Array
func (a Array) Flatten() Array
func (a Array) Include(element Element) bool
func (a Array) Index(element Element) int
func (a Array) IndexBy(block func(Element) bool) int
func (a Array) Insert(index int, elements ...Element) Array
func (a Array) IsEmpty() bool
func (a Array) IsEq(other Array) bool
func (a Array) KeepIf(block func(Element) bool) Array
func (a Array) Last() Element
func (a Array) Lasts(count int) Array
func (a Array) Len() int
func (a Array) Map(block func(Element) Element) Array
func (a Array) Max(block func(Element) int) Element
func (a Array) Min(block func(Element) int) Element
func (a Array) Pop() (Array, Element)
func (a Array) Push(element Element) Array
func (a Array) Reverse() Array
func (a Array) Shift() (Element, Array)
func (a Array) Shuffle() Array
func (a Array) Unshift(element Element) Array
func (a Array) Reduce(block func(Element) bool) Array
func (a Array) Select(block func(Element) bool) Array
```

## CLI generator for your type

To install the `array` command, which generate the type specific array code

```
go get github.com/emad-elsaid/types/cmd/array
```

## Generating the methods for your type

Lets say you have a type `Record`, and you want to generate a type that
represent an array of records `Records` that has all of the previous methods,
then write in in a file called `records.go` you'll need to execute the following

```
array -element Record -array Records -package db -output records.go
```

The usage for `array` is as follows:

```
Usage of ./array:
  -array string
        the name of the slice of your element (default "stringArray")
  -element string
        the single element of your array (default "string")
  -output string
        where to write the output (default "/dev/stdout")
  -package string
        package name the new file will belong to (default "main")
```
