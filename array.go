// Package types provides a generic data types similar to that of Ruby
package types

// Array is an alias for a slice of variables that vary in types
// Array can hold any data type, it also allow different types
// in the same Array
type Array []interface{}

// Array At returns element by index, a negative index counts from the end of
// the Array
// if index is out of range it returns nil
func (a Array) At(index int) interface{} {
	len := len(a)

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
func (a Array) Count() int {
	return len(a)
}

// Array CountElement returns number of elements equal to "element" in Array
func (a Array) CountElement(element interface{}) (count int) {
	for _, o := range a {
		if o == element {
			count++
		}
	}
	return count
}

// Array CountBy returns number of elements which "block" returns true for
func (a Array) CountBy(block func(interface{}) bool) (count int) {
	for _, o := range a {
		if block(o) {
			count++
		}
	}
	return count
}

// Array Cycle will cycle through Array elements "count" times passing each
// element to "block" function
func (a Array) Cycle(count int, block func(interface{})) {
	for i := 0; i < count; i++ {
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
func (a Array) DeleteAt(index int) Array {
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

// Array Drop will return an array without the first "count" elements from the
// beginning
func (a Array) Drop(count int) Array {
	return a[count:]
}

// Array Each will execute "block" for each element in array
func (a Array) Each(block func(interface{})) {
	for _, o := range a {
		block(o)
	}
}

// Array EachIndex will execute "block" for each element index in array
func (a Array) EachIndex(block func(int)) {
	for i, _ := range a {
		block(i)
	}
}

// Array IsEmpty will return true of array empty, false otherwise
func (a Array) IsEmpty() bool {
	return len(a) == 0
}

// Array IsEq returns true if array the "other" array
func (a Array) IsEq(other Array) bool {
	// check length
	if len(a) != len(other) {
		return false
	}

	// check values
	for i, o := range a {
		if o != other[i] {
			return false
		}
	}

	return true
}

// Array Len returns number of elements in array
func (a Array) Len() int {
	return len(a)
}

// Array Fetch will return the element in "index", if it doesn't exist it
// returns the passed "defaultValue"
func (a Array) Fetch(index int, defaultValue interface{}) interface{} {
	val := a.At(index)
	if val != nil {
		return val
	}

	return defaultValue
}

// Array Fill will replace elements inplace starting from "start" counting
// "length" elements with the passed "element" parameter, will return same array
// object
func (a Array) Fill(element interface{}, start int, length int) Array {
	for length--; length >= 0; length-- {
		a[start+length] = element
	}
	return a
}

// Array FillWith will replace elements from start counting "length" items,
// passing every index to block and replacing the element inplace with the
// return value
func (a Array) FillWith(start int, length int, block func(int) interface{}) Array {
	for length--; length >= 0; length-- {
		a[start+length] = block(start + length)
	}
	return a
}

// Array Index returns the index of the first element in array that is equal to
// "element", returns -1 if the elements if not found
func (a Array) Index(element interface{}) int {
	for i, o := range a {
		if o == element {
			return i
		}
	}
	return -1
}

// Array IndexBy returns first element that block returns true for, -1 otherwise
func (a Array) IndexBy(block func(interface{}) bool) int {
	for i, o := range a {
		if block(o) {
			return i
		}
	}

	return -1
}

// Array First returns first element of array
func (a Array) First() interface{} {
	return a.At(0)
}

// Array Last returns last element of array
func (a Array) Last() interface{} {
	return a.At(len(a) - 1)
}

// Array Firsts will return an array holding the first "count" elements of the
// array
func (a Array) Firsts(count int) Array {
	return a[0:count]
}

// Array Lasts will return an array holding the lasts "count" elements of the
// array
func (a Array) Lasts(count int) Array {
	return a[len(a)-count:]
}

// Array Flatten returns a flattened array of the current one, expanding any
// element that could be casted to Array recursively until no element could be flattened
func (a Array) Flatten() Array {
	result := Array{}
	for _, o := range a {
		element, ok := o.(Array)
		if ok {
			result = append(result, element.Flatten()...)
		} else {
			result = append(result, o)
		}
	}
	return result
}

// Array Include will return true if element found in the array
func (a Array) Include(element interface{}) bool {
	return a.Index(element) != -1
}

// Array Insert will insert a set of elements in the index and will return a new
// array
func (a Array) Insert(index int, elements ...interface{}) Array {
	result := Array{}
	result = append(result, a[0:index]...)
	result = append(result, elements...)
	result = append(result, a[index:]...)
	return result
}

// Array KeepIf will return an array contains all elements where "block"
// returned true for them
func (a Array) KeepIf(block func(interface{}) bool) Array {
	result := Array{}
	for _, o := range a {
		if block(o) {
			result = append(result, o)
		}
	}
	return result
}

// Array Map will return a new array replacing every element from current array
// with the return value of the block
func (a Array) Map(block func(interface{}) interface{}) Array {
	result := Array{}
	for _, o := range a {
		result = append(result, block(o))
	}
	return result
}

// Array Max returns the element the returned the highest value when passed to
// block
func (a Array) Max(block func(interface{}) int) interface{} {
	if len(a) == 0 {
		return nil
	}
	var maxElement interface{} = a[0]
	var maxScore = block(a[0])
	for _, o := range a[1:] {
		score := block(o)
		if score > maxScore {
			maxElement = o
			maxScore = score
		}
	}

	return maxElement
}

// Array Min returns the element the returned the lowest value when passed to
// block
func (a Array) Min(block func(interface{}) int) interface{} {
	if len(a) == 0 {
		return nil
	}
	var minElement interface{} = a[0]
	var minScore = block(a[0])
	for _, o := range a[1:] {
		score := block(o)
		if score < minScore {
			minElement = o
			minScore = score
		}
	}

	return minElement
}
