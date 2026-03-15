package types

import (
	"testing"
)

func TestSliceSample(t *testing.T) {
	t.Run("returns random element from non-empty slice", func(t *testing.T) {
		s := Slice[int]{1, 2, 3, 4, 5}
		elem, ok := s.Sample()
		if !ok {
			t.Fatal("Sample should return true for non-empty slice")
		}
		if !s.Include(elem) {
			t.Errorf("Sample returned element %v not in slice", elem)
		}
	})

	t.Run("returns false for empty slice", func(t *testing.T) {
		s := Slice[int]{}
		elem, ok := s.Sample()
		if ok {
			t.Error("Sample should return false for empty slice")
		}
		if elem != 0 {
			t.Errorf("Sample should return zero value for empty slice, got %v", elem)
		}
	})

	t.Run("works with string slice", func(t *testing.T) {
		s := Slice[string]{"apple", "banana", "cherry"}
		elem, ok := s.Sample()
		if !ok {
			t.Fatal("Sample should return true for non-empty slice")
		}
		if !s.Include(elem) {
			t.Errorf("Sample returned element %v not in slice", elem)
		}
	})

	t.Run("returns different elements on multiple calls (probabilistic)", func(t *testing.T) {
		s := Slice[int]{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		seen := make(map[int]bool)
		for i := 0; i < 50; i++ {
			elem, _ := s.Sample()
			seen[elem] = true
		}
		// With 50 samples from 10 elements, we should see at least 5 different values
		if len(seen) < 5 {
			t.Errorf("Sample should return varied elements, only saw %d unique values", len(seen))
		}
	})

	t.Run("single element slice always returns that element", func(t *testing.T) {
		s := Slice[int]{42}
		for i := 0; i < 10; i++ {
			elem, ok := s.Sample()
			if !ok {
				t.Fatal("Sample should return true for non-empty slice")
			}
			if elem != 42 {
				t.Errorf("Sample should return 42, got %v", elem)
			}
		}
	})
}

func TestSliceSampleN(t *testing.T) {
	t.Run("returns n random elements", func(t *testing.T) {
		s := Slice[int]{1, 2, 3, 4, 5}
		result := s.SampleN(3)
		if len(result) != 3 {
			t.Errorf("SampleN(3) should return 3 elements, got %d", len(result))
		}
		// All returned elements should be in original slice
		for _, elem := range result {
			if !s.Include(elem) {
				t.Errorf("SampleN returned element %v not in original slice", elem)
			}
		}
		// No duplicates (sampling without replacement)
		if len(result.Unique()) != len(result) {
			t.Error("SampleN should return unique elements (no duplicates)")
		}
	})

	t.Run("returns empty slice for n <= 0", func(t *testing.T) {
		s := Slice[int]{1, 2, 3}
		if len(s.SampleN(0)) != 0 {
			t.Error("SampleN(0) should return empty slice")
		}
		if len(s.SampleN(-5)) != 0 {
			t.Error("SampleN(-5) should return empty slice")
		}
	})

	t.Run("returns all elements shuffled when n >= length", func(t *testing.T) {
		s := Slice[int]{1, 2, 3}
		result := s.SampleN(5)
		if len(result) != 3 {
			t.Errorf("SampleN(5) should return all 3 elements, got %d", len(result))
		}
		// Should contain all original elements
		for _, elem := range s {
			if !result.Include(elem) {
				t.Errorf("SampleN should include element %v", elem)
			}
		}
	})

	t.Run("works with empty slice", func(t *testing.T) {
		s := Slice[int]{}
		result := s.SampleN(3)
		if len(result) != 0 {
			t.Error("SampleN on empty slice should return empty slice")
		}
	})

	t.Run("returns all elements when n equals length", func(t *testing.T) {
		s := Slice[int]{10, 20, 30}
		result := s.SampleN(3)
		if len(result) != 3 {
			t.Errorf("SampleN(3) should return 3 elements, got %d", len(result))
		}
		// Should contain all original elements
		for _, elem := range s {
			if !result.Include(elem) {
				t.Errorf("SampleN should include element %v", elem)
			}
		}
	})

	t.Run("sampling is random (probabilistic test)", func(t *testing.T) {
		s := Slice[int]{1, 2, 3, 4, 5}
		results := make(map[int]int) // Count how often each element appears first

		for i := 0; i < 100; i++ {
			sample := s.SampleN(2)
			if len(sample) > 0 {
				results[sample[0]]++
			}
		}

		// With 100 samples, each of the 5 elements should appear at least a few times
		// (probabilistic, but with 100 samples this should almost always pass)
		if len(results) < 3 {
			t.Errorf("SampleN should show randomness, only saw %d different first elements", len(results))
		}
	})

	t.Run("works with string slice", func(t *testing.T) {
		s := Slice[string]{"red", "green", "blue", "yellow"}
		result := s.SampleN(2)
		if len(result) != 2 {
			t.Errorf("SampleN(2) should return 2 elements, got %d", len(result))
		}
		for _, elem := range result {
			if !s.Include(elem) {
				t.Errorf("SampleN returned element %v not in original slice", elem)
			}
		}
	})

	t.Run("does not modify original slice", func(t *testing.T) {
		original := Slice[int]{1, 2, 3, 4, 5}
		expected := Slice[int]{1, 2, 3, 4, 5}
		_ = original.SampleN(3)
		if !original.IsEq(expected) {
			t.Error("SampleN should not modify original slice")
		}
	})
}
