package types

import (
	"reflect"
	"testing"
)

func TestSlice_Tally(t *testing.T) {
	tests := []struct {
		name     string
		slice    Slice[int]
		expected map[int]int
	}{
		{
			name:     "empty slice",
			slice:    Slice[int]{},
			expected: map[int]int{},
		},
		{
			name:     "single element",
			slice:    Slice[int]{1},
			expected: map[int]int{1: 1},
		},
		{
			name:     "all unique elements",
			slice:    Slice[int]{1, 2, 3, 4, 5},
			expected: map[int]int{1: 1, 2: 1, 3: 1, 4: 1, 5: 1},
		},
		{
			name:     "all same elements",
			slice:    Slice[int]{5, 5, 5, 5},
			expected: map[int]int{5: 4},
		},
		{
			name:     "mixed frequencies",
			slice:    Slice[int]{1, 2, 2, 3, 3, 3},
			expected: map[int]int{1: 1, 2: 2, 3: 3},
		},
		{
			name:     "negative numbers",
			slice:    Slice[int]{-1, -1, 0, 1, 1, 1},
			expected: map[int]int{-1: 2, 0: 1, 1: 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.slice.Tally()
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Tally() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestSlice_TallyString(t *testing.T) {
	tests := []struct {
		name     string
		slice    Slice[string]
		expected map[string]int
	}{
		{
			name:     "empty slice",
			slice:    Slice[string]{},
			expected: map[string]int{},
		},
		{
			name:     "word frequency",
			slice:    Slice[string]{"apple", "banana", "apple", "cherry", "banana", "apple"},
			expected: map[string]int{"apple": 3, "banana": 2, "cherry": 1},
		},
		{
			name:     "empty strings",
			slice:    Slice[string]{"", "a", "", "b", ""},
			expected: map[string]int{"": 3, "a": 1, "b": 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.slice.Tally()
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Tally() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// Benchmark to ensure performance is reasonable
func BenchmarkSlice_Tally(b *testing.B) {
	// Create a slice with 1000 elements, mixed frequencies
	slice := make(Slice[int], 1000)
	for i := range slice {
		slice[i] = i % 100 // 100 unique values, each appearing 10 times
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = slice.Tally()
	}
}

// Test that Tally doesn't modify the original slice
func TestSlice_TallyImmutable(t *testing.T) {
	original := Slice[int]{1, 2, 3, 2, 1}
	originalCopy := make(Slice[int], len(original))
	copy(originalCopy, original)

	_ = original.Tally()

	if !reflect.DeepEqual(original, originalCopy) {
		t.Errorf("Tally() modified the original slice: got %v, expected %v", original, originalCopy)
	}
}

// Test integration with other Slice methods
func TestSlice_TallyWithOtherMethods(t *testing.T) {
	// Example: count frequencies after filtering
	slice := Slice[int]{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	evenOnly := slice.KeepIf(func(x int) bool { return x%2 == 0 })
	tally := evenOnly.Tally()

	expected := map[int]int{2: 1, 4: 1, 6: 1, 8: 1, 10: 1}
	if !reflect.DeepEqual(tally, expected) {
		t.Errorf("Tally() after KeepIf() = %v, expected %v", tally, expected)
	}

	// Example: count frequencies after mapping
	slice2 := Slice[int]{1, 2, 3, 4, 5}
	doubled := slice2.Map(func(x int) int { return x * 2 })
	tally2 := doubled.Tally()

	expected2 := map[int]int{2: 1, 4: 1, 6: 1, 8: 1, 10: 1}
	if !reflect.DeepEqual(tally2, expected2) {
		t.Errorf("Tally() after Map() = %v, expected %v", tally2, expected2)
	}
}
