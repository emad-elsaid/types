package types

import "testing"

// TestSliceFirst tests the First method
func TestSliceFirst(t *testing.T) {
	t.Run("returns first element", func(t *testing.T) {
		s := Slice[int]{1, 2, 3}
		result := s.First()
		if result == nil || *result != 1 {
			t.Errorf("Expected 1, got %v", result)
		}
	})

	t.Run("returns nil for empty slice", func(t *testing.T) {
		s := Slice[int]{}
		result := s.First()
		if result != nil {
			t.Errorf("Expected nil for empty slice, got %v", result)
		}
	})

	t.Run("works with strings", func(t *testing.T) {
		s := Slice[string]{"a", "b", "c"}
		result := s.First()
		if result == nil || *result != "a" {
			t.Errorf("Expected 'a', got %v", result)
		}
	})
}

// TestSliceLast tests the Last method
func TestSliceLast(t *testing.T) {
	t.Run("returns last element", func(t *testing.T) {
		s := Slice[int]{1, 2, 3}
		result := s.Last()
		if result == nil || *result != 3 {
			t.Errorf("Expected 3, got %v", result)
		}
	})

	t.Run("returns nil for empty slice", func(t *testing.T) {
		s := Slice[int]{}
		result := s.Last()
		if result != nil {
			t.Errorf("Expected nil for empty slice, got %v", result)
		}
	})

	t.Run("works with single element", func(t *testing.T) {
		s := Slice[int]{42}
		result := s.Last()
		if result == nil || *result != 42 {
			t.Errorf("Expected 42, got %v", result)
		}
	})
}

// TestSliceFirsts tests the Firsts method
func TestSliceFirsts(t *testing.T) {
	t.Run("returns first n elements", func(t *testing.T) {
		s := Slice[int]{1, 2, 3, 4, 5}
		result := s.Firsts(3)
		expected := Slice[int]{1, 2, 3}
		if !result.IsEq(expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("returns empty for zero count", func(t *testing.T) {
		s := Slice[int]{1, 2, 3}
		result := s.Firsts(0)
		if len(result) != 0 {
			t.Errorf("Expected empty slice, got %v", result)
		}
	})

	t.Run("returns empty for negative count", func(t *testing.T) {
		s := Slice[int]{1, 2, 3}
		result := s.Firsts(-1)
		if len(result) != 0 {
			t.Errorf("Expected empty slice, got %v", result)
		}
	})

	t.Run("returns all elements when count exceeds length", func(t *testing.T) {
		s := Slice[int]{1, 2, 3}
		result := s.Firsts(10)
		if !result.IsEq(s) {
			t.Errorf("Expected %v, got %v", s, result)
		}
	})
}

// TestSliceLasts tests the Lasts method
func TestSliceLasts(t *testing.T) {
	t.Run("returns last n elements", func(t *testing.T) {
		s := Slice[int]{1, 2, 3, 4, 5}
		result := s.Lasts(3)
		expected := Slice[int]{3, 4, 5}
		if !result.IsEq(expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("returns empty for zero count", func(t *testing.T) {
		s := Slice[int]{1, 2, 3}
		result := s.Lasts(0)
		if len(result) != 0 {
			t.Errorf("Expected empty slice, got %v", result)
		}
	})

	t.Run("returns empty for negative count", func(t *testing.T) {
		s := Slice[int]{1, 2, 3}
		result := s.Lasts(-1)
		if len(result) != 0 {
			t.Errorf("Expected empty slice, got %v", result)
		}
	})

	t.Run("returns all elements when count exceeds length", func(t *testing.T) {
		s := Slice[int]{1, 2, 3}
		result := s.Lasts(10)
		if !result.IsEq(s) {
			t.Errorf("Expected %v, got %v", s, result)
		}
	})
}

// TestSliceIsEmpty tests the IsEmpty method
func TestSliceIsEmpty(t *testing.T) {
	t.Run("returns true for empty slice", func(t *testing.T) {
		s := Slice[int]{}
		if !s.IsEmpty() {
			t.Error("Expected true for empty slice")
		}
	})

	t.Run("returns false for non-empty slice", func(t *testing.T) {
		s := Slice[int]{1}
		if s.IsEmpty() {
			t.Error("Expected false for non-empty slice")
		}
	})

	t.Run("returns false for slice with multiple elements", func(t *testing.T) {
		s := Slice[string]{"a", "b", "c"}
		if s.IsEmpty() {
			t.Error("Expected false for non-empty slice")
		}
	})
}

// TestSlicePush tests the Push method
func TestSlicePush(t *testing.T) {
	t.Run("appends element to slice", func(t *testing.T) {
		s := Slice[int]{1, 2}
		result := s.Push(3)
		expected := Slice[int]{1, 2, 3}
		if !result.IsEq(expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("appends to empty slice", func(t *testing.T) {
		s := Slice[int]{}
		result := s.Push(1)
		expected := Slice[int]{1}
		if !result.IsEq(expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("works with strings", func(t *testing.T) {
		s := Slice[string]{"a", "b"}
		result := s.Push("c")
		expected := Slice[string]{"a", "b", "c"}
		if !result.IsEq(expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})
}

// TestSlicePop tests the Pop method
func TestSlicePop(t *testing.T) {
	t.Run("removes and returns last element", func(t *testing.T) {
		s := Slice[int]{1, 2, 3}
		newSlice, popped := s.Pop()
		expected := Slice[int]{1, 2}
		if !newSlice.IsEq(expected) {
			t.Errorf("Expected slice %v, got %v", expected, newSlice)
		}
		if popped != 3 {
			t.Errorf("Expected popped value 3, got %v", popped)
		}
	})

	t.Run("returns zero value for empty slice", func(t *testing.T) {
		s := Slice[int]{}
		newSlice, popped := s.Pop()
		if len(newSlice) != 0 {
			t.Errorf("Expected empty slice, got %v", newSlice)
		}
		if popped != 0 {
			t.Errorf("Expected zero value 0, got %v", popped)
		}
	})

	t.Run("works with single element", func(t *testing.T) {
		s := Slice[int]{42}
		newSlice, popped := s.Pop()
		if len(newSlice) != 0 {
			t.Errorf("Expected empty slice, got %v", newSlice)
		}
		if popped != 42 {
			t.Errorf("Expected 42, got %v", popped)
		}
	})
}

// TestSliceUnshift tests the Unshift method
func TestSliceUnshift(t *testing.T) {
	t.Run("prepends element to slice", func(t *testing.T) {
		s := Slice[int]{2, 3}
		result := s.Unshift(1)
		expected := Slice[int]{1, 2, 3}
		if !result.IsEq(expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("prepends to empty slice", func(t *testing.T) {
		s := Slice[int]{}
		result := s.Unshift(1)
		expected := Slice[int]{1}
		if !result.IsEq(expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("works with strings", func(t *testing.T) {
		s := Slice[string]{"b", "c"}
		result := s.Unshift("a")
		expected := Slice[string]{"a", "b", "c"}
		if !result.IsEq(expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})
}

// TestSliceShift tests the Shift method
func TestSliceShift(t *testing.T) {
	t.Run("removes and returns first element", func(t *testing.T) {
		s := Slice[int]{1, 2, 3}
		shifted, newSlice := s.Shift()
		expected := Slice[int]{2, 3}
		if !newSlice.IsEq(expected) {
			t.Errorf("Expected slice %v, got %v", expected, newSlice)
		}
		if shifted != 1 {
			t.Errorf("Expected shifted value 1, got %v", shifted)
		}
	})

	t.Run("returns zero value for empty slice", func(t *testing.T) {
		s := Slice[int]{}
		shifted, newSlice := s.Shift()
		if len(newSlice) != 0 {
			t.Errorf("Expected empty slice, got %v", newSlice)
		}
		if shifted != 0 {
			t.Errorf("Expected zero value 0, got %v", shifted)
		}
	})

	t.Run("works with single element", func(t *testing.T) {
		s := Slice[int]{42}
		shifted, newSlice := s.Shift()
		if len(newSlice) != 0 {
			t.Errorf("Expected empty slice, got %v", newSlice)
		}
		if shifted != 42 {
			t.Errorf("Expected 42, got %v", shifted)
		}
	})

	t.Run("works with strings", func(t *testing.T) {
		s := Slice[string]{"a", "b", "c"}
		shifted, newSlice := s.Shift()
		expected := Slice[string]{"b", "c"}
		if !newSlice.IsEq(expected) {
			t.Errorf("Expected slice %v, got %v", expected, newSlice)
		}
		if shifted != "a" {
			t.Errorf("Expected 'a', got %v", shifted)
		}
	})
}
