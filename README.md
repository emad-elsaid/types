# Types

Go implementation for generic types, imitating Ruby types

# Data structures and methods

* func (this Array) At(index int64) interface{}
* func (this Array) Count() int64
* func (this Array) CountBy(block func(interface{}) bool) (count int64)
* func (this Array) CountElement(element interface{}) (count int64)
* func (this Array) Cycle(count int64, block func(interface{}))
