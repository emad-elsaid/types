package types

import "testing"

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
