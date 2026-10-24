package types

// SetMap applies a transformation function to each element and returns a new set with the results.
// The transformation function must return a comparable type.
//
// Deprecated: Use (*Set[T]).Map instead.
func SetMap[T, U comparable](s *Set[T], fn func(T) U) *Set[U] {
	return s.Map(fn)
}

// SetReduce applies a reduction function to all elements in the set, starting with an initial value.
// The reduction function takes the accumulated value and the current element, returning the new accumulated value.
//
// Deprecated: Use (*Set[T]).Reduce instead.
func SetReduce[T comparable, U any](s *Set[T], initial U, fn func(U, T) U) U {
	return s.Reduce(initial, fn)
}

// Zip combines two slices element-wise, creating pairs.
// Returns a slice of pairs (2-element arrays) where each pair contains
// one element from the first slice and one from the second slice at the same index.
// The resulting slice will have the length of the shorter input slice.
// Inspired by Ruby's Array#zip method.
// Example: Zip(Slice[int]{1, 2, 3}, Slice[string]{"a", "b"}) returns [[1,"a"], [2,"b"]]
//
// Deprecated: Use Slice[T].Zip instead.
func Zip[T, U comparable](a Slice[T], b Slice[U]) [][2]any {
	return a.Zip(b)
}

// SliceReduce applies a function against an accumulator and each element in the slice
//
// Deprecated: Use Slice[T].Fold instead.
func SliceReduce[T comparable, U any](s Slice[T], initial U, fn func(U, T) U) U {
	return s.Fold(initial, fn)
}
