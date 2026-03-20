package types

import (
	"testing"
)

func TestFlatten(t *testing.T) {
	tests := []struct {
		name     string
		input    [][]int
		expected []int
	}{
		{
			name:     "basic flatten",
			input:    [][]int{{1, 2}, {3, 4}, {5}},
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			name:     "empty nested slices",
			input:    [][]int{{}, {1, 2}, {}, {3}},
			expected: []int{1, 2, 3},
		},
		{
			name:     "single nested slice",
			input:    [][]int{{1, 2, 3}},
			expected: []int{1, 2, 3},
		},
		{
			name:     "empty outer slice",
			input:    [][]int{},
			expected: []int{},
		},
		{
			name:     "all empty nested slices",
			input:    [][]int{{}, {}, {}},
			expected: []int{},
		},
		{
			name:     "single elements in nested slices",
			input:    [][]int{{1}, {2}, {3}},
			expected: []int{1, 2, 3},
		},
		{
			name:     "varying lengths",
			input:    [][]int{{1}, {2, 3, 4}, {5, 6}},
			expected: []int{1, 2, 3, 4, 5, 6},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Flatten(tt.input)
			if !Slice[int](result).IsEq(Slice[int](tt.expected)) {
				t.Errorf("Flatten() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestFlattenStrings(t *testing.T) {
	input := [][]string{
		{"hello", "world"},
		{"foo", "bar"},
		{"baz"},
	}
	expected := []string{"hello", "world", "foo", "bar", "baz"}
	result := Flatten(input)

	if !Slice[string](result).IsEq(Slice[string](expected)) {
		t.Errorf("Flatten() = %v, want %v", result, expected)
	}
}

func TestFlattenPreservesOrder(t *testing.T) {
	input := [][]int{
		{9, 8, 7},
		{6, 5},
		{4, 3, 2, 1},
	}
	expected := []int{9, 8, 7, 6, 5, 4, 3, 2, 1}
	result := Flatten(input)

	if !Slice[int](result).IsEq(Slice[int](expected)) {
		t.Errorf("Flatten() should preserve order: got %v, want %v", result, expected)
	}
}

