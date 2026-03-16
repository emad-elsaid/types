// Package types provides a generic data types similar to that of Ruby
package types

import (
	"math/rand"
	"slices"
)

// Slice is an alias for a slice of variables that vary in types
// Slice can hold any comparable data type
type Slice[T comparable] []T

// At returns element by index, a negative index counts from the end of
// the Slice
// if index is out of range it returns nil
func (a Slice[T]) At(index int) *T {
	len := len(a)

	if index < 0 {
		if -index <= len {
			return &a[len+index]
		}
		return nil
	}

	if index < len {
		return &a[index]
	}

	return nil
}

// CountElement returns number of elements equal to "element" in Slice
func (a Slice[T]) CountElement(element T) (count int) {
	for _, o := range a {
		if o == element {
			count++
		}
	}
	return count
}

// CountBy returns number of elements which "block" returns true for
func (a Slice[T]) CountBy(block func(T) bool) (count int) {
	for _, o := range a {
		if block(o) {
			count++
		}
	}
	return count
}

// Compact returns a new slice with all zero values removed.
// For comparable types, zero values are: 0 for numbers, "" for strings, false for bools, nil for pointers.
// This is inspired by Ruby's Array#compact method.
func (a Slice[T]) Compact() Slice[T] {
	var zero T
	result := make(Slice[T], 0, len(a))
	for _, element := range a {
		if element != zero {
			result = append(result, element)
		}
	}
	return result
}

// Cycle will cycle through Slice elements "count" times passing each
// element to "block" function
func (a Slice[T]) Cycle(count int, block func(T)) {
	for range count {
		for _, v := range a {
			block(v)
		}
	}
}

// Any returns true if "block" returned true for any of the Slice elements
// and false otherwise
func (a Slice[T]) Any(block func(T) bool) bool {
	return slices.ContainsFunc(a, block)
}

// All returns true if "block" returned true for all elements in Slice and
// false otherwise
func (a Slice[T]) All(block func(T) bool) bool {
	for _, o := range a {
		if !block(o) {
			return false
		}
	}

	return true
}

// None returns true if no elements in the slice satisfy the predicate function.
// Returns true for empty slices.
func (a Slice[T]) None(predicate func(T) bool) bool {
	return !a.Any(predicate)
}

// Find returns the first element that satisfies the predicate function and true.
// If no element satisfies the predicate, it returns the zero value of T and false.
func (a Slice[T]) Find(predicate func(T) bool) (T, bool) {
	for _, item := range a {
		if predicate(item) {
			return item, true
		}
	}
	var zero T
	return zero, false
}

// Delete will remove all elements that are equal to the passed element
func (a Slice[T]) Delete(element T) Slice[T] {
	result := Slice[T]{}
	for _, o := range a {
		if o != element {
			result = append(result, o)
		}
	}
	return result
}

// DeleteAt will delete an element by index
// If index is out of bounds (negative or >= length), returns original slice unchanged
func (a Slice[T]) DeleteAt(index int) Slice[T] {
	if index < 0 || index >= len(a) {
		return a
	}
	result := make(Slice[T], 0, len(a)-1)
	result = append(result, a[:index]...)
	result = append(result, a[index+1:]...)
	return result
}

// DeleteIf will delete all elements which "block" returns true for
func (a Slice[T]) DeleteIf(block func(T) bool) Slice[T] {
	result := Slice[T]{}
	for _, o := range a {
		if !block(o) {
			result = append(result, o)
		}
	}
	return result
}

// Drop will return an array without the first "count" elements from the
// beginning
func (a Slice[T]) Drop(count int) Slice[T] {
	if count <= 0 {
		return a
	}
	if count >= len(a) {
		return Slice[T]{}
	}
	return a[count:]
}

// TakeWhile returns a new slice containing elements from the start of the slice
// while the predicate returns true. Stops at the first element where predicate returns false.
func (a Slice[T]) TakeWhile(predicate func(T) bool) Slice[T] {
	for i, elem := range a {
		if !predicate(elem) {
			return a[:i]
		}
	}
	return a
}

// DropWhile returns a new slice with elements dropped from the start while the
// predicate returns true. Returns the remaining elements starting from the first
// element where predicate returns false.
func (a Slice[T]) DropWhile(predicate func(T) bool) Slice[T] {
	for i, elem := range a {
		if !predicate(elem) {
			return a[i:]
		}
	}
	return Slice[T]{}
}

// Each will execute "block" for each element in array
func (a Slice[T]) Each(block func(T)) {
	for _, o := range a {
		block(o)
	}
}

// EachIndex will execute "block" for each element index in array
func (a Slice[T]) EachIndex(block func(int)) {
	for i := range a {
		block(i)
	}
}

// IsEmpty will return true of array empty, false otherwise
func (a Slice[T]) IsEmpty() bool {
	return len(a) == 0
}

// IsEq returns true if array the "other" array
func (a Slice[T]) IsEq(other Slice[T]) bool {
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
func (a Slice[T]) Len() int {
	return len(a)
}

// Fetch will return the element in "index", if it doesn't exist it
// returns the passed "defaultValue"
func (a Slice[T]) Fetch(index int, defaultValue T) T {
	val := a.At(index)
	if val != nil {
		return *val
	}

	return defaultValue
}

// Fill will replace elements inplace starting from "start" counting
// "length" elements with the passed "element" parameter, will return same array
// object
func (a Slice[T]) Fill(element T, start int, length int) Slice[T] {
	if start < 0 || start >= len(a) || length <= 0 {
		return a
	}
	// Adjust length if it would exceed slice bounds
	if start+length > len(a) {
		length = len(a) - start
	}
	for length--; length >= 0; length-- {
		a[start+length] = element
	}
	return a
}

// FillWith will replace elements from start counting "length" items,
// passing every index to block and replacing the element inplace with the
// return value
func (a Slice[T]) FillWith(start int, length int, block func(int) T) Slice[T] {
	if start < 0 || start >= len(a) || length <= 0 {
		return a
	}
	// Adjust length if it would exceed slice bounds
	if start+length > len(a) {
		length = len(a) - start
	}
	for length--; length >= 0; length-- {
		a[start+length] = block(start + length)
	}
	return a
}

// Index returns the index of the first element in array that is equal to
// "element", returns -1 if the elements if not found
func (a Slice[T]) Index(element T) int {
	for i, o := range a {
		if o == element {
			return i
		}
	}
	return -1
}

// IndexBy returns first element that block returns true for, -1 otherwise
func (a Slice[T]) IndexBy(block func(T) bool) int {
	for i, o := range a {
		if block(o) {
			return i
		}
	}

	return -1
}

// First returns first element of array
func (a Slice[T]) First() *T {
	return a.At(0)
}

// Last returns last element of array
func (a Slice[T]) Last() *T {
	return a.At(len(a) - 1)
}

// Firsts will return an array holding the first "count" elements of the
// array
func (a Slice[T]) Firsts(count int) Slice[T] {
	if count <= 0 {
		return Slice[T]{}
	}
	if count >= len(a) {
		return a
	}
	return a[0:count]
}

// Lasts will return an array holding the lasts "count" elements of the
// array
func (a Slice[T]) Lasts(count int) Slice[T] {
	if count <= 0 {
		return Slice[T]{}
	}
	if count >= len(a) {
		return a
	}
	return a[len(a)-count:]
}

// Include will return true if element found in the array
func (a Slice[T]) Include(element T) bool {
	return a.Index(element) != -1
}

// Insert will insert a set of elements in the index and will return a new
// array
func (a Slice[T]) Insert(index int, elements ...T) Slice[T] {
	result := Slice[T]{}
	result = append(result, a[0:index]...)
	result = append(result, elements...)
	result = append(result, a[index:]...)
	return result
}

// KeepIf will return an array contains all elements where "block"
// returned true for them
func (a Slice[T]) KeepIf(block func(T) bool) Slice[T] {
	result := Slice[T]{}
	for _, o := range a {
		if block(o) {
			result = append(result, o)
		}
	}
	return result
}

// Select is an alias for KeepIf
func (a Slice[T]) Select(block func(T) bool) Slice[T] {
	return a.KeepIf(block)
}

// SelectUntil selects from the start of the slice until the block returns true,
// excluding the item that returned true.
func (a Slice[T]) SelectUntil(block func(T) bool) Slice[T] {
	index := a.IndexBy(block)
	if index == -1 {
		return a
	}

	return a[0:index]
}

// Reduce is an alias for KeepIf
func (a Slice[T]) Reduce(block func(T) bool) Slice[T] {
	return a.KeepIf(block)
}

// Map will return a new array replacing every element from current array
// with the return value of the block
func (a Slice[T]) Map(block func(T) T) Slice[T] {
	result := Slice[T]{}
	for _, o := range a {
		result = append(result, block(o))
	}
	return result
}

// Max returns the element the returned the highest value when passed to
// block
func (a Slice[T]) Max(block func(T) int) T {
	if len(a) == 0 {
		return *new(T)
	}
	var maxElement T = a[0]
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
func (a Slice[T]) Min(block func(T) int) T {
	if len(a) == 0 {
		return *new(T)
	}
	var minElement T = a[0]
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
func (a Slice[T]) Push(element T) Slice[T] {
	return append(a, element)
}

// Pop removes the last element from the array, returning new array and
// the element. Returns the original slice and zero value if slice is empty.
func (a Slice[T]) Pop() (Slice[T], T) {
	if len(a) == 0 {
		var zero T
		return a, zero
	}
	return a[:len(a)-1], a[len(a)-1]
}

// Unshift adds an element to the array returning a new modified array
func (a Slice[T]) Unshift(element T) Slice[T] {
	return append(Slice[T]{element}, a...)
}

// Shift will remove the first element of the array returning the element
// and a modified array
func (a Slice[T]) Shift() (T, Slice[T]) {
	if len(a) == 0 {
		return *new(T), a
	}
	return a[0], a[1:]
}

// Reverse returns a new slice with elements in reverse order.
// The original slice is not modified.
func (a Slice[T]) Reverse() Slice[T] {
	result := make(Slice[T], len(a))
	copy(result, a)
	for i := len(result)/2 - 1; i >= 0; i-- {
		opp := len(result) - 1 - i
		result[i], result[opp] = result[opp], result[i]
	}
	return result
}

// Shuffle returns a new slice with elements in random order.
// The original slice is not modified.
func (a Slice[T]) Shuffle() Slice[T] {
	result := make(Slice[T], len(a))
	copy(result, a)
	for i := len(result) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		result[i], result[j] = result[j], result[i]
	}
	return result
}

// Sample returns a random element from the slice.
// If the slice is empty, returns the zero value of T and false.
// The original slice is not modified.
func (a Slice[T]) Sample() (T, bool) {
	if len(a) == 0 {
		var zero T
		return zero, false
	}
	return a[rand.Intn(len(a))], true
}

// SampleN returns n random elements from the slice without replacement.
// If n is greater than the slice length, returns all elements in random order.
// If n is less than or equal to 0, returns an empty slice.
// The original slice is not modified.
func (a Slice[T]) SampleN(n int) Slice[T] {
	if n <= 0 {
		return Slice[T]{}
	}
	if n >= len(a) {
		return a.Shuffle()
	}

	// Create a shuffled copy and take first n elements
	shuffled := a.Shuffle()
	return shuffled[:n]
}

// Unique returns a unique list of elements in the slice. order or the result is
// same order of the items in the source slice
func (a Slice[T]) Unique() Slice[T] {
	keys := map[T]struct{}{}
	out := Slice[T]{}
	for _, v := range a {
		if _, ok := keys[v]; !ok {
			keys[v] = struct{}{}
			out = append(out, v)
		}
	}

	return out
}

// Partition divides the slice into two new slices based on the predicate function.
// Returns two slices: the first contains elements that satisfy the predicate,
// the second contains elements that do not satisfy the predicate.
func (s Slice[T]) Partition(predicate func(T) bool) (Slice[T], Slice[T]) {
	trueSet := Slice[T]{}
	falseSet := Slice[T]{}

	for _, item := range s {
		if predicate(item) {
			trueSet = append(trueSet, item)
		} else {
			falseSet = append(falseSet, item)
		}
	}

	return trueSet, falseSet
}

// Tally returns a map counting the occurrences of each element in the slice.
// This is inspired by Ruby's Array#tally method.
// Example: Slice[int]{1, 2, 2, 3, 3, 3}.Tally() returns map[int]int{1: 1, 2: 2, 3: 3}
func (a Slice[T]) Tally() map[T]int {
	result := make(map[T]int)
	for _, element := range a {
		result[element]++
	}
	return result
}

// SliceReduce applies a function against an accumulator and each element in the slice
func SliceReduce[T comparable, U any](s Slice[T], initial U, fn func(U, T) U) U {
	result := initial
	for _, item := range s {
		result = fn(result, item)
	}

	return result
}
