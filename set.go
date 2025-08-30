package types

import (
	"fmt"
	"slices"
	"strings"
)

// Set represents a generic set data structure that stores unique elements of type T.
// T must be comparable to be used as map keys.
type Set[T comparable] struct {
	order []T
	items map[T]struct{}
}

// NewSet creates and returns a new Set initialized with elements from the given slice.
// Duplicate elements in the slice will be automatically deduplicated.
func NewSet[T comparable](slice ...T) *Set[T] {
	s := Set[T]{
		order: make([]T, 0, len(slice)),
		items: make(map[T]struct{}, len(slice)),
	}

	for _, item := range slice {
		s.Add(item)
	}

	return &s
}

// Add inserts an element into the set.
// Returns true if the element was added (wasn't already present), false otherwise.
func (s *Set[T]) Add(item T) bool {
	if s.Contains(item) {
		return false
	}

	s.items[item] = struct{}{}
	s.order = append(s.order, item)

	return true
}

// Remove deletes an element from the set.
// Returns true if the element was removed (was present), false otherwise.
func (s *Set[T]) Remove(item T) bool {
	if !s.Contains(item) {
		return false
	}

	delete(s.items, item)
	index := slices.Index(s.order, item)
	s.order = slices.Delete(s.order, index, index+1)

	return true
}

// Contains checks if an element exists in the set.
// Returns true if the element is present, false otherwise.
func (s *Set[T]) Contains(item T) bool {
	_, exists := s.items[item]
	return exists
}

// Size returns the number of elements in the set.
func (s *Set[T]) Size() int {
	return len(s.order)
}

// IsEmpty returns true if the set contains no elements, false otherwise.
func (s *Set[T]) IsEmpty() bool {
	return len(s.order) == 0
}

// Clear removes all elements from the set.
func (s *Set[T]) Clear() {
	s.items = make(map[T]struct{})
	s.order = []T{}
}

// ToSlice returns a slice containing all elements in the set in the order they were added.
// the slice is not a copy, modifying it will affect the set.
func (s *Set[T]) ToSlice() []T {
	return s.order
}

// Clone creates and returns a shallow copy of the set.
func (s *Set[T]) Clone() *Set[T] {
	return NewSet(s.order...)
}

// Union returns a new set containing all elements that are in either this set or the other set.
func (s *Set[T]) Union(other *Set[T]) *Set[T] {
	result := s.Clone()
	for _, item := range other.order {
		result.Add(item)
	}

	return result
}

// Intersection returns a new set containing only elements that are in both this
// set and the other set in the order they were added to the first set.
func (s *Set[T]) Intersection(other *Set[T]) *Set[T] {
	result := NewSet[T]()
	for _, item := range s.order {
		if other.Contains(item) {
			result.Add(item)
		}
	}
	return result
}

// Difference returns a new set containing elements that are in this set but not
// in the other set in the order they were added to the first set.
func (s *Set[T]) Difference(other *Set[T]) *Set[T] {
	result := NewSet[T]()
	for _, item := range s.order {
		if !other.Contains(item) {
			result.Add(item)
		}
	}
	return result
}

// SymmetricDifference returns a new set containing elements that are in either this set or the other set, but not in both.
func (s *Set[T]) SymmetricDifference(other *Set[T]) *Set[T] {
	return s.Union(other).Difference(s.Intersection(other))
}

// IsSubset returns true if this set is a subset of the other set (all elements of this set are in the other set).
func (s *Set[T]) IsSubset(other *Set[T]) bool {
	for _, item := range s.order {
		if !other.Contains(item) {
			return false
		}
	}
	return true
}

// IsSuperset returns true if this set is a superset of the other set (all elements of the other set are in this set).
func (s *Set[T]) IsSuperset(other *Set[T]) bool {
	return other.IsSubset(s)
}

// IsDisjoint returns true if this set has no elements in common with the other set.
func (s *Set[T]) IsDisjoint(other *Set[T]) bool {
	return s.Intersection(other).IsEmpty()
}

// Equal returns true if this set contains exactly the same elements as the other set.
func (s *Set[T]) Equal(other *Set[T]) bool {
	return s.Size() == other.Size() && s.IsSubset(other)
}

// Each iterates over all elements in the set and calls the provided function for each element.
// The order of iteration is not guaranteed.
func (s *Set[T]) Each(fn func(T)) {
	for _, item := range s.order {
		fn(item)
	}
}

// SetMap applies a transformation function to each element and returns a new set with the results.
// The transformation function must return a comparable type.
func SetMap[T, U comparable](s *Set[T], fn func(T) U) *Set[U] {
	result := NewSet[U]()
	result.order = make([]U, 0, len(s.order))
	result.items = make(map[U]struct{}, len(s.order))

	for _, item := range s.order {
		result.Add(fn(item))
	}

	return result
}

// Filter returns a new set containing only elements that satisfy the predicate function.
func (s *Set[T]) Filter(predicate func(T) bool) *Set[T] {
	result := NewSet[T]()
	result.order = make([]T, 0, len(s.order))
	result.items = make(map[T]struct{}, len(s.order))

	for _, item := range s.order {
		if predicate(item) {
			result.Add(item)
		}
	}

	return result
}

// Reject returns a new set containing only elements that do not satisfy the predicate function.
// This is the opposite of Filter.
func (s *Set[T]) Reject(predicate func(T) bool) *Set[T] {
	return s.Filter(func(item T) bool {
		return !predicate(item)
	})
}

// Find returns the first element that satisfies the predicate function and true.
// If no element satisfies the predicate, it returns the zero value of T and false.
func (s *Set[T]) Find(predicate func(T) bool) (T, bool) {
	for _, item := range s.order {
		if predicate(item) {
			return item, true
		}
	}
	var zero T
	return zero, false
}

// All returns true if all elements in the set satisfy the predicate function.
// Returns true for empty sets.
func (s *Set[T]) All(predicate func(T) bool) bool {
	for _, item := range s.order {
		if !predicate(item) {
			return false
		}
	}
	return true
}

// Any returns true if at least one element in the set satisfies the predicate function.
// Returns false for empty sets.
func (s *Set[T]) Any(predicate func(T) bool) bool {
	return slices.ContainsFunc(s.order, predicate)
}

// None returns true if no elements in the set satisfy the predicate function.
// Returns true for empty sets.
func (s *Set[T]) None(predicate func(T) bool) bool {
	return !s.Any(predicate)
}

// Count returns the number of elements that satisfy the predicate function.
func (s *Set[T]) Count(predicate func(T) bool) int {
	count := 0
	for _, item := range s.order {
		if predicate(item) {
			count++
		}
	}
	return count
}

// SetReduce applies a reduction function to all elements in the set, starting with an initial value.
// The reduction function takes the accumulated value and the current element, returning the new accumulated value.
func SetReduce[T comparable, U any](s *Set[T], initial U, fn func(U, T) U) U {
	result := initial
	for _, item := range s.order {
		result = fn(result, item)
	}
	return result
}

// Partition divides the set into two new sets based on the predicate function.
// Returns two sets: the first contains elements that satisfy the predicate,
// the second contains elements that do not satisfy the predicate.
func (s *Set[T]) Partition(predicate func(T) bool) (*Set[T], *Set[T]) {
	trueSet := NewSet[T]()
	falseSet := NewSet[T]()

	for _, item := range s.order {
		if predicate(item) {
			trueSet.Add(item)
		} else {
			falseSet.Add(item)
		}
	}

	return trueSet, falseSet
}

// Take returns a new set containing up to n elements from this set in the order they were added.
func (s *Set[T]) Take(n int) *Set[T] {
	if n <= 0 {
		return NewSet[T]()
	}

	result := NewSet[T]()
	count := 0

	for _, item := range s.order {
		if count >= n {
			break
		}
		result.Add(item)
		count++
	}

	return result
}

// Drop returns a new set with the first n elements removed in the order they were added.
func (s *Set[T]) Drop(n int) *Set[T] {
	if n <= 0 {
		return s.Clone()
	}

	result := NewSet[T]()
	count := 0

	for _, item := range s.order {
		if count >= n {
			result.Add(item)
		}
		count++
	}

	return result
}

// String returns a string representation of the set.
func (s *Set[T]) String() string {
	if s.IsEmpty() {
		return "Set{}"
	}

	items := s.ToSlice()
	strs := make([]string, len(items))
	for i, item := range items {
		strs[i] = fmt.Sprintf("%v", item)
	}

	return fmt.Sprintf("Set{%s}", strings.Join(strs, ", "))
}
