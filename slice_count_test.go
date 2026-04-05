package types

import "testing"

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
