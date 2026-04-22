package types

import (
	"reflect"
	"testing"
)

func AssertSlicesEquals[T comparable](t *testing.T, a1 Slice[T], a2 Slice[T]) {
	if len(a1) != len(a2) {
		t.Errorf("Slices doesn't have the same length %d, %d", len(a1), len(a2))
		return
	}

	for i := range a1 {
		if a1[i] != a2[i] {
			t.Errorf("Expected (%v) = %v but got %v", i, a1[i], a2[i])
		}
	}
}

func TestSliceReject(t *testing.T) {
	s := Slice[int]{1, 2, 3, 4, 5, 6}

	rejected := s.Reject(func(i int) bool {
		return i%2 == 0
	})

	expected := Slice[int]{1, 3, 5}
	if !rejected.IsEq(expected) {
		t.Errorf("Expected %v but got %v", expected, rejected)
	}
}

func TestSliceAll(t *testing.T) {
	s := Slice[int]{2, 4, 6}
	if !s.All(func(i int) bool { return i%2 == 0 }) {
		t.Error("all should be even")
	}

	s = Slice[int]{2, 4, 5}
	if s.All(func(i int) bool { return i%2 == 0 }) {
		t.Error("not all should be even")
	}
}

func TestSliceDuplicates(t *testing.T) {
	tests := []struct {
		name  string
		slice Slice[int]
		want  Slice[int]
	}{
		{"no duplicates", Slice[int]{1, 2, 3, 4}, Slice[int]{}},
		{"one duplicate pair", Slice[int]{1, 2, 2, 3}, Slice[int]{2}},
		{"multiple duplicate pairs", Slice[int]{1, 2, 2, 3, 3}, Slice[int]{2, 3}},
		{"multiple occurrences", Slice[int]{1, 1, 1, 2, 2}, Slice[int]{1, 2}},
		{"empty slice", Slice[int]{}, Slice[int]{}},
		{"all unique", Slice[int]{5, 6, 7}, Slice[int]{}},
		{"interspersed duplicates", Slice[int]{1, 2, 1, 3, 2}, Slice[int]{1, 2}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.slice.Duplicates(); !got.IsEq(tt.want) {
				t.Errorf("Slice.Duplicates() = %v, want %v", got, tt.want)
			}
		})
	}
}
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

func TestSliceCountElement(t *testing.T) {
	tests := []struct {
		name    string
		slice   Slice[int]
		element int
		want    int
	}{
		{
			name:    "empty slice",
			slice:   Slice[int]{},
			element: 1,
			want:    0,
		},
		{
			name:    "element not found",
			slice:   Slice[int]{1, 2, 3, 4, 5},
			element: 10,
			want:    0,
		},
		{
			name:    "single occurrence",
			slice:   Slice[int]{1, 2, 3, 4, 5},
			element: 3,
			want:    1,
		},
		{
			name:    "multiple occurrences",
			slice:   Slice[int]{1, 2, 3, 2, 4, 2, 5},
			element: 2,
			want:    3,
		},
		{
			name:    "all elements same",
			slice:   Slice[int]{7, 7, 7, 7},
			element: 7,
			want:    4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.slice.CountElement(tt.element)
			if got != tt.want {
				t.Errorf("CountElement() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSliceCountElementString(t *testing.T) {
	tests := []struct {
		name    string
		slice   Slice[string]
		element string
		want    int
	}{
		{
			name:    "count hello",
			slice:   Slice[string]{"hello", "world", "hello", "go"},
			element: "hello",
			want:    2,
		},
		{
			name:    "not found",
			slice:   Slice[string]{"hello", "world"},
			element: "missing",
			want:    0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.slice.CountElement(tt.element)
			if got != tt.want {
				t.Errorf("CountElement() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSliceCountBy(t *testing.T) {
	tests := []struct {
		name  string
		slice Slice[int]
		block func(int) bool
		want  int
	}{
		{
			name:  "empty slice",
			slice: Slice[int]{},
			block: func(i int) bool { return i%2 == 0 },
			want:  0,
		},
		{
			name:  "count even numbers",
			slice: Slice[int]{1, 2, 3, 4, 5, 6},
			block: func(i int) bool { return i%2 == 0 },
			want:  3,
		},
		{
			name:  "count odd numbers",
			slice: Slice[int]{1, 2, 3, 4, 5, 6},
			block: func(i int) bool { return i%2 != 0 },
			want:  3,
		},
		{
			name:  "count greater than threshold",
			slice: Slice[int]{1, 5, 10, 15, 20, 3},
			block: func(i int) bool { return i > 10 },
			want:  2,
		},
		{
			name:  "none match",
			slice: Slice[int]{1, 2, 3, 4, 5},
			block: func(i int) bool { return i > 100 },
			want:  0,
		},
		{
			name:  "all match",
			slice: Slice[int]{2, 4, 6, 8},
			block: func(i int) bool { return i%2 == 0 },
			want:  4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.slice.CountBy(tt.block)
			if got != tt.want {
				t.Errorf("CountBy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSliceCountByString(t *testing.T) {
	tests := []struct {
		name  string
		slice Slice[string]
		block func(string) bool
		want  int
	}{
		{
			name:  "count long strings",
			slice: Slice[string]{"a", "hello", "hi", "world", "go"},
			block: func(s string) bool { return len(s) > 2 },
			want:  2,
		},
		{
			name:  "count strings starting with h",
			slice: Slice[string]{"hello", "world", "hi", "go", "hey"},
			block: func(s string) bool { return len(s) > 0 && s[0] == 'h' },
			want:  3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.slice.CountBy(tt.block)
			if got != tt.want {
				t.Errorf("CountBy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSliceDelete(t *testing.T) {
	tests := []struct {
		name    string
		slice   Slice[int]
		element int
		want    Slice[int]
	}{
		{
			name:    "delete existing element",
			slice:   Slice[int]{1, 2, 3, 4, 5},
			element: 3,
			want:    Slice[int]{1, 2, 4, 5},
		},
		{
			name:    "delete first element",
			slice:   Slice[int]{1, 2, 3, 4, 5},
			element: 1,
			want:    Slice[int]{2, 3, 4, 5},
		},
		{
			name:    "delete last element",
			slice:   Slice[int]{1, 2, 3, 4, 5},
			element: 5,
			want:    Slice[int]{1, 2, 3, 4},
		},
		{
			name:    "delete non-existent element",
			slice:   Slice[int]{1, 2, 3, 4, 5},
			element: 10,
			want:    Slice[int]{1, 2, 3, 4, 5},
		},
		{
			name:    "delete from empty slice",
			slice:   Slice[int]{},
			element: 1,
			want:    Slice[int]{},
		},
		{
			name:    "delete duplicate elements",
			slice:   Slice[int]{1, 2, 2, 3, 2, 4},
			element: 2,
			want:    Slice[int]{1, 3, 4},
		},
		{
			name:    "delete all elements",
			slice:   Slice[int]{5, 5, 5, 5},
			element: 5,
			want:    Slice[int]{},
		},
		{
			name:    "delete from single element slice",
			slice:   Slice[int]{42},
			element: 42,
			want:    Slice[int]{},
		},
		{
			name:    "delete from single element slice - not found",
			slice:   Slice[int]{42},
			element: 43,
			want:    Slice[int]{42},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.slice.Delete(tt.element)
			if !got.IsEq(tt.want) {
				t.Errorf("Delete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSliceDeleteString(t *testing.T) {
	tests := []struct {
		name    string
		slice   Slice[string]
		element string
		want    Slice[string]
	}{
		{
			name:    "delete string element",
			slice:   Slice[string]{"apple", "banana", "cherry", "date"},
			element: "banana",
			want:    Slice[string]{"apple", "cherry", "date"},
		},
		{
			name:    "delete empty string",
			slice:   Slice[string]{"", "a", "b", ""},
			element: "",
			want:    Slice[string]{"a", "b"},
		},
		{
			name:    "delete case-sensitive",
			slice:   Slice[string]{"Apple", "apple", "APPLE"},
			element: "apple",
			want:    Slice[string]{"Apple", "APPLE"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.slice.Delete(tt.element)
			if !got.IsEq(tt.want) {
				t.Errorf("Delete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSliceDeleteAt(t *testing.T) {
	tests := []struct {
		name  string
		slice Slice[int]
		index int
		want  Slice[int]
	}{
		{
			name:  "delete at middle",
			slice: Slice[int]{1, 2, 3, 4, 5},
			index: 2,
			want:  Slice[int]{1, 2, 4, 5},
		},
		{
			name:  "delete at beginning",
			slice: Slice[int]{1, 2, 3, 4, 5},
			index: 0,
			want:  Slice[int]{2, 3, 4, 5},
		},
		{
			name:  "delete at end",
			slice: Slice[int]{1, 2, 3, 4, 5},
			index: 4,
			want:  Slice[int]{1, 2, 3, 4},
		},
		{
			name:  "negative index - out of bounds",
			slice: Slice[int]{1, 2, 3, 4, 5},
			index: -1,
			want:  Slice[int]{1, 2, 3, 4, 5},
		},
		{
			name:  "index too large - out of bounds",
			slice: Slice[int]{1, 2, 3, 4, 5},
			index: 10,
			want:  Slice[int]{1, 2, 3, 4, 5},
		},
		{
			name:  "index equals length - out of bounds",
			slice: Slice[int]{1, 2, 3},
			index: 3,
			want:  Slice[int]{1, 2, 3},
		},
		{
			name:  "delete from single element",
			slice: Slice[int]{42},
			index: 0,
			want:  Slice[int]{},
		},
		{
			name:  "delete from two elements - first",
			slice: Slice[int]{1, 2},
			index: 0,
			want:  Slice[int]{2},
		},
		{
			name:  "delete from two elements - second",
			slice: Slice[int]{1, 2},
			index: 1,
			want:  Slice[int]{1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.slice.DeleteAt(tt.index)
			if !got.IsEq(tt.want) {
				t.Errorf("DeleteAt(%d) = %v, want %v", tt.index, got, tt.want)
			}
		})
	}
}

func TestSliceDeleteIf(t *testing.T) {
	tests := []struct {
		name  string
		slice Slice[int]
		block func(int) bool
		want  Slice[int]
	}{
		{
			name:  "delete even numbers",
			slice: Slice[int]{1, 2, 3, 4, 5, 6},
			block: func(i int) bool { return i%2 == 0 },
			want:  Slice[int]{1, 3, 5},
		},
		{
			name:  "delete odd numbers",
			slice: Slice[int]{1, 2, 3, 4, 5, 6},
			block: func(i int) bool { return i%2 == 1 },
			want:  Slice[int]{2, 4, 6},
		},
		{
			name:  "delete greater than threshold",
			slice: Slice[int]{1, 5, 10, 15, 20},
			block: func(i int) bool { return i > 10 },
			want:  Slice[int]{1, 5, 10},
		},
		{
			name:  "delete all elements",
			slice: Slice[int]{1, 2, 3, 4, 5},
			block: func(i int) bool { return true },
			want:  Slice[int]{},
		},
		{
			name:  "delete no elements",
			slice: Slice[int]{1, 2, 3, 4, 5},
			block: func(i int) bool { return false },
			want:  Slice[int]{1, 2, 3, 4, 5},
		},
		{
			name:  "delete from empty slice",
			slice: Slice[int]{},
			block: func(i int) bool { return i > 0 },
			want:  Slice[int]{},
		},
		{
			name:  "delete negative numbers",
			slice: Slice[int]{-3, -2, -1, 0, 1, 2, 3},
			block: func(i int) bool { return i < 0 },
			want:  Slice[int]{0, 1, 2, 3},
		},
		{
			name:  "delete specific value",
			slice: Slice[int]{1, 2, 3, 2, 4, 2, 5},
			block: func(i int) bool { return i == 2 },
			want:  Slice[int]{1, 3, 4, 5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.slice.DeleteIf(tt.block)
			if !got.IsEq(tt.want) {
				t.Errorf("DeleteIf() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSliceDeleteIfString(t *testing.T) {
	tests := []struct {
		name  string
		slice Slice[string]
		block func(string) bool
		want  Slice[string]
	}{
		{
			name:  "delete empty strings",
			slice: Slice[string]{"a", "", "b", "", "c"},
			block: func(s string) bool { return s == "" },
			want:  Slice[string]{"a", "b", "c"},
		},
		{
			name:  "delete strings longer than 3",
			slice: Slice[string]{"cat", "dog", "elephant", "ant", "butterfly"},
			block: func(s string) bool { return len(s) > 3 },
			want:  Slice[string]{"cat", "dog", "ant"},
		},
		{
			name:  "delete strings starting with 'a'",
			slice: Slice[string]{"apple", "banana", "apricot", "cherry", "avocado"},
			block: func(s string) bool { return len(s) > 0 && s[0] == 'a' },
			want:  Slice[string]{"banana", "cherry"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.slice.DeleteIf(tt.block)
			if !got.IsEq(tt.want) {
				t.Errorf("DeleteIf() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkSliceDelete(b *testing.B) {
	slice := Slice[int]{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = slice.Delete(5)
	}
}

func BenchmarkSliceDeleteAt(b *testing.B) {
	slice := Slice[int]{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = slice.DeleteAt(4)
	}
}

func BenchmarkSliceDeleteIf(b *testing.B) {
	slice := Slice[int]{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = slice.DeleteIf(func(x int) bool { return x%2 == 0 })
	}
}


func TestSliceDrop(t *testing.T) {
	tests := []struct {
		name  string
		slice Slice[int]
		count int
		want  Slice[int]
	}{
		{
			name:  "drop 0 returns original slice",
			slice: Slice[int]{1, 2, 3, 4, 5},
			count: 0,
			want:  Slice[int]{1, 2, 3, 4, 5},
		},
		{
			name:  "drop negative returns original slice",
			slice: Slice[int]{1, 2, 3, 4, 5},
			count: -5,
			want:  Slice[int]{1, 2, 3, 4, 5},
		},
		{
			name:  "drop first element",
			slice: Slice[int]{1, 2, 3, 4, 5},
			count: 1,
			want:  Slice[int]{2, 3, 4, 5},
		},
		{
			name:  "drop first three elements",
			slice: Slice[int]{1, 2, 3, 4, 5},
			count: 3,
			want:  Slice[int]{4, 5},
		},
		{
			name:  "drop all elements",
			slice: Slice[int]{1, 2, 3, 4, 5},
			count: 5,
			want:  Slice[int]{},
		},
		{
			name:  "drop more than length returns empty slice",
			slice: Slice[int]{1, 2, 3},
			count: 10,
			want:  Slice[int]{},
		},
		{
			name:  "drop from empty slice returns empty slice",
			slice: Slice[int]{},
			count: 5,
			want:  Slice[int]{},
		},
		{
			name:  "drop from single element slice",
			slice: Slice[int]{42},
			count: 1,
			want:  Slice[int]{},
		},
		{
			name:  "drop 0 from single element",
			slice: Slice[int]{42},
			count: 0,
			want:  Slice[int]{42},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.slice.Drop(tt.count)
			if !got.IsEq(tt.want) {
				t.Errorf("Drop(%d) = %v, want %v", tt.count, got, tt.want)
			}
		})
	}
}

func TestSliceDropString(t *testing.T) {
	tests := []struct {
		name  string
		slice Slice[string]
		count int
		want  Slice[string]
	}{
		{
			name:  "drop strings from start",
			slice: Slice[string]{"a", "b", "c", "d", "e"},
			count: 2,
			want:  Slice[string]{"c", "d", "e"},
		},
		{
			name:  "drop all strings",
			slice: Slice[string]{"hello", "world"},
			count: 2,
			want:  Slice[string]{},
		},
		{
			name:  "drop more than length",
			slice: Slice[string]{"one"},
			count: 5,
			want:  Slice[string]{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.slice.Drop(tt.count)
			if !got.IsEq(tt.want) {
				t.Errorf("Drop(%d) = %v, want %v", tt.count, got, tt.want)
			}
		})
	}
}

func TestSliceDrop_DoesNotModifyOriginal(t *testing.T) {
	original := Slice[int]{1, 2, 3, 4, 5}
	expected := Slice[int]{1, 2, 3, 4, 5}

	_ = original.Drop(2)

	if !original.IsEq(expected) {
		t.Error("Drop should not modify the original slice")
	}
}

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

func TestSliceFill(t *testing.T) {
	t.Run("fills middle portion", func(t *testing.T) {
		s := Slice[int]{1, 2, 3, 4, 5}
		result := s.Fill(0, 1, 3)
		expected := Slice[int]{1, 0, 0, 0, 5}
		AssertSlicesEquals(t, expected, result)
	})

	t.Run("fills from start", func(t *testing.T) {
		s := Slice[int]{1, 2, 3, 4, 5}
		result := s.Fill(9, 0, 2)
		expected := Slice[int]{9, 9, 3, 4, 5}
		AssertSlicesEquals(t, expected, result)
	})

	t.Run("fills to end", func(t *testing.T) {
		s := Slice[int]{1, 2, 3, 4, 5}
		result := s.Fill(7, 3, 2)
		expected := Slice[int]{1, 2, 3, 7, 7}
		AssertSlicesEquals(t, expected, result)
	})

	t.Run("length exceeds bounds", func(t *testing.T) {
		s := Slice[int]{1, 2, 3, 4, 5}
		result := s.Fill(0, 2, 10)
		expected := Slice[int]{1, 2, 0, 0, 0}
		AssertSlicesEquals(t, expected, result)
	})

	t.Run("negative start returns unchanged", func(t *testing.T) {
		s := Slice[int]{1, 2, 3, 4, 5}
		result := s.Fill(0, -1, 3)
		expected := Slice[int]{1, 2, 3, 4, 5}
		AssertSlicesEquals(t, expected, result)
	})

	t.Run("start beyond bounds returns unchanged", func(t *testing.T) {
		s := Slice[int]{1, 2, 3, 4, 5}
		result := s.Fill(0, 10, 3)
		expected := Slice[int]{1, 2, 3, 4, 5}
		AssertSlicesEquals(t, expected, result)
	})

	t.Run("zero length returns unchanged", func(t *testing.T) {
		s := Slice[int]{1, 2, 3, 4, 5}
		result := s.Fill(0, 2, 0)
		expected := Slice[int]{1, 2, 3, 4, 5}
		AssertSlicesEquals(t, expected, result)
	})

	t.Run("negative length returns unchanged", func(t *testing.T) {
		s := Slice[int]{1, 2, 3, 4, 5}
		result := s.Fill(0, 2, -5)
		expected := Slice[int]{1, 2, 3, 4, 5}
		AssertSlicesEquals(t, expected, result)
	})

	t.Run("single element", func(t *testing.T) {
		s := Slice[int]{1, 2, 3, 4, 5}
		result := s.Fill(99, 2, 1)
		expected := Slice[int]{1, 2, 99, 4, 5}
		AssertSlicesEquals(t, expected, result)
	})

	t.Run("fills entire slice", func(t *testing.T) {
		s := Slice[int]{1, 2, 3, 4, 5}
		result := s.Fill(42, 0, 5)
		expected := Slice[int]{42, 42, 42, 42, 42}
		AssertSlicesEquals(t, expected, result)
	})

	t.Run("works with strings", func(t *testing.T) {
		s := Slice[string]{"a", "b", "c", "d"}
		result := s.Fill("x", 1, 2)
		expected := Slice[string]{"a", "x", "x", "d"}
		AssertSlicesEquals(t, expected, result)
	})

	t.Run("empty slice returns unchanged", func(t *testing.T) {
		s := Slice[int]{}
		result := s.Fill(0, 0, 5)
		expected := Slice[int]{}
		AssertSlicesEquals(t, expected, result)
	})
}

func TestSliceEach(t *testing.T) {
	t.Run("executes block for each element", func(t *testing.T) {
		s := Slice[int]{1, 2, 3, 4, 5}
		sum := 0
		s.Each(func(val int) {
			sum += val
		})
		if sum != 15 {
			t.Errorf("Expected sum to be 15, got %d", sum)
		}
	})

	t.Run("works with empty slice", func(t *testing.T) {
		s := Slice[int]{}
		count := 0
		s.Each(func(val int) {
			count++
		})
		if count != 0 {
			t.Errorf("Expected count to be 0, got %d", count)
		}
	})

	t.Run("receives elements in order", func(t *testing.T) {
		s := Slice[string]{"a", "b", "c"}
		result := []string{}
		s.Each(func(val string) {
			result = append(result, val)
		})
		expected := []string{"a", "b", "c"}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("allows mutation of external state", func(t *testing.T) {
		s := Slice[int]{10, 20, 30}
		collected := Slice[int]{}
		s.Each(func(val int) {
			collected = append(collected, val*2)
		})
		expected := Slice[int]{20, 40, 60}
		AssertSlicesEquals(t, expected, collected)
	})
}

func TestSliceEachIndex(t *testing.T) {
	t.Run("executes block for each index", func(t *testing.T) {
		s := Slice[int]{10, 20, 30, 40}
		indexSum := 0
		s.EachIndex(func(idx int) {
			indexSum += idx
		})
		if indexSum != 6 {
			t.Errorf("Expected index sum to be 6, got %d", indexSum)
		}
	})

	t.Run("works with empty slice", func(t *testing.T) {
		s := Slice[string]{}
		count := 0
		s.EachIndex(func(idx int) {
			count++
		})
		if count != 0 {
			t.Errorf("Expected count to be 0, got %d", count)
		}
	})

	t.Run("receives indices in order", func(t *testing.T) {
		s := Slice[bool]{true, false, true}
		indices := []int{}
		s.EachIndex(func(idx int) {
			indices = append(indices, idx)
		})
		expected := []int{0, 1, 2}
		if !reflect.DeepEqual(indices, expected) {
			t.Errorf("Expected %v, got %v", expected, indices)
		}
	})

	t.Run("can be used to build index map", func(t *testing.T) {
		s := Slice[string]{"a", "b", "c"}
		indexMap := make(map[int]bool)
		s.EachIndex(func(idx int) {
			indexMap[idx] = true
		})
		if len(indexMap) != 3 {
			t.Errorf("Expected map length 3, got %d", len(indexMap))
		}
		for i := 0; i < 3; i++ {
			if !indexMap[i] {
				t.Errorf("Expected index %d to be in map", i)
			}
		}
	})
}
func TestSliceFetch(t *testing.T) {
	tests := []struct {
		name         string
		slice        Slice[int]
		index        int
		defaultValue int
		want         int
	}{
		{
			name:         "fetch valid positive index",
			slice:        Slice[int]{10, 20, 30, 40},
			index:        2,
			defaultValue: 99,
			want:         30,
		},
		{
			name:         "fetch first element",
			slice:        Slice[int]{10, 20, 30},
			index:        0,
			defaultValue: 99,
			want:         10,
		},
		{
			name:         "fetch last element",
			slice:        Slice[int]{10, 20, 30},
			index:        2,
			defaultValue: 99,
			want:         30,
		},
		{
			name:         "fetch out of bounds positive returns default",
			slice:        Slice[int]{10, 20, 30},
			index:        10,
			defaultValue: 99,
			want:         99,
		},
		{
			name:         "fetch negative index from end",
			slice:        Slice[int]{10, 20, 30, 40},
			index:        -1,
			defaultValue: 99,
			want:         40,
		},
		{
			name:         "fetch negative index middle",
			slice:        Slice[int]{10, 20, 30, 40},
			index:        -2,
			defaultValue: 99,
			want:         30,
		},
		{
			name:         "fetch out of bounds negative returns default",
			slice:        Slice[int]{10, 20, 30},
			index:        -10,
			defaultValue: 99,
			want:         99,
		},
		{
			name:         "fetch from empty slice returns default",
			slice:        Slice[int]{},
			index:        0,
			defaultValue: 99,
			want:         99,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.slice.Fetch(tt.index, tt.defaultValue)
			if got != tt.want {
				t.Errorf("Slice.Fetch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSliceReverse(t *testing.T) {
	tests := []struct {
		name     string
		slice    Slice[int]
		expected Slice[int]
	}{
		{
			name:     "reverse normal slice",
			slice:    Slice[int]{1, 2, 3, 4, 5},
			expected: Slice[int]{5, 4, 3, 2, 1},
		},
		{
			name:     "reverse empty slice",
			slice:    Slice[int]{},
			expected: Slice[int]{},
		},
		{
			name:     "reverse single element",
			slice:    Slice[int]{42},
			expected: Slice[int]{42},
		},
		{
			name:     "reverse two elements",
			slice:    Slice[int]{1, 2},
			expected: Slice[int]{2, 1},
		},
		{
			name:     "reverse even length slice",
			slice:    Slice[int]{10, 20, 30, 40},
			expected: Slice[int]{40, 30, 20, 10},
		},
		{
			name:     "reverse odd length slice",
			slice:    Slice[int]{10, 20, 30, 40, 50},
			expected: Slice[int]{50, 40, 30, 20, 10},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a copy of original for mutation check
			original := make(Slice[int], len(tt.slice))
			copy(original, tt.slice)
			
			result := tt.slice.Reverse()
			
			// Verify result matches expected
			if !result.IsEq(tt.expected) {
				t.Errorf("Reverse() = %v, want %v", result, tt.expected)
			}
			
			// Verify original slice is unchanged
			if !tt.slice.IsEq(original) {
				t.Errorf("Original slice was modified: got %v, want %v", tt.slice, original)
			}
		})
	}
}

func TestSliceReverseStrings(t *testing.T) {
	input := Slice[string]{"apple", "banana", "cherry", "date"}
	expected := Slice[string]{"date", "cherry", "banana", "apple"}
	result := input.Reverse()
	
	if !result.IsEq(expected) {
		t.Errorf("Reverse() with strings = %v, want %v", result, expected)
	}
	
	// Verify original unchanged
	expectedOriginal := Slice[string]{"apple", "banana", "cherry", "date"}
	if !input.IsEq(expectedOriginal) {
		t.Errorf("Original slice was modified")
	}
}
func TestSliceFillWith(t *testing.T) {
	tests := []struct {
		name     string
		slice    Slice[int]
		start    int
		length   int
		block    func(int) int
		expected Slice[int]
	}{
		{
			name:     "fill middle of slice",
			slice:    Slice[int]{1, 2, 3, 4, 5},
			start:    1,
			length:   3,
			block:    func(i int) int { return i * 10 },
			expected: Slice[int]{1, 10, 20, 30, 5},
		},
		{
			name:     "fill from beginning",
			slice:    Slice[int]{1, 2, 3, 4, 5},
			start:    0,
			length:   2,
			block:    func(i int) int { return i + 100 },
			expected: Slice[int]{100, 101, 3, 4, 5},
		},
		{
			name:     "fill to end of slice",
			slice:    Slice[int]{1, 2, 3, 4, 5},
			start:    3,
			length:   2,
			block:    func(i int) int { return 0 },
			expected: Slice[int]{1, 2, 3, 0, 0},
		},
		{
			name:     "length exceeds slice bounds - should truncate",
			slice:    Slice[int]{1, 2, 3, 4, 5},
			start:    3,
			length:   10,
			block:    func(i int) int { return 99 },
			expected: Slice[int]{1, 2, 3, 99, 99},
		},
		{
			name:     "negative start - no modification",
			slice:    Slice[int]{1, 2, 3, 4, 5},
			start:    -1,
			length:   3,
			block:    func(i int) int { return 0 },
			expected: Slice[int]{1, 2, 3, 4, 5},
		},
		{
			name:     "start beyond slice length - no modification",
			slice:    Slice[int]{1, 2, 3, 4, 5},
			start:    10,
			length:   3,
			block:    func(i int) int { return 0 },
			expected: Slice[int]{1, 2, 3, 4, 5},
		},
		{
			name:     "zero length - no modification",
			slice:    Slice[int]{1, 2, 3, 4, 5},
			start:    2,
			length:   0,
			block:    func(i int) int { return 0 },
			expected: Slice[int]{1, 2, 3, 4, 5},
		},
		{
			name:     "negative length - no modification",
			slice:    Slice[int]{1, 2, 3, 4, 5},
			start:    2,
			length:   -5,
			block:    func(i int) int { return 0 },
			expected: Slice[int]{1, 2, 3, 4, 5},
		},
		{
			name:     "single element fill",
			slice:    Slice[int]{1, 2, 3, 4, 5},
			start:    2,
			length:   1,
			block:    func(i int) int { return i * i },
			expected: Slice[int]{1, 2, 4, 4, 5},
		},
		{
			name:     "fill entire slice",
			slice:    Slice[int]{1, 2, 3, 4, 5},
			start:    0,
			length:   5,
			block:    func(i int) int { return i },
			expected: Slice[int]{0, 1, 2, 3, 4},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.slice.FillWith(tt.start, tt.length, tt.block)
			AssertSlicesEquals(t, result, tt.expected)
		})
	}
}

func TestSliceFillWithStrings(t *testing.T) {
	slice := Slice[string]{"a", "b", "c", "d", "e"}
	result := slice.FillWith(1, 3, func(i int) string {
		return "x"
	})
	expected := Slice[string]{"a", "x", "x", "x", "e"}
	AssertSlicesEquals(t, result, expected)
}

func TestSliceFillWithEmptySlice(t *testing.T) {
	slice := Slice[int]{}
	result := slice.FillWith(0, 5, func(i int) int { return i })
	expected := Slice[int]{}
	AssertSlicesEquals(t, result, expected)
}

func TestSliceInsert(t *testing.T) {
	tests := []struct {
		name     string
		slice    Slice[int]
		index    int
		elements []int
		want     Slice[int]
	}{
		{
			name:     "insert at beginning",
			slice:    Slice[int]{3, 4, 5},
			index:    0,
			elements: []int{1, 2},
			want:     Slice[int]{1, 2, 3, 4, 5},
		},
		{
			name:     "insert at end",
			slice:    Slice[int]{1, 2, 3},
			index:    3,
			elements: []int{4, 5},
			want:     Slice[int]{1, 2, 3, 4, 5},
		},
		{
			name:     "insert in middle",
			slice:    Slice[int]{1, 2, 5, 6},
			index:    2,
			elements: []int{3, 4},
			want:     Slice[int]{1, 2, 3, 4, 5, 6},
		},
		{
			name:     "insert single element at beginning",
			slice:    Slice[int]{2, 3, 4},
			index:    0,
			elements: []int{1},
			want:     Slice[int]{1, 2, 3, 4},
		},
		{
			name:     "insert single element in middle",
			slice:    Slice[int]{1, 3, 4},
			index:    1,
			elements: []int{2},
			want:     Slice[int]{1, 2, 3, 4},
		},
		{
			name:     "insert single element at end",
			slice:    Slice[int]{1, 2, 3},
			index:    3,
			elements: []int{4},
			want:     Slice[int]{1, 2, 3, 4},
		},
		{
			name:     "insert into empty slice",
			slice:    Slice[int]{},
			index:    0,
			elements: []int{1, 2, 3},
			want:     Slice[int]{1, 2, 3},
		},
		{
			name:     "insert empty elements",
			slice:    Slice[int]{1, 2, 3},
			index:    1,
			elements: []int{},
			want:     Slice[int]{1, 2, 3},
		},
		{
			name:     "insert multiple elements at beginning",
			slice:    Slice[int]{4, 5},
			index:    0,
			elements: []int{1, 2, 3},
			want:     Slice[int]{1, 2, 3, 4, 5},
		},
		{
			name:     "insert multiple elements at end",
			slice:    Slice[int]{1, 2},
			index:    2,
			elements: []int{3, 4, 5},
			want:     Slice[int]{1, 2, 3, 4, 5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.slice.Insert(tt.index, tt.elements...)
			AssertSlicesEquals(t, got, tt.want)
		})
	}
}

func TestSliceInsertWithStrings(t *testing.T) {
	tests := []struct {
		name     string
		slice    Slice[string]
		index    int
		elements []string
		want     Slice[string]
	}{
		{
			name:     "insert string at beginning",
			slice:    Slice[string]{"world"},
			index:    0,
			elements: []string{"hello"},
			want:     Slice[string]{"hello", "world"},
		},
		{
			name:     "insert multiple strings in middle",
			slice:    Slice[string]{"a", "d"},
			index:    1,
			elements: []string{"b", "c"},
			want:     Slice[string]{"a", "b", "c", "d"},
		},
		{
			name:     "insert empty string",
			slice:    Slice[string]{"a", "b"},
			index:    1,
			elements: []string{""},
			want:     Slice[string]{"a", "", "b"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.slice.Insert(tt.index, tt.elements...)
			AssertSlicesEquals(t, got, tt.want)
		})
	}
}

func TestSliceInsertImmutability(t *testing.T) {
	original := Slice[int]{1, 2, 3}
	originalCopy := make(Slice[int], len(original))
	copy(originalCopy, original)

	result := original.Insert(1, 99)

	AssertSlicesEquals(t, original, originalCopy)
	AssertSlicesEquals(t, result, Slice[int]{1, 99, 2, 3})
}

func TestSlicePartition(t *testing.T) {
	tests := []struct {
		name          string
		slice         Slice[int]
		predicate     func(int) bool
		expectedTrue  Slice[int]
		expectedFalse Slice[int]
	}{
		{
			name:          "partition even and odd numbers",
			slice:         Slice[int]{1, 2, 3, 4, 5, 6},
			predicate:     func(n int) bool { return n%2 == 0 },
			expectedTrue:  Slice[int]{2, 4, 6},
			expectedFalse: Slice[int]{1, 3, 5},
		},
		{
			name:          "partition by threshold",
			slice:         Slice[int]{1, 5, 10, 15, 20},
			predicate:     func(n int) bool { return n > 10 },
			expectedTrue:  Slice[int]{15, 20},
			expectedFalse: Slice[int]{1, 5, 10},
		},
		{
			name:          "all match predicate",
			slice:         Slice[int]{2, 4, 6, 8},
			predicate:     func(n int) bool { return n%2 == 0 },
			expectedTrue:  Slice[int]{2, 4, 6, 8},
			expectedFalse: Slice[int]{},
		},
		{
			name:          "none match predicate",
			slice:         Slice[int]{1, 3, 5, 7},
			predicate:     func(n int) bool { return n%2 == 0 },
			expectedTrue:  Slice[int]{},
			expectedFalse: Slice[int]{1, 3, 5, 7},
		},
		{
			name:          "empty slice",
			slice:         Slice[int]{},
			predicate:     func(n int) bool { return n%2 == 0 },
			expectedTrue:  Slice[int]{},
			expectedFalse: Slice[int]{},
		},
		{
			name:          "single element matching",
			slice:         Slice[int]{2},
			predicate:     func(n int) bool { return n%2 == 0 },
			expectedTrue:  Slice[int]{2},
			expectedFalse: Slice[int]{},
		},
		{
			name:          "single element not matching",
			slice:         Slice[int]{1},
			predicate:     func(n int) bool { return n%2 == 0 },
			expectedTrue:  Slice[int]{},
			expectedFalse: Slice[int]{1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			trueSet, falseSet := tt.slice.Partition(tt.predicate)

			if !trueSet.IsEq(tt.expectedTrue) {
				t.Errorf("Partition() true set = %v, want %v", trueSet, tt.expectedTrue)
			}
			if !falseSet.IsEq(tt.expectedFalse) {
				t.Errorf("Partition() false set = %v, want %v", falseSet, tt.expectedFalse)
			}
		})
	}
}

func TestSlicePartitionStrings(t *testing.T) {
	words := Slice[string]{"apple", "banana", "cherry", "apricot", "blueberry"}
	predicate := func(s string) bool { return s[0] == 'a' }

	startsWithA, other := words.Partition(predicate)

	expectedA := Slice[string]{"apple", "apricot"}
	expectedOther := Slice[string]{"banana", "cherry", "blueberry"}

	if !startsWithA.IsEq(expectedA) {
		t.Errorf("Partition() true set = %v, want %v", startsWithA, expectedA)
	}
	if !other.IsEq(expectedOther) {
		t.Errorf("Partition() false set = %v, want %v", other, expectedOther)
	}
}

func TestSlicePartitionPreservesOrder(t *testing.T) {
	slice := Slice[int]{5, 1, 8, 3, 9, 2, 7, 4, 6}
	predicate := func(n int) bool { return n > 5 }

	greater, lessOrEqual := slice.Partition(predicate)

	expectedGreater := Slice[int]{8, 9, 7, 6}
	expectedLessOrEqual := Slice[int]{5, 1, 3, 2, 4}

	if !greater.IsEq(expectedGreater) {
		t.Errorf("Partition() should preserve order: got %v, want %v", greater, expectedGreater)
	}
	if !lessOrEqual.IsEq(expectedLessOrEqual) {
		t.Errorf("Partition() should preserve order: got %v, want %v", lessOrEqual, expectedLessOrEqual)
	}
}

func TestSliceMax(t *testing.T) {
	t.Run("returns element with max score", func(t *testing.T) {
		s := Slice[string]{"apple", "banana", "kiwi", "strawberry"}
		result := s.Max(func(s string) int {
			return len(s)
		})
		if result != "strawberry" {
			t.Errorf("Expected 'strawberry' but got '%s'", result)
		}
	})

	t.Run("returns first element for slice of size 1", func(t *testing.T) {
		s := Slice[int]{42}
		result := s.Max(func(i int) int { return i })
		if result != 42 {
			t.Errorf("Expected 42 but got %d", result)
		}
	})

	t.Run("returns zero value for empty slice", func(t *testing.T) {
		s := Slice[int]{}
		result := s.Max(func(i int) int { return i })
		if result != 0 {
			t.Errorf("Expected 0 but got %d", result)
		}
	})

	t.Run("handles negative scores", func(t *testing.T) {
		s := Slice[int]{-10, -5, -20, -3}
		result := s.Max(func(i int) int { return i })
		if result != -3 {
			t.Errorf("Expected -3 but got %d", result)
		}
	})

	t.Run("handles tie - returns first max", func(t *testing.T) {
		s := Slice[int]{1, 3, 3, 2}
		result := s.Max(func(i int) int { return i })
		if result != 3 {
			t.Errorf("Expected first 3 but got %d", result)
		}
	})
}

func TestSliceMin(t *testing.T) {
	t.Run("returns element with min score", func(t *testing.T) {
		s := Slice[string]{"apple", "banana", "kiwi", "strawberry"}
		result := s.Min(func(s string) int {
			return len(s)
		})
		if result != "kiwi" {
			t.Errorf("Expected 'kiwi' but got '%s'", result)
		}
	})

	t.Run("returns first element for slice of size 1", func(t *testing.T) {
		s := Slice[int]{42}
		result := s.Min(func(i int) int { return i })
		if result != 42 {
			t.Errorf("Expected 42 but got %d", result)
		}
	})

	t.Run("returns zero value for empty slice", func(t *testing.T) {
		s := Slice[int]{}
		result := s.Min(func(i int) int { return i })
		if result != 0 {
			t.Errorf("Expected 0 but got %d", result)
		}
	})

	t.Run("handles negative scores", func(t *testing.T) {
		s := Slice[int]{-10, -5, -20, -3}
		result := s.Min(func(i int) int { return i })
		if result != -20 {
			t.Errorf("Expected -20 but got %d", result)
		}
	})

	t.Run("handles tie - returns first min", func(t *testing.T) {
		s := Slice[int]{3, 1, 1, 2}
		result := s.Min(func(i int) int { return i })
		if result != 1 {
			t.Errorf("Expected first 1 but got %d", result)
		}
	})
}

func TestSliceReduce(t *testing.T) {
	t.Run("sums integers", func(t *testing.T) {
		s := Slice[int]{1, 2, 3, 4, 5}
		result := SliceReduce(s, 0, func(acc int, val int) int {
			return acc + val
		})
		if result != 15 {
			t.Errorf("Expected 15 but got %d", result)
		}
	})

	t.Run("concatenates strings", func(t *testing.T) {
		s := Slice[string]{"Hello", " ", "World"}
		result := SliceReduce(s, "", func(acc string, val string) string {
			return acc + val
		})
		if result != "Hello World" {
			t.Errorf("Expected 'Hello World' but got '%s'", result)
		}
	})

	t.Run("returns initial value for empty slice", func(t *testing.T) {
		s := Slice[int]{}
		result := SliceReduce(s, 100, func(acc int, val int) int {
			return acc + val
		})
		if result != 100 {
			t.Errorf("Expected 100 but got %d", result)
		}
	})

	t.Run("reduces to different type", func(t *testing.T) {
		s := Slice[int]{1, 2, 3, 4}
		result := SliceReduce(s, []int{}, func(acc []int, val int) []int {
			if val%2 == 0 {
				return append(acc, val)
			}
			return acc
		})
		expected := []int{2, 4}
		if len(result) != len(expected) {
			t.Errorf("Expected %v but got %v", expected, result)
		}
		for i := range expected {
			if result[i] != expected[i] {
				t.Errorf("Expected %v but got %v", expected, result)
			}
		}
	})

	t.Run("multiplies with initial value", func(t *testing.T) {
		s := Slice[int]{2, 3, 4}
		result := SliceReduce(s, 1, func(acc int, val int) int {
			return acc * val
		})
		if result != 24 {
			t.Errorf("Expected 24 but got %d", result)
		}
	})
}

func TestSliceSelect(t *testing.T) {
	t.Run("alias for KeepIf - filters even numbers", func(t *testing.T) {
		s := Slice[int]{1, 2, 3, 4, 5, 6}
		result := s.Select(func(i int) bool {
			return i%2 == 0
		})
		expected := Slice[int]{2, 4, 6}
		if !result.IsEq(expected) {
			t.Errorf("Expected %v but got %v", expected, result)
		}
	})

	t.Run("returns empty slice when nothing matches", func(t *testing.T) {
		s := Slice[int]{1, 3, 5}
		result := s.Select(func(i int) bool {
			return i%2 == 0
		})
		if len(result) != 0 {
			t.Errorf("Expected empty slice but got %v", result)
		}
	})
}

func TestIsEq(t *testing.T) {
	tests := []struct {
		name     string
		slice    Slice[int]
		other    Slice[int]
		expected bool
	}{
		{
			name:     "equal slices",
			slice:    Slice[int]{1, 2, 3},
			other:    Slice[int]{1, 2, 3},
			expected: true,
		},
		{
			name:     "different lengths",
			slice:    Slice[int]{1, 2, 3},
			other:    Slice[int]{1, 2},
			expected: false,
		},
		{
			name:     "same length different values",
			slice:    Slice[int]{1, 2, 3},
			other:    Slice[int]{1, 2, 4},
			expected: false,
		},
		{
			name:     "both empty",
			slice:    Slice[int]{},
			other:    Slice[int]{},
			expected: true,
		},
		{
			name:     "one empty one not",
			slice:    Slice[int]{1},
			other:    Slice[int]{},
			expected: false,
		},
		{
			name:     "different values at start",
			slice:    Slice[int]{5, 2, 3},
			other:    Slice[int]{1, 2, 3},
			expected: false,
		},
		{
			name:     "different values at middle",
			slice:    Slice[int]{1, 5, 3},
			other:    Slice[int]{1, 2, 3},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.slice.IsEq(tt.other)
			if result != tt.expected {
				t.Errorf("IsEq() = %v, want %v", result, tt.expected)
			}
		})
	}
}
