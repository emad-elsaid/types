
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

import (
	"testing"
)

func TestSliceCompact(t *testing.T) {
	t.Run("removes zero values from int slice", func(t *testing.T) {
		s := Slice[int]{1, 0, 2, 0, 3, 0}
		result := s.Compact()
		expected := Slice[int]{1, 2, 3}
		if !result.IsEq(expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("removes empty strings from string slice", func(t *testing.T) {
		s := Slice[string]{"hello", "", "world", "", "!"}
		result := s.Compact()
		expected := Slice[string]{"hello", "world", "!"}
		if !result.IsEq(expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("removes false from bool slice", func(t *testing.T) {
		s := Slice[bool]{true, false, true, false, true}
		result := s.Compact()
		expected := Slice[bool]{true, true, true}
		if !result.IsEq(expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("returns empty slice when all values are zero", func(t *testing.T) {
		s := Slice[int]{0, 0, 0, 0}
		result := s.Compact()
		if len(result) != 0 {
			t.Errorf("Expected empty slice, got %v", result)
		}
	})

	t.Run("returns same slice when no zero values", func(t *testing.T) {
		s := Slice[int]{1, 2, 3, 4, 5}
		result := s.Compact()
		if !result.IsEq(s) {
			t.Errorf("Expected %v, got %v", s, result)
		}
	})

	t.Run("returns empty slice when input is empty", func(t *testing.T) {
		s := Slice[int]{}
		result := s.Compact()
		if len(result) != 0 {
			t.Errorf("Expected empty slice, got %v", result)
		}
	})

	t.Run("works with float64", func(t *testing.T) {
		s := Slice[float64]{1.5, 0.0, 2.5, 0.0, 3.5}
		result := s.Compact()
		expected := Slice[float64]{1.5, 2.5, 3.5}
		if !result.IsEq(expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("removes zero runes from rune slice", func(t *testing.T) {
		s := Slice[rune]{'a', 0, 'b', 0, 'c'}
		result := s.Compact()
		expected := Slice[rune]{'a', 'b', 'c'}
		if !result.IsEq(expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("preserves order of non-zero elements", func(t *testing.T) {
		s := Slice[int]{5, 0, 4, 0, 3, 0, 2, 0, 1}
		result := s.Compact()
		expected := Slice[int]{5, 4, 3, 2, 1}
		if !result.IsEq(expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("single zero value returns empty slice", func(t *testing.T) {
		s := Slice[int]{0}
		result := s.Compact()
		if len(result) != 0 {
			t.Errorf("Expected empty slice, got %v", result)
		}
	})

	t.Run("single non-zero value returns same value", func(t *testing.T) {
		s := Slice[int]{42}
		result := s.Compact()
		expected := Slice[int]{42}
		if !result.IsEq(expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})
}


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


import (
	"reflect"
	"testing"
)

func TestSliceReduce_Basic(t *testing.T) {
	tests := []struct {
		name  string
		input Slice[int]
		block func(int) bool
		want  Slice[int]
	}{
		{
			name:  "keep even numbers",
			input: Slice[int]{1, 2, 3, 4, 5, 6},
			block: func(n int) bool { return n%2 == 0 },
			want:  Slice[int]{2, 4, 6},
		},
		{
			name:  "keep odd numbers",
			input: Slice[int]{1, 2, 3, 4, 5, 6},
			block: func(n int) bool { return n%2 != 0 },
			want:  Slice[int]{1, 3, 5},
		},
		{
			name:  "keep numbers greater than 3",
			input: Slice[int]{1, 2, 3, 4, 5, 6},
			block: func(n int) bool { return n > 3 },
			want:  Slice[int]{4, 5, 6},
		},
		{
			name:  "keep none",
			input: Slice[int]{1, 2, 3},
			block: func(n int) bool { return false },
			want:  Slice[int]{},
		},
		{
			name:  "keep all",
			input: Slice[int]{1, 2, 3},
			block: func(n int) bool { return true },
			want:  Slice[int]{1, 2, 3},
		},
		{
			name:  "empty slice",
			input: Slice[int]{},
			block: func(n int) bool { return true },
			want:  Slice[int]{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.input.Reduce(tt.block)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Reduce() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSliceReduce_String(t *testing.T) {
	tests := []struct {
		name  string
		input Slice[string]
		block func(string) bool
		want  Slice[string]
	}{
		{
			name:  "keep long strings",
			input: Slice[string]{"a", "hello", "hi", "world", "go"},
			block: func(s string) bool { return len(s) > 2 },
			want:  Slice[string]{"hello", "world"},
		},
		{
			name:  "keep strings starting with h",
			input: Slice[string]{"hello", "world", "hi", "golang"},
			block: func(s string) bool { return len(s) > 0 && s[0] == 'h' },
			want:  Slice[string]{"hello", "hi"},
		},
		{
			name:  "empty result",
			input: Slice[string]{"a", "b", "c"},
			block: func(s string) bool { return len(s) > 5 },
			want:  Slice[string]{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.input.Reduce(tt.block)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Reduce() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSliceReduce_Struct(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	people := Slice[Person]{
		{Name: "Alice", Age: 25},
		{Name: "Bob", Age: 30},
		{Name: "Charlie", Age: 20},
		{Name: "Diana", Age: 35},
	}

	tests := []struct {
		name  string
		input Slice[Person]
		block func(Person) bool
		want  Slice[Person]
	}{
		{
			name:  "keep adults over 25",
			input: people,
			block: func(p Person) bool { return p.Age > 25 },
			want: Slice[Person]{
				{Name: "Bob", Age: 30},
				{Name: "Diana", Age: 35},
			},
		},
		{
			name:  "keep names starting with A",
			input: people,
			block: func(p Person) bool { return len(p.Name) > 0 && p.Name[0] == 'A' },
			want: Slice[Person]{
				{Name: "Alice", Age: 25},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.input.Reduce(tt.block)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Reduce() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSliceReduce_ChainedOperations(t *testing.T) {
	// Test that Reduce works well in method chains
	input := Slice[int]{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// Chain: keep evens, then keep those > 4
	result := input.
		Reduce(func(n int) bool { return n%2 == 0 }).
		Reduce(func(n int) bool { return n > 4 })

	want := Slice[int]{6, 8, 10}

	if !reflect.DeepEqual(result, want) {
		t.Errorf("Chained Reduce() = %v, want %v", result, want)
	}
}

func TestSliceReduce_PreservesOrder(t *testing.T) {
	input := Slice[int]{5, 1, 4, 2, 3}

	// Keep odds - should preserve original order
	result := input.Reduce(func(n int) bool { return n%2 != 0 })

	want := Slice[int]{5, 1, 3}

	if !reflect.DeepEqual(result, want) {
		t.Errorf("Reduce() = %v, want %v (order should be preserved)", result, want)
	}
}

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

import (
	"testing"
)

// TestSliceSelectUntil_Found tests SelectUntil when condition is met
func TestSliceSelectUntil_Found(t *testing.T) {
	a := Slice[int]{1, 2, 3, 4, 5, 6}
	result := a.SelectUntil(func(e int) bool {
		return e == 3
	})
	expected := Slice[int]{1, 2}
	AssertSlicesEquals(t, expected, result)
}

// TestSliceSelectUntil_NotFound tests SelectUntil when condition is never met
func TestSliceSelectUntil_NotFound(t *testing.T) {
	a := Slice[int]{1, 2, 3, 4, 5}
	result := a.SelectUntil(func(e int) bool {
		return e == 10
	})
	// Should return the entire slice when condition is never met
	AssertSlicesEquals(t, a, result)
}

// TestSliceSelectUntil_FirstElement tests SelectUntil when first element matches
func TestSliceSelectUntil_FirstElement(t *testing.T) {
	a := Slice[int]{1, 2, 3, 4}
	result := a.SelectUntil(func(e int) bool {
		return e == 1
	})
	expected := Slice[int]{}
	AssertSlicesEquals(t, expected, result)
}

// TestSliceSelectUntil_LastElement tests SelectUntil when last element matches
func TestSliceSelectUntil_LastElement(t *testing.T) {
	a := Slice[int]{1, 2, 3, 4, 5}
	result := a.SelectUntil(func(e int) bool {
		return e == 5
	})
	expected := Slice[int]{1, 2, 3, 4}
	AssertSlicesEquals(t, expected, result)
}

// TestSliceSelectUntil_Empty tests SelectUntil on empty slice
func TestSliceSelectUntil_Empty(t *testing.T) {
	a := Slice[int]{}
	result := a.SelectUntil(func(e int) bool {
		return e == 1
	})
	expected := Slice[int]{}
	AssertSlicesEquals(t, expected, result)
}

// TestSliceSelectUntil_Strings tests SelectUntil with string type
func TestSliceSelectUntil_Strings(t *testing.T) {
	a := Slice[string]{"apple", "banana", "cherry", "date"}
	result := a.SelectUntil(func(e string) bool {
		return e == "cherry"
	})
	expected := Slice[string]{"apple", "banana"}
	AssertSlicesEquals(t, expected, result)
}

// TestSliceSelectUntil_ComplexCondition tests SelectUntil with complex predicate
func TestSliceSelectUntil_ComplexCondition(t *testing.T) {
	a := Slice[int]{1, 3, 5, 8, 10, 12}
	result := a.SelectUntil(func(e int) bool {
		return e%2 == 0 // First even number
	})
	expected := Slice[int]{1, 3, 5}
	AssertSlicesEquals(t, expected, result)
}

// TestSliceSelectUntil_Chaining tests SelectUntil with method chaining
func TestSliceSelectUntil_Chaining(t *testing.T) {
	a := Slice[int]{1, 2, 3, 4, 5, 6, 7, 8}
	result := a.SelectUntil(func(e int) bool {
		return e > 5
	}).Map(func(e int) int {
		return e * 2
	})
	expected := Slice[int]{2, 4, 6, 8, 10}
	AssertSlicesEquals(t, expected, result)
}


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

import (
	"testing"
)

func TestSliceTakeWhile(t *testing.T) {
	t.Run("takes elements while predicate is true", func(t *testing.T) {
		s := Slice[int]{1, 2, 3, 4, 5, 6}
		result := s.TakeWhile(func(x int) bool { return x < 4 })
		expected := Slice[int]{1, 2, 3}
		if !result.IsEq(expected) {
			t.Errorf("TakeWhile: expected %v, got %v", expected, result)
		}
	})

	t.Run("returns empty slice when first element fails predicate", func(t *testing.T) {
		s := Slice[int]{5, 6, 7, 8}
		result := s.TakeWhile(func(x int) bool { return x < 5 })
		expected := Slice[int]{}
		if !result.IsEq(expected) {
			t.Errorf("TakeWhile: expected %v, got %v", expected, result)
		}
	})

	t.Run("returns all elements when predicate always true", func(t *testing.T) {
		s := Slice[int]{1, 2, 3, 4}
		result := s.TakeWhile(func(x int) bool { return x < 10 })
		if !result.IsEq(s) {
			t.Errorf("TakeWhile: expected %v, got %v", s, result)
		}
	})

	t.Run("returns empty slice for empty input", func(t *testing.T) {
		s := Slice[int]{}
		result := s.TakeWhile(func(x int) bool { return x < 5 })
		expected := Slice[int]{}
		if !result.IsEq(expected) {
			t.Errorf("TakeWhile: expected %v, got %v", expected, result)
		}
	})

	t.Run("works with strings", func(t *testing.T) {
		s := Slice[string]{"apple", "banana", "cherry", "date"}
		result := s.TakeWhile(func(x string) bool { return len(x) < 6 })
		expected := Slice[string]{"apple"}
		if !result.IsEq(expected) {
			t.Errorf("TakeWhile: expected %v, got %v", expected, result)
		}
	})

	t.Run("stops at first false predicate", func(t *testing.T) {
		s := Slice[int]{2, 4, 6, 3, 8, 10}
		result := s.TakeWhile(func(x int) bool { return x%2 == 0 })
		expected := Slice[int]{2, 4, 6}
		if !result.IsEq(expected) {
			t.Errorf("TakeWhile: expected %v, got %v", expected, result)
		}
	})
}

func TestSliceDropWhile(t *testing.T) {
	t.Run("drops elements while predicate is true", func(t *testing.T) {
		s := Slice[int]{1, 2, 3, 4, 5, 6}
		result := s.DropWhile(func(x int) bool { return x < 4 })
		expected := Slice[int]{4, 5, 6}
		if !result.IsEq(expected) {
			t.Errorf("DropWhile: expected %v, got %v", expected, result)
		}
	})

	t.Run("returns original slice when first element fails predicate", func(t *testing.T) {
		s := Slice[int]{5, 6, 7, 8}
		result := s.DropWhile(func(x int) bool { return x < 5 })
		if !result.IsEq(s) {
			t.Errorf("DropWhile: expected %v, got %v", s, result)
		}
	})

	t.Run("returns empty slice when predicate always true", func(t *testing.T) {
		s := Slice[int]{1, 2, 3, 4}
		result := s.DropWhile(func(x int) bool { return x < 10 })
		expected := Slice[int]{}
		if !result.IsEq(expected) {
			t.Errorf("DropWhile: expected %v, got %v", expected, result)
		}
	})

	t.Run("returns empty slice for empty input", func(t *testing.T) {
		s := Slice[int]{}
		result := s.DropWhile(func(x int) bool { return x < 5 })
		expected := Slice[int]{}
		if !result.IsEq(expected) {
			t.Errorf("DropWhile: expected %v, got %v", expected, result)
		}
	})

	t.Run("works with strings", func(t *testing.T) {
		s := Slice[string]{"apple", "banana", "cherry", "date"}
		result := s.DropWhile(func(x string) bool { return len(x) < 6 })
		expected := Slice[string]{"banana", "cherry", "date"}
		if !result.IsEq(expected) {
			t.Errorf("DropWhile: expected %v, got %v", expected, result)
		}
	})

	t.Run("stops dropping at first false predicate", func(t *testing.T) {
		s := Slice[int]{2, 4, 6, 3, 8, 10}
		result := s.DropWhile(func(x int) bool { return x%2 == 0 })
		expected := Slice[int]{3, 8, 10}
		if !result.IsEq(expected) {
			t.Errorf("DropWhile: expected %v, got %v", expected, result)
		}
	})
}

func TestTakeWhileDropWhileComplement(t *testing.T) {
	t.Run("TakeWhile and DropWhile are complementary", func(t *testing.T) {
		s := Slice[int]{1, 2, 3, 4, 5, 6, 7, 8}
		predicate := func(x int) bool { return x < 5 }

		taken := s.TakeWhile(predicate)
		dropped := s.DropWhile(predicate)

		// Concatenate taken and dropped should equal original
		combined := append(taken, dropped...)

		if !combined.IsEq(s) {
			t.Errorf("TakeWhile and DropWhile should be complementary: expected %v, got %v", s, combined)
		}
	})

	t.Run("lengths should sum to original length", func(t *testing.T) {
		s := Slice[int]{10, 20, 30, 40, 50}
		predicate := func(x int) bool { return x <= 30 }

		taken := s.TakeWhile(predicate)
		dropped := s.DropWhile(predicate)

		if taken.Len()+dropped.Len() != s.Len() {
			t.Errorf("Lengths don't match: taken %d + dropped %d != original %d",
				taken.Len(), dropped.Len(), s.Len())
		}
	})
}

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

import (
	"testing"
)

func TestSliceCompact(t *testing.T) {
	t.Run("removes zero values from int slice", func(t *testing.T) {
		s := Slice[int]{1, 0, 2, 0, 3, 0}
		result := s.Compact()
		expected := Slice[int]{1, 2, 3}
		if !result.IsEq(expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("removes empty strings from string slice", func(t *testing.T) {
		s := Slice[string]{"hello", "", "world", "", "!"}
		result := s.Compact()
		expected := Slice[string]{"hello", "world", "!"}
		if !result.IsEq(expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("removes false from bool slice", func(t *testing.T) {
		s := Slice[bool]{true, false, true, false, true}
		result := s.Compact()
		expected := Slice[bool]{true, true, true}
		if !result.IsEq(expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("returns empty slice when all values are zero", func(t *testing.T) {
		s := Slice[int]{0, 0, 0, 0}
		result := s.Compact()
		if len(result) != 0 {
			t.Errorf("Expected empty slice, got %v", result)
		}
	})

	t.Run("returns same slice when no zero values", func(t *testing.T) {
		s := Slice[int]{1, 2, 3, 4, 5}
		result := s.Compact()
		if !result.IsEq(s) {
			t.Errorf("Expected %v, got %v", s, result)
		}
	})

	t.Run("returns empty slice when input is empty", func(t *testing.T) {
		s := Slice[int]{}
		result := s.Compact()
		if len(result) != 0 {
			t.Errorf("Expected empty slice, got %v", result)
		}
	})

	t.Run("works with float64", func(t *testing.T) {
		s := Slice[float64]{1.5, 0.0, 2.5, 0.0, 3.5}
		result := s.Compact()
		expected := Slice[float64]{1.5, 2.5, 3.5}
		if !result.IsEq(expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("removes zero runes from rune slice", func(t *testing.T) {
		s := Slice[rune]{'a', 0, 'b', 0, 'c'}
		result := s.Compact()
		expected := Slice[rune]{'a', 'b', 'c'}
		if !result.IsEq(expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("preserves order of non-zero elements", func(t *testing.T) {
		s := Slice[int]{5, 0, 4, 0, 3, 0, 2, 0, 1}
		result := s.Compact()
		expected := Slice[int]{5, 4, 3, 2, 1}
		if !result.IsEq(expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("single zero value returns empty slice", func(t *testing.T) {
		s := Slice[int]{0}
		result := s.Compact()
		if len(result) != 0 {
			t.Errorf("Expected empty slice, got %v", result)
		}
	})

	t.Run("single non-zero value returns same value", func(t *testing.T) {
		s := Slice[int]{42}
		result := s.Compact()
		expected := Slice[int]{42}
		if !result.IsEq(expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})
}


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


import (
	"reflect"
	"testing"
)

func TestSliceReduce_Basic(t *testing.T) {
	tests := []struct {
		name  string
		input Slice[int]
		block func(int) bool
		want  Slice[int]
	}{
		{
			name:  "keep even numbers",
			input: Slice[int]{1, 2, 3, 4, 5, 6},
			block: func(n int) bool { return n%2 == 0 },
			want:  Slice[int]{2, 4, 6},
		},
		{
			name:  "keep odd numbers",
			input: Slice[int]{1, 2, 3, 4, 5, 6},
			block: func(n int) bool { return n%2 != 0 },
			want:  Slice[int]{1, 3, 5},
		},
		{
			name:  "keep numbers greater than 3",
			input: Slice[int]{1, 2, 3, 4, 5, 6},
			block: func(n int) bool { return n > 3 },
			want:  Slice[int]{4, 5, 6},
		},
		{
			name:  "keep none",
			input: Slice[int]{1, 2, 3},
			block: func(n int) bool { return false },
			want:  Slice[int]{},
		},
		{
			name:  "keep all",
			input: Slice[int]{1, 2, 3},
			block: func(n int) bool { return true },
			want:  Slice[int]{1, 2, 3},
		},
		{
			name:  "empty slice",
			input: Slice[int]{},
			block: func(n int) bool { return true },
			want:  Slice[int]{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.input.Reduce(tt.block)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Reduce() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSliceReduce_String(t *testing.T) {
	tests := []struct {
		name  string
		input Slice[string]
		block func(string) bool
		want  Slice[string]
	}{
		{
			name:  "keep long strings",
			input: Slice[string]{"a", "hello", "hi", "world", "go"},
			block: func(s string) bool { return len(s) > 2 },
			want:  Slice[string]{"hello", "world"},
		},
		{
			name:  "keep strings starting with h",
			input: Slice[string]{"hello", "world", "hi", "golang"},
			block: func(s string) bool { return len(s) > 0 && s[0] == 'h' },
			want:  Slice[string]{"hello", "hi"},
		},
		{
			name:  "empty result",
			input: Slice[string]{"a", "b", "c"},
			block: func(s string) bool { return len(s) > 5 },
			want:  Slice[string]{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.input.Reduce(tt.block)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Reduce() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSliceReduce_Struct(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	people := Slice[Person]{
		{Name: "Alice", Age: 25},
		{Name: "Bob", Age: 30},
		{Name: "Charlie", Age: 20},
		{Name: "Diana", Age: 35},
	}

	tests := []struct {
		name  string
		input Slice[Person]
		block func(Person) bool
		want  Slice[Person]
	}{
		{
			name:  "keep adults over 25",
			input: people,
			block: func(p Person) bool { return p.Age > 25 },
			want: Slice[Person]{
				{Name: "Bob", Age: 30},
				{Name: "Diana", Age: 35},
			},
		},
		{
			name:  "keep names starting with A",
			input: people,
			block: func(p Person) bool { return len(p.Name) > 0 && p.Name[0] == 'A' },
			want: Slice[Person]{
				{Name: "Alice", Age: 25},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.input.Reduce(tt.block)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Reduce() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSliceReduce_ChainedOperations(t *testing.T) {
	// Test that Reduce works well in method chains
	input := Slice[int]{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// Chain: keep evens, then keep those > 4
	result := input.
		Reduce(func(n int) bool { return n%2 == 0 }).
		Reduce(func(n int) bool { return n > 4 })

	want := Slice[int]{6, 8, 10}

	if !reflect.DeepEqual(result, want) {
		t.Errorf("Chained Reduce() = %v, want %v", result, want)
	}
}

func TestSliceReduce_PreservesOrder(t *testing.T) {
	input := Slice[int]{5, 1, 4, 2, 3}

	// Keep odds - should preserve original order
	result := input.Reduce(func(n int) bool { return n%2 != 0 })

	want := Slice[int]{5, 1, 3}

	if !reflect.DeepEqual(result, want) {
		t.Errorf("Reduce() = %v, want %v (order should be preserved)", result, want)
	}
}

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

import (
	"testing"
)

// TestSliceSelectUntil_Found tests SelectUntil when condition is met
func TestSliceSelectUntil_Found(t *testing.T) {
	a := Slice[int]{1, 2, 3, 4, 5, 6}
	result := a.SelectUntil(func(e int) bool {
		return e == 3
	})
	expected := Slice[int]{1, 2}
	AssertSlicesEquals(t, expected, result)
}

// TestSliceSelectUntil_NotFound tests SelectUntil when condition is never met
func TestSliceSelectUntil_NotFound(t *testing.T) {
	a := Slice[int]{1, 2, 3, 4, 5}
	result := a.SelectUntil(func(e int) bool {
		return e == 10
	})
	// Should return the entire slice when condition is never met
	AssertSlicesEquals(t, a, result)
}

// TestSliceSelectUntil_FirstElement tests SelectUntil when first element matches
func TestSliceSelectUntil_FirstElement(t *testing.T) {
	a := Slice[int]{1, 2, 3, 4}
	result := a.SelectUntil(func(e int) bool {
		return e == 1
	})
	expected := Slice[int]{}
	AssertSlicesEquals(t, expected, result)
}

// TestSliceSelectUntil_LastElement tests SelectUntil when last element matches
func TestSliceSelectUntil_LastElement(t *testing.T) {
	a := Slice[int]{1, 2, 3, 4, 5}
	result := a.SelectUntil(func(e int) bool {
		return e == 5
	})
	expected := Slice[int]{1, 2, 3, 4}
	AssertSlicesEquals(t, expected, result)
}

// TestSliceSelectUntil_Empty tests SelectUntil on empty slice
func TestSliceSelectUntil_Empty(t *testing.T) {
	a := Slice[int]{}
	result := a.SelectUntil(func(e int) bool {
		return e == 1
	})
	expected := Slice[int]{}
	AssertSlicesEquals(t, expected, result)
}

// TestSliceSelectUntil_Strings tests SelectUntil with string type
func TestSliceSelectUntil_Strings(t *testing.T) {
	a := Slice[string]{"apple", "banana", "cherry", "date"}
	result := a.SelectUntil(func(e string) bool {
		return e == "cherry"
	})
	expected := Slice[string]{"apple", "banana"}
	AssertSlicesEquals(t, expected, result)
}

// TestSliceSelectUntil_ComplexCondition tests SelectUntil with complex predicate
func TestSliceSelectUntil_ComplexCondition(t *testing.T) {
	a := Slice[int]{1, 3, 5, 8, 10, 12}
	result := a.SelectUntil(func(e int) bool {
		return e%2 == 0 // First even number
	})
	expected := Slice[int]{1, 3, 5}
	AssertSlicesEquals(t, expected, result)
}

// TestSliceSelectUntil_Chaining tests SelectUntil with method chaining
func TestSliceSelectUntil_Chaining(t *testing.T) {
	a := Slice[int]{1, 2, 3, 4, 5, 6, 7, 8}
	result := a.SelectUntil(func(e int) bool {
		return e > 5
	}).Map(func(e int) int {
		return e * 2
	})
	expected := Slice[int]{2, 4, 6, 8, 10}
	AssertSlicesEquals(t, expected, result)
}


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

import (
	"testing"
)

func TestSliceTakeWhile(t *testing.T) {
	t.Run("takes elements while predicate is true", func(t *testing.T) {
		s := Slice[int]{1, 2, 3, 4, 5, 6}
		result := s.TakeWhile(func(x int) bool { return x < 4 })
		expected := Slice[int]{1, 2, 3}
		if !result.IsEq(expected) {
			t.Errorf("TakeWhile: expected %v, got %v", expected, result)
		}
	})

	t.Run("returns empty slice when first element fails predicate", func(t *testing.T) {
		s := Slice[int]{5, 6, 7, 8}
		result := s.TakeWhile(func(x int) bool { return x < 5 })
		expected := Slice[int]{}
		if !result.IsEq(expected) {
			t.Errorf("TakeWhile: expected %v, got %v", expected, result)
		}
	})

	t.Run("returns all elements when predicate always true", func(t *testing.T) {
		s := Slice[int]{1, 2, 3, 4}
		result := s.TakeWhile(func(x int) bool { return x < 10 })
		if !result.IsEq(s) {
			t.Errorf("TakeWhile: expected %v, got %v", s, result)
		}
	})

	t.Run("returns empty slice for empty input", func(t *testing.T) {
		s := Slice[int]{}
		result := s.TakeWhile(func(x int) bool { return x < 5 })
		expected := Slice[int]{}
		if !result.IsEq(expected) {
			t.Errorf("TakeWhile: expected %v, got %v", expected, result)
		}
	})

	t.Run("works with strings", func(t *testing.T) {
		s := Slice[string]{"apple", "banana", "cherry", "date"}
		result := s.TakeWhile(func(x string) bool { return len(x) < 6 })
		expected := Slice[string]{"apple"}
		if !result.IsEq(expected) {
			t.Errorf("TakeWhile: expected %v, got %v", expected, result)
		}
	})

	t.Run("stops at first false predicate", func(t *testing.T) {
		s := Slice[int]{2, 4, 6, 3, 8, 10}
		result := s.TakeWhile(func(x int) bool { return x%2 == 0 })
		expected := Slice[int]{2, 4, 6}
		if !result.IsEq(expected) {
			t.Errorf("TakeWhile: expected %v, got %v", expected, result)
		}
	})
}

func TestSliceDropWhile(t *testing.T) {
	t.Run("drops elements while predicate is true", func(t *testing.T) {
		s := Slice[int]{1, 2, 3, 4, 5, 6}
		result := s.DropWhile(func(x int) bool { return x < 4 })
		expected := Slice[int]{4, 5, 6}
		if !result.IsEq(expected) {
			t.Errorf("DropWhile: expected %v, got %v", expected, result)
		}
	})

	t.Run("returns original slice when first element fails predicate", func(t *testing.T) {
		s := Slice[int]{5, 6, 7, 8}
		result := s.DropWhile(func(x int) bool { return x < 5 })
		if !result.IsEq(s) {
			t.Errorf("DropWhile: expected %v, got %v", s, result)
		}
	})

	t.Run("returns empty slice when predicate always true", func(t *testing.T) {
		s := Slice[int]{1, 2, 3, 4}
		result := s.DropWhile(func(x int) bool { return x < 10 })
		expected := Slice[int]{}
		if !result.IsEq(expected) {
			t.Errorf("DropWhile: expected %v, got %v", expected, result)
		}
	})

	t.Run("returns empty slice for empty input", func(t *testing.T) {
		s := Slice[int]{}
		result := s.DropWhile(func(x int) bool { return x < 5 })
		expected := Slice[int]{}
		if !result.IsEq(expected) {
			t.Errorf("DropWhile: expected %v, got %v", expected, result)
		}
	})

	t.Run("works with strings", func(t *testing.T) {
		s := Slice[string]{"apple", "banana", "cherry", "date"}
		result := s.DropWhile(func(x string) bool { return len(x) < 6 })
		expected := Slice[string]{"banana", "cherry", "date"}
		if !result.IsEq(expected) {
			t.Errorf("DropWhile: expected %v, got %v", expected, result)
		}
	})

	t.Run("stops dropping at first false predicate", func(t *testing.T) {
		s := Slice[int]{2, 4, 6, 3, 8, 10}
		result := s.DropWhile(func(x int) bool { return x%2 == 0 })
		expected := Slice[int]{3, 8, 10}
		if !result.IsEq(expected) {
			t.Errorf("DropWhile: expected %v, got %v", expected, result)
		}
	})
}

func TestTakeWhileDropWhileComplement(t *testing.T) {
	t.Run("TakeWhile and DropWhile are complementary", func(t *testing.T) {
		s := Slice[int]{1, 2, 3, 4, 5, 6, 7, 8}
		predicate := func(x int) bool { return x < 5 }

		taken := s.TakeWhile(predicate)
		dropped := s.DropWhile(predicate)

		// Concatenate taken and dropped should equal original
		combined := append(taken, dropped...)

		if !combined.IsEq(s) {
			t.Errorf("TakeWhile and DropWhile should be complementary: expected %v, got %v", s, combined)
		}
	})

	t.Run("lengths should sum to original length", func(t *testing.T) {
		s := Slice[int]{10, 20, 30, 40, 50}
		predicate := func(x int) bool { return x <= 30 }

		taken := s.TakeWhile(predicate)
		dropped := s.DropWhile(predicate)

		if taken.Len()+dropped.Len() != s.Len() {
			t.Errorf("Lengths don't match: taken %d + dropped %d != original %d",
				taken.Len(), dropped.Len(), s.Len())
		}
	})
}

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
