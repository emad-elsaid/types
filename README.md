# Types

Go implementation for generic types, imitating Ruby types

# Data structures and methods

* func (a Array) All(block func(interface{}) bool) bool
* func (a Array) Any(block func(interface{}) bool) bool
* func (a Array) At(index int64) interface{}
* func (a Array) Count() int64
* func (a Array) CountBy(block func(interface{}) bool) (count int64)
* func (a Array) CountElement(element interface{}) (count int64)
* func (a Array) Cycle(count int64, block func(interface{}))
