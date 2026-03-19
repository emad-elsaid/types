package types

import (
	"testing"
)

func TestSliceRotate(t *testing.T) {
	tests := []struct {
		name     string
		input    Slice[int]
		count    int
		expected Slice[int]
	}{
		{
			name:     "rotate left by 2",
			input:    Slice[int]{1, 2, 3, 4, 5},
			count:    2,
			expected: Slice[int]{3, 4, 5, 1, 2},
		},
		{
			name:     "rotate right by 2",
			input:    Slice[int]{1, 2, 3, 4, 5},
			count:    -2,
			expected: Slice[int]{4, 5, 1, 2, 3},
		},
		{
			name:     "rotate left by 1",
			input:    Slice[int]{1, 2, 3, 4, 5},
			count:    1,
			expected: Slice[int]{2, 3, 4, 5, 1},
		},
		{
			name:     "rotate right by 1",
			input:    Slice[int]{1, 2, 3, 4, 5},
			count:    -1,
			expected: Slice[int]{5, 1, 2, 3, 4},
		},
		{
			name:     "rotate by 0 (no change)",
			input:    Slice[int]{1, 2, 3, 4, 5},
			count:    0,
			expected: Slice[int]{1, 2, 3, 4, 5},
		},
		{
			name:     "rotate by length (full rotation, no change)",
			input:    Slice[int]{1, 2, 3, 4, 5},
			count:    5,
			expected: Slice[int]{1, 2, 3, 4, 5},
		},
		{
			name:     "rotate by negative length (full rotation, no change)",
			input:    Slice[int]{1, 2, 3, 4, 5},
			count:    -5,
			expected: Slice[int]{1, 2, 3, 4, 5},
		},
		{
			name:     "rotate by more than length",
			input:    Slice[int]{1, 2, 3, 4, 5},
			count:    7, // equivalent to rotating by 2
			expected: Slice[int]{3, 4, 5, 1, 2},
		},
		{
			name:     "rotate by negative more than length",
			input:    Slice[int]{1, 2, 3, 4, 5},
			count:    -7, // equivalent to rotating by -2
			expected: Slice[int]{4, 5, 1, 2, 3},
		},
		{
			name:     "empty slice",
			input:    Slice[int]{},
			count:    3,
			expected: Slice[int]{},
		},
		{
			name:     "single element",
			input:    Slice[int]{42},
			count:    1,
			expected: Slice[int]{42},
		},
		{
			name:     "two elements rotate left",
			input:    Slice[int]{1, 2},
			count:    1,
			expected: Slice[int]{2, 1},
		},
		{
			name:     "two elements rotate right",
			input:    Slice[int]{1, 2},
			count:    -1,
			expected: Slice[int]{2, 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.Rotate(tt.count)

			// Check length
			if len(result) != len(tt.expected) {
				t.Errorf("expected length %d, got %d", len(tt.expected), len(result))
				return
			}

			// Check each element
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("at index %d: expected %v, got %v", i, tt.expected[i], result[i])
				}
			}

			// Verify original slice is not modified
			if len(tt.input) > 0 {
				// Just check the first element to ensure immutability
				firstElement := tt.input[0]
				_ = tt.input.Rotate(tt.count)
				if tt.input[0] != firstElement {
					t.Errorf("original slice was modified")
				}
			}
		})
	}
}

func TestSliceRotateWithStrings(t *testing.T) {
	input := Slice[string]{"a", "b", "c", "d", "e"}

	tests := []struct {
		name     string
		count    int
		expected Slice[string]
	}{
		{
			name:     "rotate strings left by 2",
			count:    2,
			expected: Slice[string]{"c", "d", "e", "a", "b"},
		},
		{
			name:     "rotate strings right by 2",
			count:    -2,
			expected: Slice[string]{"d", "e", "a", "b", "c"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := input.Rotate(tt.count)

			if !result.IsEq(tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
