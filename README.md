# Types

Go implementation for generic types, imitating Ruby types

# Data structures and methods

## Array

```
All(block func(interface{}) bool) bool
Any(block func(interface{}) bool) bool
At(index int) interface{}
Collect(block func(interface{}) interface{}) Array
Compact() Array
Count() int
CountBy(block func(interface{}) bool) (count int)
CountElement(element interface{}) (count int)
Cycle(count int, block func(interface{}))
Delete(element interface{}) Array
DeleteAt(index int) Array
DeleteIf(block func(interface{}) bool) Array
Drop(count int) Array
Each(block func(interface{}))
EachIndex(block func(int))
Fetch(index int, defaultValue interface{}) interface{}
Fill(element interface{}, start int, length int) Array
FillWith(start int, length int, block func(int) interface{}) Array
First() interface{}
Firsts(count int) Array
Flatten() Array
Include(element interface{}) bool
Index(element interface{}) int
IndexBy(block func(interface{}) bool) int
Insert(index int, elements ...interface{}) Array
IsEmpty() bool
IsEq(other Array) bool
KeepIf(block func(interface{}) bool) Array
Last() interface{}
Lasts(count int) Array
Len() int
Map(block func(interface{}) interface{}) Array
Max(block func(interface{}) int) interface{}
Min(block func(interface{}) int) interface{}
```
