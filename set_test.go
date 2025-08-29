package types

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestNewSet(t *testing.T) {
	tests := []struct {
		name  string
		slice []int
		want  []int
	}{
		{
			name:  "empty slice",
			slice: []int{},
			want:  []int{},
		},
		{
			name:  "single element",
			slice: []int{1},
			want:  []int{1},
		},
		{
			name:  "multiple unique elements",
			slice: []int{1, 2, 3},
			want:  []int{1, 2, 3},
		},
		{
			name:  "duplicate elements",
			slice: []int{1, 2, 2, 3, 1},
			want:  []int{1, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set := NewSet(tt.slice...)
			got := set.ToSlice()
			sort.Ints(got)
			sort.Ints(tt.want)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSet_Add(t *testing.T) {
	tests := []struct {
		name     string
		initial  []int
		add      int
		want     bool
		wantSize int
	}{
		{
			name:     "add to empty set",
			initial:  []int{},
			add:      1,
			want:     true,
			wantSize: 1,
		},
		{
			name:     "add unique element",
			initial:  []int{1, 2},
			add:      3,
			want:     true,
			wantSize: 3,
		},
		{
			name:     "add duplicate element",
			initial:  []int{1, 2, 3},
			add:      2,
			want:     false,
			wantSize: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set := NewSet(tt.initial...)
			got := set.Add(tt.add)

			if got != tt.want {
				t.Errorf("Add() = %v, want %v", got, tt.want)
			}
			if set.Size() != tt.wantSize {
				t.Errorf("Add() size = %v, want %v", set.Size(), tt.wantSize)
			}
		})
	}
}

func TestSet_Remove(t *testing.T) {
	tests := []struct {
		name     string
		initial  []int
		remove   int
		want     bool
		wantSize int
	}{
		{
			name:     "remove from empty set",
			initial:  []int{},
			remove:   1,
			want:     false,
			wantSize: 0,
		},
		{
			name:     "remove existing element",
			initial:  []int{1, 2, 3},
			remove:   2,
			want:     true,
			wantSize: 2,
		},
		{
			name:     "remove non-existing element",
			initial:  []int{1, 2, 3},
			remove:   4,
			want:     false,
			wantSize: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set := NewSet(tt.initial...)
			got := set.Remove(tt.remove)

			if got != tt.want {
				t.Errorf("Remove() = %v, want %v", got, tt.want)
			}
			if set.Size() != tt.wantSize {
				t.Errorf("Remove() size = %v, want %v", set.Size(), tt.wantSize)
			}
		})
	}
}

func TestSet_Contains(t *testing.T) {
	tests := []struct {
		name    string
		initial []int
		check   int
		want    bool
	}{
		{
			name:    "empty set",
			initial: []int{},
			check:   1,
			want:    false,
		},
		{
			name:    "element exists",
			initial: []int{1, 2, 3},
			check:   2,
			want:    true,
		},
		{
			name:    "element does not exist",
			initial: []int{1, 2, 3},
			check:   4,
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set := NewSet(tt.initial...)
			got := set.Contains(tt.check)

			if got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSet_Size(t *testing.T) {
	tests := []struct {
		name    string
		initial []int
		want    int
	}{
		{
			name:    "empty set",
			initial: []int{},
			want:    0,
		},
		{
			name:    "single element",
			initial: []int{1},
			want:    1,
		},
		{
			name:    "multiple elements",
			initial: []int{1, 2, 3, 4, 5},
			want:    5,
		},
		{
			name:    "duplicates removed",
			initial: []int{1, 1, 2, 2, 3},
			want:    3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set := NewSet(tt.initial...)
			got := set.Size()

			if got != tt.want {
				t.Errorf("Size() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSet_IsEmpty(t *testing.T) {
	tests := []struct {
		name    string
		initial []int
		want    bool
	}{
		{
			name:    "empty set",
			initial: []int{},
			want:    true,
		},
		{
			name:    "non-empty set",
			initial: []int{1, 2, 3},
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set := NewSet(tt.initial...)
			got := set.IsEmpty()

			if got != tt.want {
				t.Errorf("IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSet_Clear(t *testing.T) {
	tests := []struct {
		name    string
		initial []int
	}{
		{
			name:    "clear empty set",
			initial: []int{},
		},
		{
			name:    "clear non-empty set",
			initial: []int{1, 2, 3, 4, 5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set := NewSet(tt.initial...)
			set.Clear()

			if !set.IsEmpty() {
				t.Errorf("Clear() should result in empty set")
			}
			if set.Size() != 0 {
				t.Errorf("Clear() size = %v, want 0", set.Size())
			}
		})
	}
}

func TestSet_ToSlice(t *testing.T) {
	tests := []struct {
		name    string
		initial []int
		want    []int
	}{
		{
			name:    "empty set",
			initial: []int{},
			want:    []int{},
		},
		{
			name:    "single element",
			initial: []int{1},
			want:    []int{1},
		},
		{
			name:    "multiple elements",
			initial: []int{3, 1, 2},
			want:    []int{1, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set := NewSet(tt.initial...)
			got := set.ToSlice()
			sort.Ints(got)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSet_Clone(t *testing.T) {
	tests := []struct {
		name    string
		initial []int
	}{
		{
			name:    "clone empty set",
			initial: []int{},
		},
		{
			name:    "clone non-empty set",
			initial: []int{1, 2, 3, 4, 5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			original := NewSet(tt.initial...)
			clone := original.Clone()

			if !original.Equal(clone) {
				t.Errorf("Clone() should be equal to original")
			}

			// Test independence
			clone.Add(999)
			if original.Contains(999) {
				t.Errorf("Clone() should be independent of original")
			}
		})
	}
}

func TestSet_Union(t *testing.T) {
	tests := []struct {
		name string
		set1 []int
		set2 []int
		want []int
	}{
		{
			name: "union with empty sets",
			set1: []int{},
			set2: []int{},
			want: []int{},
		},
		{
			name: "union with one empty set",
			set1: []int{1, 2, 3},
			set2: []int{},
			want: []int{1, 2, 3},
		},
		{
			name: "union with no overlap",
			set1: []int{1, 2},
			set2: []int{3, 4},
			want: []int{1, 2, 3, 4},
		},
		{
			name: "union with overlap",
			set1: []int{1, 2, 3},
			set2: []int{2, 3, 4},
			want: []int{1, 2, 3, 4},
		},
		{
			name: "union with identical sets",
			set1: []int{1, 2, 3},
			set2: []int{1, 2, 3},
			want: []int{1, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s1 := NewSet(tt.set1...)
			s2 := NewSet(tt.set2...)
			result := s1.Union(s2)
			got := result.ToSlice()
			sort.Ints(got)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Union() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSet_Intersection(t *testing.T) {
	tests := []struct {
		name string
		set1 []int
		set2 []int
		want []int
	}{
		{
			name: "intersection with empty sets",
			set1: []int{},
			set2: []int{},
			want: []int{},
		},
		{
			name: "intersection with one empty set",
			set1: []int{1, 2, 3},
			set2: []int{},
			want: []int{},
		},
		{
			name: "intersection with no overlap",
			set1: []int{1, 2},
			set2: []int{3, 4},
			want: []int{},
		},
		{
			name: "intersection with overlap",
			set1: []int{1, 2, 3},
			set2: []int{2, 3, 4},
			want: []int{2, 3},
		},
		{
			name: "intersection with identical sets",
			set1: []int{1, 2, 3},
			set2: []int{1, 2, 3},
			want: []int{1, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s1 := NewSet(tt.set1...)
			s2 := NewSet(tt.set2...)
			result := s1.Intersection(s2)
			got := result.ToSlice()
			sort.Ints(got)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Intersection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSet_Difference(t *testing.T) {
	tests := []struct {
		name string
		set1 []int
		set2 []int
		want []int
	}{
		{
			name: "difference with empty sets",
			set1: []int{},
			set2: []int{},
			want: []int{},
		},
		{
			name: "difference with empty second set",
			set1: []int{1, 2, 3},
			set2: []int{},
			want: []int{1, 2, 3},
		},
		{
			name: "difference with empty first set",
			set1: []int{},
			set2: []int{1, 2, 3},
			want: []int{},
		},
		{
			name: "difference with no overlap",
			set1: []int{1, 2},
			set2: []int{3, 4},
			want: []int{1, 2},
		},
		{
			name: "difference with overlap",
			set1: []int{1, 2, 3},
			set2: []int{2, 3, 4},
			want: []int{1},
		},
		{
			name: "difference with identical sets",
			set1: []int{1, 2, 3},
			set2: []int{1, 2, 3},
			want: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s1 := NewSet(tt.set1...)
			s2 := NewSet(tt.set2...)
			result := s1.Difference(s2)
			got := result.ToSlice()
			sort.Ints(got)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Difference() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSet_SymmetricDifference(t *testing.T) {
	tests := []struct {
		name string
		set1 []int
		set2 []int
		want []int
	}{
		{
			name: "symmetric difference with empty sets",
			set1: []int{},
			set2: []int{},
			want: []int{},
		},
		{
			name: "symmetric difference with one empty set",
			set1: []int{1, 2, 3},
			set2: []int{},
			want: []int{1, 2, 3},
		},
		{
			name: "symmetric difference with no overlap",
			set1: []int{1, 2},
			set2: []int{3, 4},
			want: []int{1, 2, 3, 4},
		},
		{
			name: "symmetric difference with overlap",
			set1: []int{1, 2, 3},
			set2: []int{2, 3, 4},
			want: []int{1, 4},
		},
		{
			name: "symmetric difference with identical sets",
			set1: []int{1, 2, 3},
			set2: []int{1, 2, 3},
			want: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s1 := NewSet(tt.set1...)
			s2 := NewSet(tt.set2...)
			result := s1.SymmetricDifference(s2)
			got := result.ToSlice()
			sort.Ints(got)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SymmetricDifference() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSet_IsSubset(t *testing.T) {
	tests := []struct {
		name string
		set1 []int
		set2 []int
		want bool
	}{
		{
			name: "empty set is subset of empty set",
			set1: []int{},
			set2: []int{},
			want: true,
		},
		{
			name: "empty set is subset of non-empty set",
			set1: []int{},
			set2: []int{1, 2, 3},
			want: true,
		},
		{
			name: "non-empty set is not subset of empty set",
			set1: []int{1, 2},
			set2: []int{},
			want: false,
		},
		{
			name: "proper subset",
			set1: []int{1, 2},
			set2: []int{1, 2, 3},
			want: true,
		},
		{
			name: "not a subset",
			set1: []int{1, 2, 4},
			set2: []int{1, 2, 3},
			want: false,
		},
		{
			name: "identical sets are subsets",
			set1: []int{1, 2, 3},
			set2: []int{1, 2, 3},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s1 := NewSet(tt.set1...)
			s2 := NewSet(tt.set2...)
			got := s1.IsSubset(s2)

			if got != tt.want {
				t.Errorf("IsSubset() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSet_IsSuperset(t *testing.T) {
	tests := []struct {
		name string
		set1 []int
		set2 []int
		want bool
	}{
		{
			name: "empty set is superset of empty set",
			set1: []int{},
			set2: []int{},
			want: true,
		},
		{
			name: "non-empty set is superset of empty set",
			set1: []int{1, 2, 3},
			set2: []int{},
			want: true,
		},
		{
			name: "empty set is not superset of non-empty set",
			set1: []int{},
			set2: []int{1, 2},
			want: false,
		},
		{
			name: "proper superset",
			set1: []int{1, 2, 3},
			set2: []int{1, 2},
			want: true,
		},
		{
			name: "not a superset",
			set1: []int{1, 2, 3},
			set2: []int{1, 2, 4},
			want: false,
		},
		{
			name: "identical sets are supersets",
			set1: []int{1, 2, 3},
			set2: []int{1, 2, 3},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s1 := NewSet(tt.set1...)
			s2 := NewSet(tt.set2...)
			got := s1.IsSuperset(s2)

			if got != tt.want {
				t.Errorf("IsSuperset() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSet_IsDisjoint(t *testing.T) {
	tests := []struct {
		name string
		set1 []int
		set2 []int
		want bool
	}{
		{
			name: "empty sets are disjoint",
			set1: []int{},
			set2: []int{},
			want: true,
		},
		{
			name: "empty and non-empty sets are disjoint",
			set1: []int{},
			set2: []int{1, 2, 3},
			want: true,
		},
		{
			name: "disjoint sets",
			set1: []int{1, 2},
			set2: []int{3, 4},
			want: true,
		},
		{
			name: "overlapping sets are not disjoint",
			set1: []int{1, 2, 3},
			set2: []int{2, 3, 4},
			want: false,
		},
		{
			name: "identical sets are not disjoint",
			set1: []int{1, 2, 3},
			set2: []int{1, 2, 3},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s1 := NewSet(tt.set1...)
			s2 := NewSet(tt.set2...)
			got := s1.IsDisjoint(s2)

			if got != tt.want {
				t.Errorf("IsDisjoint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSet_Equal(t *testing.T) {
	tests := []struct {
		name string
		set1 []int
		set2 []int
		want bool
	}{
		{
			name: "empty sets are equal",
			set1: []int{},
			set2: []int{},
			want: true,
		},
		{
			name: "empty and non-empty sets are not equal",
			set1: []int{},
			set2: []int{1, 2, 3},
			want: false,
		},
		{
			name: "identical sets are equal",
			set1: []int{1, 2, 3},
			set2: []int{1, 2, 3},
			want: true,
		},
		{
			name: "different order same elements are equal",
			set1: []int{3, 1, 2},
			set2: []int{1, 2, 3},
			want: true,
		},
		{
			name: "different sets are not equal",
			set1: []int{1, 2, 3},
			set2: []int{1, 2, 4},
			want: false,
		},
		{
			name: "different size sets are not equal",
			set1: []int{1, 2},
			set2: []int{1, 2, 3},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s1 := NewSet(tt.set1...)
			s2 := NewSet(tt.set2...)
			got := s1.Equal(s2)

			if got != tt.want {
				t.Errorf("Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSet_Each(t *testing.T) {
	tests := []struct {
		name    string
		initial []int
		want    map[int]bool // track which elements were visited
	}{
		{
			name:    "empty set",
			initial: []int{},
			want:    map[int]bool{},
		},
		{
			name:    "single element",
			initial: []int{1},
			want:    map[int]bool{1: true},
		},
		{
			name:    "multiple elements",
			initial: []int{1, 2, 3},
			want:    map[int]bool{1: true, 2: true, 3: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set := NewSet(tt.initial...)
			visited := make(map[int]bool)

			set.Each(func(item int) {
				visited[item] = true
			})

			if !reflect.DeepEqual(visited, tt.want) {
				t.Errorf("Each() visited = %v, want %v", visited, tt.want)
			}
		})
	}
}

func TestSetMap(t *testing.T) {
	tests := []struct {
		name    string
		initial []int
		mapFn   func(int) string
		want    []string
	}{
		{
			name:    "empty set",
			initial: []int{},
			mapFn:   func(n int) string { return fmt.Sprintf("%d", n) },
			want:    []string{},
		},
		{
			name:    "int to string",
			initial: []int{1, 2, 3},
			mapFn:   func(n int) string { return fmt.Sprintf("num_%d", n) },
			want:    []string{"num_1", "num_2", "num_3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set := NewSet(tt.initial...)

			if tt.name == "square numbers" {
				// Special case for int to int mapping
				mapFn := func(n int) int { return n * n }
				result := SetMap(set, mapFn)
				got := result.ToSlice()
				sort.Ints(got)
				want := tt.want
				sort.Strings(want)

				if !reflect.DeepEqual(got, want) {
					t.Errorf("Map() = %v, want %v", got, want)
				}
			} else {
				// For string mapping tests
				result := SetMap(set, tt.mapFn)
				got := result.ToSlice()
				sort.Strings(got)
				want := tt.want
				sort.Strings(want)

				if !reflect.DeepEqual(got, want) {
					t.Errorf("Map() = %v, want %v", got, want)
				}
			}
		})
	}
}

func TestSet_Filter(t *testing.T) {
	tests := []struct {
		name      string
		initial   []int
		predicate func(int) bool
		want      []int
	}{
		{
			name:      "empty set",
			initial:   []int{},
			predicate: func(n int) bool { return n > 0 },
			want:      []int{},
		},
		{
			name:      "filter even numbers",
			initial:   []int{1, 2, 3, 4, 5, 6},
			predicate: func(n int) bool { return n%2 == 0 },
			want:      []int{2, 4, 6},
		},
		{
			name:      "filter greater than 3",
			initial:   []int{1, 2, 3, 4, 5},
			predicate: func(n int) bool { return n > 3 },
			want:      []int{4, 5},
		},
		{
			name:      "filter all elements (none match)",
			initial:   []int{1, 2, 3},
			predicate: func(n int) bool { return n > 10 },
			want:      []int{},
		},
		{
			name:      "filter no elements (all match)",
			initial:   []int{1, 2, 3},
			predicate: func(n int) bool { return n > 0 },
			want:      []int{1, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set := NewSet(tt.initial...)
			result := set.Filter(tt.predicate)
			got := result.ToSlice()
			sort.Ints(got)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Filter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSet_Reject(t *testing.T) {
	tests := []struct {
		name      string
		initial   []int
		predicate func(int) bool
		want      []int
	}{
		{
			name:      "empty set",
			initial:   []int{},
			predicate: func(n int) bool { return n > 0 },
			want:      []int{},
		},
		{
			name:      "reject even numbers",
			initial:   []int{1, 2, 3, 4, 5, 6},
			predicate: func(n int) bool { return n%2 == 0 },
			want:      []int{1, 3, 5},
		},
		{
			name:      "reject greater than 3",
			initial:   []int{1, 2, 3, 4, 5},
			predicate: func(n int) bool { return n > 3 },
			want:      []int{1, 2, 3},
		},
		{
			name:      "reject all elements (all match predicate)",
			initial:   []int{1, 2, 3},
			predicate: func(n int) bool { return n > 0 },
			want:      []int{},
		},
		{
			name:      "reject no elements (none match predicate)",
			initial:   []int{1, 2, 3},
			predicate: func(n int) bool { return n > 10 },
			want:      []int{1, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set := NewSet(tt.initial...)
			result := set.Reject(tt.predicate)
			got := result.ToSlice()
			sort.Ints(got)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Reject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSet_Find(t *testing.T) {
	tests := []struct {
		name      string
		initial   []int
		predicate func(int) bool
		wantValue int
		wantFound bool
	}{
		{
			name:      "empty set",
			initial:   []int{},
			predicate: func(n int) bool { return n > 0 },
			wantValue: 0,
			wantFound: false,
		},
		{
			name:      "find first even number",
			initial:   []int{1, 3, 5, 2, 4, 6},
			predicate: func(n int) bool { return n%2 == 0 },
			wantValue: 0, // Can be any even number since order is not guaranteed
			wantFound: true,
		},
		{
			name:      "find element greater than 10",
			initial:   []int{1, 2, 3, 15, 20},
			predicate: func(n int) bool { return n > 10 },
			wantValue: 0, // Can be 15 or 20
			wantFound: true,
		},
		{
			name:      "no element matches",
			initial:   []int{1, 2, 3},
			predicate: func(n int) bool { return n > 10 },
			wantValue: 0,
			wantFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set := NewSet(tt.initial...)
			gotValue, gotFound := set.Find(tt.predicate)

			if gotFound != tt.wantFound {
				t.Errorf("Find() found = %v, want %v", gotFound, tt.wantFound)
			}

			if tt.wantFound && gotFound {
				// Verify the found value satisfies the predicate
				if !tt.predicate(gotValue) {
					t.Errorf("Find() returned value %v that doesn't satisfy predicate", gotValue)
				}
				// Verify the found value was in the original set
				if !set.Contains(gotValue) {
					t.Errorf("Find() returned value %v that wasn't in the set", gotValue)
				}
			}
		})
	}
}

func TestSet_All(t *testing.T) {
	tests := []struct {
		name      string
		initial   []int
		predicate func(int) bool
		want      bool
	}{
		{
			name:      "empty set - should return true",
			initial:   []int{},
			predicate: func(n int) bool { return n > 0 },
			want:      true,
		},
		{
			name:      "all elements are positive",
			initial:   []int{1, 2, 3, 4, 5},
			predicate: func(n int) bool { return n > 0 },
			want:      true,
		},
		{
			name:      "not all elements are even",
			initial:   []int{2, 3, 4, 6},
			predicate: func(n int) bool { return n%2 == 0 },
			want:      false,
		},
		{
			name:      "all elements are greater than 10",
			initial:   []int{11, 12, 13},
			predicate: func(n int) bool { return n > 10 },
			want:      true,
		},
		{
			name:      "not all elements are greater than 5",
			initial:   []int{1, 6, 7, 8},
			predicate: func(n int) bool { return n > 5 },
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set := NewSet(tt.initial...)
			got := set.All(tt.predicate)

			if got != tt.want {
				t.Errorf("All() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSet_Any(t *testing.T) {
	tests := []struct {
		name      string
		initial   []int
		predicate func(int) bool
		want      bool
	}{
		{
			name:      "empty set - should return false",
			initial:   []int{},
			predicate: func(n int) bool { return n > 0 },
			want:      false,
		},
		{
			name:      "at least one element is even",
			initial:   []int{1, 3, 4, 5},
			predicate: func(n int) bool { return n%2 == 0 },
			want:      true,
		},
		{
			name:      "no element is greater than 10",
			initial:   []int{1, 2, 3, 4, 5},
			predicate: func(n int) bool { return n > 10 },
			want:      false,
		},
		{
			name:      "at least one element is negative",
			initial:   []int{-1, 2, 3},
			predicate: func(n int) bool { return n < 0 },
			want:      true,
		},
		{
			name:      "all elements satisfy predicate",
			initial:   []int{2, 4, 6, 8},
			predicate: func(n int) bool { return n%2 == 0 },
			want:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set := NewSet(tt.initial...)
			got := set.Any(tt.predicate)

			if got != tt.want {
				t.Errorf("Any() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSet_None(t *testing.T) {
	tests := []struct {
		name      string
		initial   []int
		predicate func(int) bool
		want      bool
	}{
		{
			name:      "empty set - should return true",
			initial:   []int{},
			predicate: func(n int) bool { return n > 0 },
			want:      true,
		},
		{
			name:      "no element is even",
			initial:   []int{1, 3, 5, 7},
			predicate: func(n int) bool { return n%2 == 0 },
			want:      true,
		},
		{
			name:      "some elements are greater than 5",
			initial:   []int{1, 6, 7, 8},
			predicate: func(n int) bool { return n > 5 },
			want:      false,
		},
		{
			name:      "no element is negative",
			initial:   []int{1, 2, 3},
			predicate: func(n int) bool { return n < 0 },
			want:      true,
		},
		{
			name:      "some elements are positive",
			initial:   []int{1, 2, 3},
			predicate: func(n int) bool { return n > 0 },
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set := NewSet(tt.initial...)
			got := set.None(tt.predicate)

			if got != tt.want {
				t.Errorf("None() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSet_Count(t *testing.T) {
	tests := []struct {
		name      string
		initial   []int
		predicate func(int) bool
		want      int
	}{
		{
			name:      "empty set",
			initial:   []int{},
			predicate: func(n int) bool { return n > 0 },
			want:      0,
		},
		{
			name:      "count even numbers",
			initial:   []int{1, 2, 3, 4, 5, 6},
			predicate: func(n int) bool { return n%2 == 0 },
			want:      3,
		},
		{
			name:      "count numbers greater than 3",
			initial:   []int{1, 2, 3, 4, 5},
			predicate: func(n int) bool { return n > 3 },
			want:      2,
		},
		{
			name:      "no elements match",
			initial:   []int{1, 2, 3},
			predicate: func(n int) bool { return n > 10 },
			want:      0,
		},
		{
			name:      "all elements match",
			initial:   []int{1, 2, 3},
			predicate: func(n int) bool { return n > 0 },
			want:      3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set := NewSet(tt.initial...)
			got := set.Count(tt.predicate)

			if got != tt.want {
				t.Errorf("Count() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReduce(t *testing.T) {
	tests := []struct {
		name        string
		initial     []int
		initial_val int
		reduceFn    func(int, int) int
		want        int
	}{
		{
			name:        "empty set",
			initial:     []int{},
			initial_val: 0,
			reduceFn:    func(acc, n int) int { return acc + n },
			want:        0,
		},
		{
			name:        "sum all elements",
			initial:     []int{1, 2, 3, 4, 5},
			initial_val: 0,
			reduceFn:    func(acc, n int) int { return acc + n },
			want:        15,
		},
		{
			name:        "multiply all elements",
			initial:     []int{2, 3, 4},
			initial_val: 1,
			reduceFn:    func(acc, n int) int { return acc * n },
			want:        24,
		},
		{
			name:        "find maximum",
			initial:     []int{3, 7, 2, 9, 1},
			initial_val: 0,
			reduceFn: func(acc, n int) int {
				if n > acc {
					return n
				}
				return acc
			},
			want: 9,
		},
		{
			name:        "concatenate as string lengths",
			initial:     []int{10, 100, 1000},
			initial_val: 0,
			reduceFn:    func(acc, n int) int { return acc + len(fmt.Sprintf("%d", n)) },
			want:        9, // 2 + 3 + 4 characters
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set := NewSet(tt.initial...)
			got := Reduce(set, tt.initial_val, tt.reduceFn)

			if got != tt.want {
				t.Errorf("Reduce() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSet_Partition(t *testing.T) {
	tests := []struct {
		name      string
		initial   []int
		predicate func(int) bool
		wantTrue  []int
		wantFalse []int
	}{
		{
			name:      "empty set",
			initial:   []int{},
			predicate: func(n int) bool { return n%2 == 0 },
			wantTrue:  []int{},
			wantFalse: []int{},
		},
		{
			name:      "partition even and odd",
			initial:   []int{1, 2, 3, 4, 5, 6},
			predicate: func(n int) bool { return n%2 == 0 },
			wantTrue:  []int{2, 4, 6},
			wantFalse: []int{1, 3, 5},
		},
		{
			name:      "partition greater than 3",
			initial:   []int{1, 2, 3, 4, 5},
			predicate: func(n int) bool { return n > 3 },
			wantTrue:  []int{4, 5},
			wantFalse: []int{1, 2, 3},
		},
		{
			name:      "all elements satisfy predicate",
			initial:   []int{2, 4, 6},
			predicate: func(n int) bool { return n%2 == 0 },
			wantTrue:  []int{2, 4, 6},
			wantFalse: []int{},
		},
		{
			name:      "no elements satisfy predicate",
			initial:   []int{1, 3, 5},
			predicate: func(n int) bool { return n%2 == 0 },
			wantTrue:  []int{},
			wantFalse: []int{1, 3, 5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set := NewSet(tt.initial...)
			gotTrue, gotFalse := set.Partition(tt.predicate)

			gotTrueSlice := gotTrue.ToSlice()
			gotFalseSlice := gotFalse.ToSlice()
			sort.Ints(gotTrueSlice)
			sort.Ints(gotFalseSlice)

			if !reflect.DeepEqual(gotTrueSlice, tt.wantTrue) {
				t.Errorf("Partition() true set = %v, want %v", gotTrueSlice, tt.wantTrue)
			}
			if !reflect.DeepEqual(gotFalseSlice, tt.wantFalse) {
				t.Errorf("Partition() false set = %v, want %v", gotFalseSlice, tt.wantFalse)
			}
		})
	}
}

func TestSet_Take(t *testing.T) {
	tests := []struct {
		name     string
		initial  []int
		n        int
		wantSize int
	}{
		{
			name:     "take from empty set",
			initial:  []int{},
			n:        3,
			wantSize: 0,
		},
		{
			name:     "take zero elements",
			initial:  []int{1, 2, 3, 4, 5},
			n:        0,
			wantSize: 0,
		},
		{
			name:     "take negative elements",
			initial:  []int{1, 2, 3, 4, 5},
			n:        -1,
			wantSize: 0,
		},
		{
			name:     "take fewer than available",
			initial:  []int{1, 2, 3, 4, 5},
			n:        3,
			wantSize: 3,
		},
		{
			name:     "take more than available",
			initial:  []int{1, 2, 3},
			n:        5,
			wantSize: 3,
		},
		{
			name:     "take exact amount",
			initial:  []int{1, 2, 3},
			n:        3,
			wantSize: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set := NewSet(tt.initial...)
			result := set.Take(tt.n)

			if result.Size() != tt.wantSize {
				t.Errorf("Take() size = %v, want %v", result.Size(), tt.wantSize)
			}

			// Verify all elements in result are from original set
			result.Each(func(item int) {
				if !set.Contains(item) {
					t.Errorf("Take() returned element %v not in original set", item)
				}
			})
		})
	}
}

func TestSet_Drop(t *testing.T) {
	tests := []struct {
		name     string
		initial  []int
		n        int
		wantSize int
	}{
		{
			name:     "drop from empty set",
			initial:  []int{},
			n:        3,
			wantSize: 0,
		},
		{
			name:     "drop zero elements",
			initial:  []int{1, 2, 3, 4, 5},
			n:        0,
			wantSize: 5,
		},
		{
			name:     "drop negative elements",
			initial:  []int{1, 2, 3, 4, 5},
			n:        -1,
			wantSize: 5,
		},
		{
			name:     "drop fewer than available",
			initial:  []int{1, 2, 3, 4, 5},
			n:        3,
			wantSize: 2,
		},
		{
			name:     "drop more than available",
			initial:  []int{1, 2, 3},
			n:        5,
			wantSize: 0,
		},
		{
			name:     "drop exact amount",
			initial:  []int{1, 2, 3},
			n:        3,
			wantSize: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set := NewSet(tt.initial...)
			result := set.Drop(tt.n)

			if result.Size() != tt.wantSize {
				t.Errorf("Drop() size = %v, want %v", result.Size(), tt.wantSize)
			}

			// Verify all elements in result are from original set
			result.Each(func(item int) {
				if !set.Contains(item) {
					t.Errorf("Drop() returned element %v not in original set", item)
				}
			})
		})
	}
}

func TestSet_String(t *testing.T) {
	tests := []struct {
		name    string
		initial []int
		want    string
	}{
		{
			name:    "empty set",
			initial: []int{},
			want:    "Set{}",
		},
		// Note: For non-empty sets, we can't test exact string matches
		// because the iteration order of maps is not guaranteed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set := NewSet(tt.initial...)
			got := set.String()

			if tt.name == "empty set" {
				if got != tt.want {
					t.Errorf("String() = %v, want %v", got, tt.want)
				}
			} else {
				// For non-empty sets, just verify format
				if !strings.HasPrefix(got, "Set{") || !strings.HasSuffix(got, "}") {
					t.Errorf("String() format incorrect: %v", got)
				}
			}
		})
	}
}

// Additional integration tests
func TestSet_Integration(t *testing.T) {
	t.Run("chaining operations", func(t *testing.T) {
		// Create a set, add elements, perform operations
		set := NewSet[int]()
		set.Add(1)
		set.Add(2)
		set.Add(3)
		set.Add(4)
		set.Add(5)

		// Filter even numbers, then check size
		evens := set.Filter(func(n int) bool { return n%2 == 0 })
		if evens.Size() != 2 {
			t.Errorf("Expected 2 even numbers, got %d", evens.Size())
		}

		// Union with another set
		other := NewSet(6, 7, 8)
		combined := set.Union(other)
		if combined.Size() != 8 {
			t.Errorf("Expected 8 elements after union, got %d", combined.Size())
		}

		// Check intersection
		intersection := evens.Intersection(other)
		if !intersection.IsEmpty() {
			t.Errorf("Expected empty intersection, got %v", intersection.ToSlice())
		}
	})

	t.Run("type safety with strings", func(t *testing.T) {
		stringSet := NewSet("apple", "banana", "cherry")

		// Test basic operations
		stringSet.Add("date")
		if !stringSet.Contains("apple") {
			t.Error("String set should contain 'apple'")
		}

		// Test filter
		longWords := stringSet.Filter(func(s string) bool { return len(s) > 5 })
		expected := []string{"banana", "cherry"}
		got := longWords.ToSlice()
		sort.Strings(got)

		if !reflect.DeepEqual(got, expected) {
			t.Errorf("Long words filter: got %v, want %v", got, expected)
		}
	})
}

func TestSetString(t *testing.T) {
	set := NewSet("apple", "banana", "cherry")
	str := set.String()

	if str != "Set{apple, banana, cherry}" {
		t.Errorf("String() format incorrect: %v", str)
	}
}
