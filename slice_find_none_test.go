package types

import "testing"

func TestSliceNone(t *testing.T) {
	tests := []struct {
		name      string
		slice     Slice[int]
		predicate func(int) bool
		expected  bool
	}{
		{
			name:      "empty slice",
			slice:     Slice[int]{},
			predicate: func(x int) bool { return x > 0 },
			expected:  true, // None returns true for empty slices
		},
		{
			name:      "no elements satisfy predicate",
			slice:     Slice[int]{1, 2, 3, 4, 5},
			predicate: func(x int) bool { return x > 10 },
			expected:  true,
		},
		{
			name:      "some elements satisfy predicate",
			slice:     Slice[int]{1, 2, 3, 4, 5},
			predicate: func(x int) bool { return x == 3 },
			expected:  false,
		},
		{
			name:      "all elements satisfy predicate",
			slice:     Slice[int]{2, 4, 6, 8},
			predicate: func(x int) bool { return x%2 == 0 },
			expected:  false,
		},
		{
			name:      "single element satisfies",
			slice:     Slice[int]{5},
			predicate: func(x int) bool { return x == 5 },
			expected:  false,
		},
		{
			name:      "single element does not satisfy",
			slice:     Slice[int]{5},
			predicate: func(x int) bool { return x == 10 },
			expected:  true,
		},
		{
			name:      "negative numbers - none negative",
			slice:     Slice[int]{1, 2, 3, 4},
			predicate: func(x int) bool { return x < 0 },
			expected:  true,
		},
		{
			name:      "negative numbers - has negative",
			slice:     Slice[int]{-1, 2, 3, 4},
			predicate: func(x int) bool { return x < 0 },
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.slice.None(tt.predicate)
			if result != tt.expected {
				t.Errorf("None() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestSliceFind(t *testing.T) {
	tests := []struct {
		name          string
		slice         Slice[int]
		predicate     func(int) bool
		expectedValue int
		expectedFound bool
	}{
		{
			name:          "empty slice",
			slice:         Slice[int]{},
			predicate:     func(x int) bool { return x > 0 },
			expectedValue: 0,
			expectedFound: false,
		},
		{
			name:          "find first even number",
			slice:         Slice[int]{1, 3, 4, 5, 6},
			predicate:     func(x int) bool { return x%2 == 0 },
			expectedValue: 4,
			expectedFound: true,
		},
		{
			name:          "find first odd number",
			slice:         Slice[int]{2, 4, 5, 6, 7},
			predicate:     func(x int) bool { return x%2 != 0 },
			expectedValue: 5,
			expectedFound: true,
		},
		{
			name:          "no element satisfies predicate",
			slice:         Slice[int]{1, 2, 3, 4, 5},
			predicate:     func(x int) bool { return x > 10 },
			expectedValue: 0,
			expectedFound: false,
		},
		{
			name:          "first element satisfies",
			slice:         Slice[int]{10, 2, 3, 4, 5},
			predicate:     func(x int) bool { return x > 5 },
			expectedValue: 10,
			expectedFound: true,
		},
		{
			name:          "last element satisfies",
			slice:         Slice[int]{1, 2, 3, 4, 10},
			predicate:     func(x int) bool { return x > 5 },
			expectedValue: 10,
			expectedFound: true,
		},
		{
			name:          "multiple elements satisfy - returns first",
			slice:         Slice[int]{1, 5, 8, 10, 12},
			predicate:     func(x int) bool { return x > 4 },
			expectedValue: 5,
			expectedFound: true,
		},
		{
			name:          "single element satisfies",
			slice:         Slice[int]{42},
			predicate:     func(x int) bool { return x == 42 },
			expectedValue: 42,
			expectedFound: true,
		},
		{
			name:          "single element does not satisfy",
			slice:         Slice[int]{42},
			predicate:     func(x int) bool { return x == 10 },
			expectedValue: 0,
			expectedFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, found := tt.slice.Find(tt.predicate)
			if found != tt.expectedFound {
				t.Errorf("Find() found = %v, expected %v", found, tt.expectedFound)
			}
			if value != tt.expectedValue {
				t.Errorf("Find() value = %v, expected %v", value, tt.expectedValue)
			}
		})
	}
}

func TestSliceFindString(t *testing.T) {
	tests := []struct {
		name          string
		slice         Slice[string]
		predicate     func(string) bool
		expectedValue string
		expectedFound bool
	}{
		{
			name:          "find string starting with 'h'",
			slice:         Slice[string]{"apple", "banana", "hello", "world"},
			predicate:     func(s string) bool { return len(s) > 0 && s[0] == 'h' },
			expectedValue: "hello",
			expectedFound: true,
		},
		{
			name:          "find string with length > 6",
			slice:         Slice[string]{"cat", "dog", "elephant", "fox"},
			predicate:     func(s string) bool { return len(s) > 6 },
			expectedValue: "elephant",
			expectedFound: true,
		},
		{
			name:          "no string matches",
			slice:         Slice[string]{"cat", "dog", "fox"},
			predicate:     func(s string) bool { return len(s) > 10 },
			expectedValue: "",
			expectedFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, found := tt.slice.Find(tt.predicate)
			if found != tt.expectedFound {
				t.Errorf("Find() found = %v, expected %v", found, tt.expectedFound)
			}
			if value != tt.expectedValue {
				t.Errorf("Find() value = %q, expected %q", value, tt.expectedValue)
			}
		})
	}
}
