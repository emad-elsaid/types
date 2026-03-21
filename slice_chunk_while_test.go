package types

import (
	"reflect"
	"testing"
)

func TestSlice_ChunkWhile(t *testing.T) {
	tests := []struct {
		name      string
		slice     Slice[int]
		predicate func(int, int) bool
		expected  []Slice[int]
	}{
		{
			name:      "consecutive numbers",
			slice:     Slice[int]{1, 2, 4, 5, 7, 9},
			predicate: func(a, b int) bool { return b-a == 1 },
			expected:  []Slice[int]{{1, 2}, {4, 5}, {7}, {9}},
		},
		{
			name:      "all consecutive",
			slice:     Slice[int]{1, 2, 3, 4, 5},
			predicate: func(a, b int) bool { return b-a == 1 },
			expected:  []Slice[int]{{1, 2, 3, 4, 5}},
		},
		{
			name:      "none consecutive",
			slice:     Slice[int]{1, 3, 5, 7, 9},
			predicate: func(a, b int) bool { return b-a == 1 },
			expected:  []Slice[int]{{1}, {3}, {5}, {7}, {9}},
		},
		{
			name:      "empty slice",
			slice:     Slice[int]{},
			predicate: func(a, b int) bool { return b-a == 1 },
			expected:  []Slice[int]{},
		},
		{
			name:      "single element",
			slice:     Slice[int]{42},
			predicate: func(a, b int) bool { return b-a == 1 },
			expected:  []Slice[int]{{42}},
		},
		{
			name:      "two elements - match",
			slice:     Slice[int]{1, 2},
			predicate: func(a, b int) bool { return b-a == 1 },
			expected:  []Slice[int]{{1, 2}},
		},
		{
			name:      "two elements - no match",
			slice:     Slice[int]{1, 5},
			predicate: func(a, b int) bool { return b-a == 1 },
			expected:  []Slice[int]{{1}, {5}},
		},
		{
			name:      "same values",
			slice:     Slice[int]{3, 3, 3, 5, 5, 7},
			predicate: func(a, b int) bool { return a == b },
			expected:  []Slice[int]{{3, 3, 3}, {5, 5}, {7}},
		},
		{
			name:      "ascending values",
			slice:     Slice[int]{1, 3, 2, 5, 4, 7, 6},
			predicate: func(a, b int) bool { return a < b },
			expected:  []Slice[int]{{1, 3}, {2, 5}, {4, 7}, {6}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.slice.ChunkWhile(tt.predicate)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ChunkWhile() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestSlice_ChunkWhile_Strings(t *testing.T) {
	tests := []struct {
		name      string
		slice     Slice[string]
		predicate func(string, string) bool
		expected  []Slice[string]
	}{
		{
			name:  "same length strings",
			slice: Slice[string]{"a", "b", "cc", "dd", "eee", "f"},
			predicate: func(a, b string) bool {
				return len(a) == len(b)
			},
			expected: []Slice[string]{{"a", "b"}, {"cc", "dd"}, {"eee"}, {"f"}},
		},
		{
			name:  "alphabetical order",
			slice: Slice[string]{"apple", "banana", "cherry", "ant", "bear"},
			predicate: func(a, b string) bool {
				return a < b
			},
			expected: []Slice[string]{{"apple", "banana", "cherry"}, {"ant", "bear"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.slice.ChunkWhile(tt.predicate)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ChunkWhile() = %v, expected %v", result, tt.expected)
			}
		})
	}
}
