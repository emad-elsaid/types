// Package types provides a generic data types similar to that of Ruby
package types

// Array is an alias for a slice of variables that vary in types
// Array can hold any data type, it also allow different types
// in the same Array
type Array []interface{}

// Array At returns element by index, a negative index counts from the end of
// the Array
// if index is out of range it returns nil
func (a Array) At(index int64) interface{} {
	len := int64(len(a))

	if index < 0 {
		if -index <= len {
			return a[len+index]
		}
		return nil
	}

	if index < len {
		return a[index]
	}

	return nil
}

// Array Count returns total number of elements in Array
func (a Array) Count() int64 {
	return int64(len(a))
}

// Array CountElement returns number of elements equal to "element" in Array
func (a Array) CountElement(element interface{}) (count int64) {
	for _, o := range a {
		if o == element {
			count++
		}
	}
	return count
}

// Array CountBy returns number of elements which "block" returns true for
func (a Array) CountBy(block func(interface{}) bool) (count int64) {
	for _, o := range a {
		if block(o) {
			count++
		}
	}
	return count
}

// Array Cycle will cycle through Array elements "count" times passing each
// element to "block" function
func (a Array) Cycle(count int64, block func(interface{})) {
	var i int64
	for i = 0; i < count; i++ {
		for _, v := range a {
			block(v)
		}
	}
}

// Array Any returns true if "block" returned true for any of the Array elements
// and false otherwise
func (a Array) Any(block func(interface{}) bool) bool {
	for _, o := range a {
		if block(o) {
			return true
		}
	}

	return false
}

// Array All returns true if "block" returned true for all elements in Array and
// false otherwise
func (a Array) All(block func(interface{}) bool) bool {
	for _, o := range a {
		if !block(o) {
			return false
		}
	}

	return true
}

// Array Collect will pass every element in array to "block" returining a new Array with the return values
func (a Array) Collect(block func(interface{}) interface{}) Array {
	result := Array{}

	for _, o := range a {
		result = append(result, block(o))
	}

	return result
}

// Array Compact will return a new array with all non-nil elements
func (a Array) Compact() Array {
	result := Array{}
	for _, o := range a {
		if o != nil {
			result = append(result, o)
		}
	}
	return result
}

// Array Delete will remove all elements that are equal to the passed element
func (a Array) Delete(element interface{}) Array {
	result := Array{}
	for _, o := range a {
		if o != element {
			result = append(result, o)
		}
	}
	return result
}

// Array DeleteAt will delete an element by index
func (a Array) DeleteAt(index int64) Array {
	return append(a[:index], a[index+1:]...)
}

// Array DeleteIf will delete all elements which "block" returns true for
func (a Array) DeleteIf(block func(interface{}) bool) Array {
	result := Array{}
	for _, o := range a {
		if !block(o) {
			result = append(result, o)
		}
	}
	return result
}

// Array Drop will return an array without the first "count" elements from the beginning
func (a Array) Drop(count int64) Array {
	return a[count:]
}

// Array Each will execute "block" for each element in array
func (a Array) Each(block func(interface{})) {
	for _, o := range a {
		block(o)
	}
}

// Array EachIndex will execute "block" for each element index in array
func (a Array) EachIndex(block func(int64)) {
	for i, _ := range a {
		block(int64(i))
	}
}

// Array IsEmpty will return true of array empty, false otherwise
func (a Array) IsEmpty() bool {
	return len(a) == 0
}
