package types

import "testing"

func TestSliceDrop(t *testing.T) {
	t.Run("drops specified number of elements from start", func(t *testing.T) {
		s := Slice[int]{1, 2, 3, 4, 5}
		result := s.Drop(2)
		expected := Slice[int]{3, 4, 5}
		AssertSlicesEquals(t, expected, result)
	})

	t.Run("returns empty slice when dropping all elements", func(t *testing.T) {
		s := Slice[int]{1, 2, 3}
		result := s.Drop(3)
		expected := Slice[int]{}
		AssertSlicesEquals(t, expected, result)
	})

	t.Run("returns empty slice when dropping more than length", func(t *testing.T) {
		s := Slice[int]{1, 2, 3}
		result := s.Drop(10)
		expected := Slice[int]{}
		AssertSlicesEquals(t, expected, result)
	})

	t.Run("returns original slice when dropping zero elements", func(t *testing.T) {
		s := Slice[int]{1, 2, 3, 4, 5}
		result := s.Drop(0)
		AssertSlicesEquals(t, s, result)
	})

	t.Run("returns original slice when dropping negative count", func(t *testing.T) {
		s := Slice[int]{1, 2, 3, 4, 5}
		result := s.Drop(-5)
		AssertSlicesEquals(t, s, result)
	})

	t.Run("handles empty slice", func(t *testing.T) {
		s := Slice[int]{}
		result := s.Drop(5)
		expected := Slice[int]{}
		AssertSlicesEquals(t, expected, result)
	})

	t.Run("drops one element", func(t *testing.T) {
		s := Slice[int]{10, 20, 30}
		result := s.Drop(1)
		expected := Slice[int]{20, 30}
		AssertSlicesEquals(t, expected, result)
	})

	t.Run("works with string slice", func(t *testing.T) {
		s := Slice[string]{"a", "b", "c", "d"}
		result := s.Drop(2)
		expected := Slice[string]{"c", "d"}
		AssertSlicesEquals(t, expected, result)
	})

	t.Run("works with single element slice", func(t *testing.T) {
		s := Slice[int]{42}
		result := s.Drop(1)
		expected := Slice[int]{}
		AssertSlicesEquals(t, expected, result)
	})

	t.Run("preserves original slice", func(t *testing.T) {
		s := Slice[int]{1, 2, 3, 4, 5}
		original := make(Slice[int], len(s))
		copy(original, s)
		_ = s.Drop(2)
		AssertSlicesEquals(t, original, s)
	})
}
