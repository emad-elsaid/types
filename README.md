# Types

Go implementation for generic types, imitating Ruby types

# Data structures and methods

## Array

```
func (a Array) All(block func(interface{}) bool) bool
func (a Array) Any(block func(interface{}) bool) bool
func (a Array) At(index int64) interface{}
func (a Array) Collect(block func(interface{}) interface{}) Array
func (a Array) Compact() Array
func (a Array) Count() int64
func (a Array) CountBy(block func(interface{}) bool) (count int64)
func (a Array) CountElement(element interface{}) (count int64)
func (a Array) Cycle(count int64, block func(interface{}))
func (a Array) Delete(element interface{}) Array
func (a Array) DeleteAt(index int64) Array
```
