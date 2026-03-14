package types

import (
	"testing"
)

func TestSliceCompact(t *testing.T) {
	t.Run("removes zero values from int slice", func(t *testing.T) {
		s := Slice[int]{1, 0, 2, 0, 3, 0}
		result := s.Compact()
		expected := Slice[int]{1, 2, 3}
		if !result.IsEq(expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("removes empty strings from string slice", func(t *testing.T) {
		s := Slice[string]{"hello", "", "world", "", "!"}
		result := s.Compact()
		expected := Slice[string]{"hello", "world", "!"}
		if !result.IsEq(expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("removes false from bool slice", func(t *testing.T) {
		s := Slice[bool]{true, false, true, false, true}
		result := s.Compact()
		expected := Slice[bool]{true, true, true}
		if !result.IsEq(expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("returns empty slice when all values are zero", func(t *testing.T) {
		s := Slice[int]{0, 0, 0, 0}
		result := s.Compact()
		if len(result) != 0 {
			t.Errorf("Expected empty slice, got %v", result)
		}
	})

	t.Run("returns same slice when no zero values", func(t *testing.T) {
		s := Slice[int]{1, 2, 3, 4, 5}
		result := s.Compact()
		if !result.IsEq(s) {
			t.Errorf("Expected %v, got %v", s, result)
		}
	})

	t.Run("returns empty slice when input is empty", func(t *testing.T) {
		s := Slice[int]{}
		result := s.Compact()
		if len(result) != 0 {
			t.Errorf("Expected empty slice, got %v", result)
		}
	})

	t.Run("works with float64", func(t *testing.T) {
		s := Slice[float64]{1.5, 0.0, 2.5, 0.0, 3.5}
		result := s.Compact()
		expected := Slice[float64]{1.5, 2.5, 3.5}
		if !result.IsEq(expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("removes zero runes from rune slice", func(t *testing.T) {
		s := Slice[rune]{'a', 0, 'b', 0, 'c'}
		result := s.Compact()
		expected := Slice[rune]{'a', 'b', 'c'}
		if !result.IsEq(expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("preserves order of non-zero elements", func(t *testing.T) {
		s := Slice[int]{5, 0, 4, 0, 3, 0, 2, 0, 1}
		result := s.Compact()
		expected := Slice[int]{5, 4, 3, 2, 1}
		if !result.IsEq(expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("single zero value returns empty slice", func(t *testing.T) {
		s := Slice[int]{0}
		result := s.Compact()
		if len(result) != 0 {
			t.Errorf("Expected empty slice, got %v", result)
		}
	})

	t.Run("single non-zero value returns same value", func(t *testing.T) {
		s := Slice[int]{42}
		result := s.Compact()
		expected := Slice[int]{42}
		if !result.IsEq(expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})
}
