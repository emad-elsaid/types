package types

import (
	"reflect"
	"testing"
)

func TestZip(t *testing.T) {
	t.Run("equal length slices", func(t *testing.T) {
		a := Slice[int]{1, 2, 3}
		b := Slice[string]{"a", "b", "c"}
		result := Zip(a, b)

		expected := [][2]any{
			{1, "a"},
			{2, "b"},
			{3, "c"},
		}

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Zip() = %v, want %v", result, expected)
		}
	})

	t.Run("first slice shorter", func(t *testing.T) {
		a := Slice[int]{1, 2}
		b := Slice[string]{"a", "b", "c", "d"}
		result := Zip(a, b)

		expected := [][2]any{
			{1, "a"},
			{2, "b"},
		}

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Zip() = %v, want %v", result, expected)
		}
	})

	t.Run("second slice shorter", func(t *testing.T) {
		a := Slice[int]{1, 2, 3, 4, 5}
		b := Slice[string]{"a", "b"}
		result := Zip(a, b)

		expected := [][2]any{
			{1, "a"},
			{2, "b"},
		}

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Zip() = %v, want %v", result, expected)
		}
	})

	t.Run("first slice empty", func(t *testing.T) {
		a := Slice[int]{}
		b := Slice[string]{"a", "b", "c"}
		result := Zip(a, b)

		expected := [][2]any{}

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Zip() = %v, want %v", result, expected)
		}
	})

	t.Run("second slice empty", func(t *testing.T) {
		a := Slice[int]{1, 2, 3}
		b := Slice[string]{}
		result := Zip(a, b)

		expected := [][2]any{}

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Zip() = %v, want %v", result, expected)
		}
	})

	t.Run("both slices empty", func(t *testing.T) {
		a := Slice[int]{}
		b := Slice[string]{}
		result := Zip(a, b)

		expected := [][2]any{}

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Zip() = %v, want %v", result, expected)
		}
	})

	t.Run("same type slices", func(t *testing.T) {
		a := Slice[int]{1, 2, 3}
		b := Slice[int]{10, 20, 30}
		result := Zip(a, b)

		expected := [][2]any{
			{1, 10},
			{2, 20},
			{3, 30},
		}

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Zip() = %v, want %v", result, expected)
		}
	})

	t.Run("different comparable types", func(t *testing.T) {
		a := Slice[float64]{1.5, 2.5, 3.5}
		b := Slice[bool]{true, false, true}
		result := Zip(a, b)

		expected := [][2]any{
			{1.5, true},
			{2.5, false},
			{3.5, true},
		}

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Zip() = %v, want %v", result, expected)
		}
	})

	t.Run("single element slices", func(t *testing.T) {
		a := Slice[int]{42}
		b := Slice[string]{"answer"}
		result := Zip(a, b)

		expected := [][2]any{
			{42, "answer"},
		}

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Zip() = %v, want %v", result, expected)
		}
	})
}

// Benchmark tests
func BenchmarkZipEqualLength(b *testing.B) {
	a := make(Slice[int], 1000)
	s := make(Slice[string], 1000)
	for i := range a {
		a[i] = i
		s[i] = "test"
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Zip(a, s)
	}
}

func BenchmarkZipDifferentLength(b *testing.B) {
	a := make(Slice[int], 100)
	s := make(Slice[string], 1000)
	for i := range a {
		a[i] = i
	}
	for i := range s {
		s[i] = "test"
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Zip(a, s)
	}
}
