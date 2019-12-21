// Package types provides a generic data types similar to that of Ruby
package types

import "math/rand"

// ElementArray is an alias for a slice of variables that vary in types
// ElementArray can hold any data type, it also allow different types
// in the same ElementArray
type Element interface{}
type ElementArray []Element

// At returns element by index, a negative index counts from the end of
// the ElementArray
// if index is out of range it returns nil
func (a ElementArray) At(index int) Element {
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

// CountElement returns number of elements equal to "element" in ElementArray
func (a ElementArray) CountElement(element Element) (count int) {
	for _, o := range a {
		if o == element {
			count++
		}
	}
	return count
}

// CountBy returns number of elements which "block" returns true for
func (a ElementArray) CountBy(block func(Element) bool) (count int) {
	for _, o := range a {
		if block(o) {
			count++
		}
	}
	return count
}

// Cycle will cycle through ElementArray elements "count" times passing each
// element to "block" function
func (a ElementArray) Cycle(count int, block func(Element)) {
	for i := 0; i < count; i++ {
		for _, v := range a {
			block(v)
		}
	}
}

// Any returns true if "block" returned true for any of the ElementArray elements
// and false otherwise
func (a ElementArray) Any(block func(Element) bool) bool {
	for _, o := range a {
		if block(o) {
			return true
		}
	}

	return false
}

// All returns true if "block" returned true for all elements in ElementArray and
// false otherwise
func (a ElementArray) All(block func(Element) bool) bool {
	for _, o := range a {
		if !block(o) {
			return false
		}
	}

	return true
}

// Compact will return a new array with all non-nil elements
func (a ElementArray) Compact() ElementArray {
	result := ElementArray{}
	for _, o := range a {
		if o != nil {
			result = append(result, o)
		}
	}
	return result
}

// Delete will remove all elements that are equal to the passed element
func (a ElementArray) Delete(element Element) ElementArray {
	result := ElementArray{}
	for _, o := range a {
		if o != element {
			result = append(result, o)
		}
	}
	return result
}

// DeleteAt will delete an element by index
func (a ElementArray) DeleteAt(index int) ElementArray {
	return append(a[:index], a[index+1:]...)
}

// DeleteIf will delete all elements which "block" returns true for
func (a ElementArray) DeleteIf(block func(Element) bool) ElementArray {
	result := ElementArray{}
	for _, o := range a {
		if !block(o) {
			result = append(result, o)
		}
	}
	return result
}

// Drop will return an array without the first "count" elements from the
// beginning
func (a ElementArray) Drop(count int) ElementArray {
	return a[count:]
}

// Each will execute "block" for each element in array
func (a ElementArray) Each(block func(Element)) {
	for _, o := range a {
		block(o)
	}
}

// EachIndex will execute "block" for each element index in array
func (a ElementArray) EachIndex(block func(int)) {
	for i := range a {
		block(i)
	}
}

// IsEmpty will return true of array empty, false otherwise
func (a ElementArray) IsEmpty() bool {
	return len(a) == 0
}

// IsEq returns true if array the "other" array
func (a ElementArray) IsEq(other ElementArray) bool {
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

// Len returns number of elements in array
func (a ElementArray) Len() int {
	return len(a)
}

// Fetch will return the element in "index", if it doesn't exist it
// returns the passed "defaultValue"
func (a ElementArray) Fetch(index int, defaultValue Element) Element {
	val := a.At(index)
	if val != nil {
		return val
	}

	return defaultValue
}

// Fill will replace elements inplace starting from "start" counting
// "length" elements with the passed "element" parameter, will return same array
// object
func (a ElementArray) Fill(element Element, start int, length int) ElementArray {
	for length--; length >= 0; length-- {
		a[start+length] = element
	}
	return a
}

// FillWith will replace elements from start counting "length" items,
// passing every index to block and replacing the element inplace with the
// return value
func (a ElementArray) FillWith(start int, length int, block func(int) Element) ElementArray {
	for length--; length >= 0; length-- {
		a[start+length] = block(start + length)
	}
	return a
}

// Index returns the index of the first element in array that is equal to
// "element", returns -1 if the elements if not found
func (a ElementArray) Index(element Element) int {
	for i, o := range a {
		if o == element {
			return i
		}
	}
	return -1
}

// IndexBy returns first element that block returns true for, -1 otherwise
func (a ElementArray) IndexBy(block func(Element) bool) int {
	for i, o := range a {
		if block(o) {
			return i
		}
	}

	return -1
}

// First returns first element of array
func (a ElementArray) First() Element {
	return a.At(0)
}

// Last returns last element of array
func (a ElementArray) Last() Element {
	return a.At(len(a) - 1)
}

// Firsts will return an array holding the first "count" elements of the
// array
func (a ElementArray) Firsts(count int) ElementArray {
	return a[0:count]
}

// Lasts will return an array holding the lasts "count" elements of the
// array
func (a ElementArray) Lasts(count int) ElementArray {
	return a[len(a)-count:]
}

// Flatten returns a flattened array of the current one, expanding any
// element that could be casted to ElementArray recursively until no element could be flattened
func (a ElementArray) Flatten() ElementArray {
	result := ElementArray{}
	for _, o := range a {
		element, ok := o.(ElementArray)
		if ok {
			result = append(result, element.Flatten()...)
		} else {
			result = append(result, o)
		}
	}
	return result
}

// Include will return true if element found in the array
func (a ElementArray) Include(element Element) bool {
	return a.Index(element) != -1
}

// Insert will insert a set of elements in the index and will return a new
// array
func (a ElementArray) Insert(index int, elements ...Element) ElementArray {
	result := ElementArray{}
	result = append(result, a[0:index]...)
	result = append(result, elements...)
	result = append(result, a[index:]...)
	return result
}

// KeepIf will return an array contains all elements where "block"
// returned true for them
func (a ElementArray) KeepIf(block func(Element) bool) ElementArray {
	result := ElementArray{}
	for _, o := range a {
		if block(o) {
			result = append(result, o)
		}
	}
	return result
}

// Select is an alias for KeepIf
func (a ElementArray) Select(block func(Element) bool) ElementArray {
	return a.KeepIf(block)
}

// Reduce is an alias for KeepIf
func (a ElementArray) Reduce(block func(Element) bool) ElementArray {
	return a.KeepIf(block)
}

// Map will return a new array replacing every element from current array
// with the return value of the block
func (a ElementArray) Map(block func(Element) Element) ElementArray {
	result := ElementArray{}
	for _, o := range a {
		result = append(result, block(o))
	}
	return result
}

// Max returns the element the returned the highest value when passed to
// block
func (a ElementArray) Max(block func(Element) int) Element {
	if len(a) == 0 {
		return nil
	}
	var maxElement Element = a[0]
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

// Min returns the element the returned the lowest value when passed to
// block
func (a ElementArray) Min(block func(Element) int) Element {
	if len(a) == 0 {
		return nil
	}
	var minElement Element = a[0]
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

// Push appends an element to the array returning a new modified array
func (a ElementArray) Push(element Element) ElementArray {
	return append(a, element)
}

// Pop removes the last element from the array, returning new array and
// the element
func (a ElementArray) Pop() (ElementArray, Element) {
	return a[:len(a)-1], a[len(a)-1]
}

// Unshift adds an element to the array returning a new modified array
func (a ElementArray) Unshift(element Element) ElementArray {
	return append(ElementArray{element}, a...)
}

// Shift will remove the first element of the array returning the element
// and a modified array
func (a ElementArray) Shift() (Element, ElementArray) {
	if len(a) == 0 {
		return nil, a
	}
	return a[0], a[1:]
}

// Reverse will reverse the array in reverse returning the array reference
// again
func (a ElementArray) Reverse() ElementArray {
	for i := len(a)/2 - 1; i >= 0; i-- {
		opp := len(a) - 1 - i
		a[i], a[opp] = a[opp], a[i]
	}
	return a
}

// Shuffle will randomly shuffle an array elements order returning array
// reference again
func (a ElementArray) Shuffle() ElementArray {
	for i := len(a) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		a[i], a[j] = a[j], a[i]
	}
	return a
}
