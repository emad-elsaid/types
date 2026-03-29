package types

import "testing"

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
