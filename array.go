// Package types provides a generic data types similar to that of Ruby
package types

// Array is an alias for a slice of variables that vary in types
// Array can hold any data type, it also allow different types
// in the same Array
type Array []interface{}

// Array At returns element by index, a negative index counts from the end of
// the Array
// if index is out of range it returns nil
func (this Array) At(index int64) interface{} {
	len := int64(len(this))

	if index < 0 {
		if -index <= len {
			return this[len+index]
		}
		return nil
	}

	if index < len {
		return this[index]
	}

	return nil
}

// Array Count returns total number of elements in Array
func (this Array) Count() int64 {
	return int64(len(this))
}

// Array CountElement returns number of elements equal to "element" in Array
func (this Array) CountElement(element interface{}) (count int64) {
	for _, o := range this {
		if o == element {
			count++
		}
	}
	return count
}

// Array CountBy returns number of elements which "block" returns true for
func (this Array) CountBy(block func(interface{}) bool) (count int64) {
	for _, o := range this {
		if block(o) {
			count++
		}
	}
	return count
}

// Array Cycle will cycle through Array elements "count" times passing each
// element to "block" function
func (this Array) Cycle(count int64, block func(interface{})) {
	var i int64
	for i = 0; i < count; i++ {
		for _, v := range this {
			block(v)
		}
	}
}
