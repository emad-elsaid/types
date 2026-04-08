package types

import "testing"

func TestSliceCycle(t *testing.T) {
	t.Run("cycles through slice elements multiple times", func(t *testing.T) {
		s := Slice[int]{1, 2, 3}
		var result []int

		s.Cycle(2, func(v int) {
			result = append(result, v)
		})

		expected := []int{1, 2, 3, 1, 2, 3}
		if len(result) != len(expected) {
			t.Errorf("Expected length %d but got %d", len(expected), len(result))
		}

		for i := range result {
			if result[i] != expected[i] {
				t.Errorf("At index %d: expected %d but got %d", i, expected[i], result[i])
			}
		}
	})

	t.Run("cycles zero times produces no output", func(t *testing.T) {
		s := Slice[int]{1, 2, 3}
		count := 0

		s.Cycle(0, func(v int) {
			count++
		})

		if count != 0 {
			t.Errorf("Expected 0 calls but got %d", count)
		}
	})

	t.Run("cycles empty slice", func(t *testing.T) {
		s := Slice[int]{}
		count := 0

		s.Cycle(3, func(v int) {
			count++
		})

		if count != 0 {
			t.Errorf("Expected 0 calls for empty slice but got %d", count)
		}
	})

	t.Run("cycles once", func(t *testing.T) {
		s := Slice[string]{"a", "b"}
		var result []string

		s.Cycle(1, func(v string) {
			result = append(result, v)
		})

		expected := []string{"a", "b"}
		if len(result) != len(expected) {
			t.Errorf("Expected length %d but got %d", len(expected), len(result))
		}

		for i := range result {
			if result[i] != expected[i] {
				t.Errorf("At index %d: expected %s but got %s", i, expected[i], result[i])
			}
		}
	})

	t.Run("cycles with side effects", func(t *testing.T) {
		s := Slice[int]{10, 20}
		sum := 0

		s.Cycle(3, func(v int) {
			sum += v
		})

		// (10 + 20) * 3 = 90
		expected := 90
		if sum != expected {
			t.Errorf("Expected sum %d but got %d", expected, sum)
		}
	})
}
