package types

import "testing"

func TestSliceSum(t *testing.T) {
	t.Run("int slice", func(t *testing.T) {
		s := Slice[int]{1, 2, 3, 4, 5}
		result := SliceSum(s)
		expected := 15
		if result != expected {
			t.Errorf("Expected %d, got %d", expected, result)
		}
	})

	t.Run("float64 slice", func(t *testing.T) {
		s := Slice[float64]{1.5, 2.5, 3.0}
		result := SliceSum(s)
		expected := 7.0
		if result != expected {
			t.Errorf("Expected %f, got %f", expected, result)
		}
	})

	t.Run("empty slice", func(t *testing.T) {
		s := Slice[int]{}
		result := SliceSum(s)
		expected := 0
		if result != expected {
			t.Errorf("Expected %d for empty slice, got %d", expected, result)
		}
	})

	t.Run("negative numbers", func(t *testing.T) {
		s := Slice[int]{-1, -2, -3, 4, 5}
		result := SliceSum(s)
		expected := 3
		if result != expected {
			t.Errorf("Expected %d, got %d", expected, result)
		}
	})

	t.Run("uint slice", func(t *testing.T) {
		s := Slice[uint]{1, 2, 3}
		result := SliceSum(s)
		var expected uint = 6
		if result != expected {
			t.Errorf("Expected %d, got %d", expected, result)
		}
	})

	t.Run("int64 slice", func(t *testing.T) {
		s := Slice[int64]{100, 200, 300}
		result := SliceSum(s)
		var expected int64 = 600
		if result != expected {
			t.Errorf("Expected %d, got %d", expected, result)
		}
	})

	t.Run("float32 slice", func(t *testing.T) {
		s := Slice[float32]{1.1, 2.2, 3.3}
		result := SliceSum(s)
		expected := float32(6.6)
		// Use approximate equality for floats
		if result < expected-0.01 || result > expected+0.01 {
			t.Errorf("Expected approximately %f, got %f", expected, result)
		}
	})

	t.Run("single element", func(t *testing.T) {
		s := Slice[int]{42}
		result := SliceSum(s)
		expected := 42
		if result != expected {
			t.Errorf("Expected %d, got %d", expected, result)
		}
	})

	t.Run("zeros", func(t *testing.T) {
		s := Slice[int]{0, 0, 0}
		result := SliceSum(s)
		expected := 0
		if result != expected {
			t.Errorf("Expected %d, got %d", expected, result)
		}
	})
}
