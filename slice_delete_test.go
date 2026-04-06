package types

import "testing"

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
