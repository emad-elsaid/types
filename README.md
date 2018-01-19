# Types

Go implementation for generic types, imitating Ruby types

# Data structures and methods

## Array

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
