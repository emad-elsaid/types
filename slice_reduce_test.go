package types

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
