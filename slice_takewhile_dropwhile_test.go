package types

import (
	"testing"
)

func TestSliceTakeWhile(t *testing.T) {
	t.Run("takes elements while predicate is true", func(t *testing.T) {
		s := Slice[int]{1, 2, 3, 4, 5, 6}
		result := s.TakeWhile(func(x int) bool { return x < 4 })
		expected := Slice[int]{1, 2, 3}
		if !result.IsEq(expected) {
			t.Errorf("TakeWhile: expected %v, got %v", expected, result)
		}
	})

	t.Run("returns empty slice when first element fails predicate", func(t *testing.T) {
		s := Slice[int]{5, 6, 7, 8}
		result := s.TakeWhile(func(x int) bool { return x < 5 })
		expected := Slice[int]{}
		if !result.IsEq(expected) {
			t.Errorf("TakeWhile: expected %v, got %v", expected, result)
		}
	})

	t.Run("returns all elements when predicate always true", func(t *testing.T) {
		s := Slice[int]{1, 2, 3, 4}
		result := s.TakeWhile(func(x int) bool { return x < 10 })
		if !result.IsEq(s) {
			t.Errorf("TakeWhile: expected %v, got %v", s, result)
		}
	})

	t.Run("returns empty slice for empty input", func(t *testing.T) {
		s := Slice[int]{}
		result := s.TakeWhile(func(x int) bool { return x < 5 })
		expected := Slice[int]{}
		if !result.IsEq(expected) {
			t.Errorf("TakeWhile: expected %v, got %v", expected, result)
		}
	})

	t.Run("works with strings", func(t *testing.T) {
		s := Slice[string]{"apple", "banana", "cherry", "date"}
		result := s.TakeWhile(func(x string) bool { return len(x) < 6 })
		expected := Slice[string]{"apple"}
		if !result.IsEq(expected) {
			t.Errorf("TakeWhile: expected %v, got %v", expected, result)
		}
	})

	t.Run("stops at first false predicate", func(t *testing.T) {
		s := Slice[int]{2, 4, 6, 3, 8, 10}
		result := s.TakeWhile(func(x int) bool { return x%2 == 0 })
		expected := Slice[int]{2, 4, 6}
		if !result.IsEq(expected) {
			t.Errorf("TakeWhile: expected %v, got %v", expected, result)
		}
	})
}

func TestSliceDropWhile(t *testing.T) {
	t.Run("drops elements while predicate is true", func(t *testing.T) {
		s := Slice[int]{1, 2, 3, 4, 5, 6}
		result := s.DropWhile(func(x int) bool { return x < 4 })
		expected := Slice[int]{4, 5, 6}
		if !result.IsEq(expected) {
			t.Errorf("DropWhile: expected %v, got %v", expected, result)
		}
	})

	t.Run("returns original slice when first element fails predicate", func(t *testing.T) {
		s := Slice[int]{5, 6, 7, 8}
		result := s.DropWhile(func(x int) bool { return x < 5 })
		if !result.IsEq(s) {
			t.Errorf("DropWhile: expected %v, got %v", s, result)
		}
	})

	t.Run("returns empty slice when predicate always true", func(t *testing.T) {
		s := Slice[int]{1, 2, 3, 4}
		result := s.DropWhile(func(x int) bool { return x < 10 })
		expected := Slice[int]{}
		if !result.IsEq(expected) {
			t.Errorf("DropWhile: expected %v, got %v", expected, result)
		}
	})

	t.Run("returns empty slice for empty input", func(t *testing.T) {
		s := Slice[int]{}
		result := s.DropWhile(func(x int) bool { return x < 5 })
		expected := Slice[int]{}
		if !result.IsEq(expected) {
			t.Errorf("DropWhile: expected %v, got %v", expected, result)
		}
	})

	t.Run("works with strings", func(t *testing.T) {
		s := Slice[string]{"apple", "banana", "cherry", "date"}
		result := s.DropWhile(func(x string) bool { return len(x) < 6 })
		expected := Slice[string]{"banana", "cherry", "date"}
		if !result.IsEq(expected) {
			t.Errorf("DropWhile: expected %v, got %v", expected, result)
		}
	})

	t.Run("stops dropping at first false predicate", func(t *testing.T) {
		s := Slice[int]{2, 4, 6, 3, 8, 10}
		result := s.DropWhile(func(x int) bool { return x%2 == 0 })
		expected := Slice[int]{3, 8, 10}
		if !result.IsEq(expected) {
			t.Errorf("DropWhile: expected %v, got %v", expected, result)
		}
	})
}

func TestTakeWhileDropWhileComplement(t *testing.T) {
	t.Run("TakeWhile and DropWhile are complementary", func(t *testing.T) {
		s := Slice[int]{1, 2, 3, 4, 5, 6, 7, 8}
		predicate := func(x int) bool { return x < 5 }

		taken := s.TakeWhile(predicate)
		dropped := s.DropWhile(predicate)

		// Concatenate taken and dropped should equal original
		combined := append(taken, dropped...)

		if !combined.IsEq(s) {
			t.Errorf("TakeWhile and DropWhile should be complementary: expected %v, got %v", s, combined)
		}
	})

	t.Run("lengths should sum to original length", func(t *testing.T) {
		s := Slice[int]{10, 20, 30, 40, 50}
		predicate := func(x int) bool { return x <= 30 }

		taken := s.TakeWhile(predicate)
		dropped := s.DropWhile(predicate)

		if taken.Len()+dropped.Len() != s.Len() {
			t.Errorf("Lengths don't match: taken %d + dropped %d != original %d",
				taken.Len(), dropped.Len(), s.Len())
		}
	})
}
