package types

import "testing"

func TestSliceAt(t *testing.T) {
	s := Slice[int]{1, 2, 3}

	// Test valid index
	val1 := s.At(1)
	if val1 == nil || *val1 != 2 {
		t.Errorf("Expected At(1) to return a pointer to 2, but got %v", val1)
	}

    // Test valid negative index
	val2 := s.At(-1)
	if val2 == nil || *val2 != 3 {
		t.Errorf("Expected At(-1) to return a pointer to 3, but got %v", val2)
	}

	// Test out of bounds (upper)
	if val := s.At(3); val != nil {
		t.Errorf("Expected At(3) to be nil, but got %v", val)
	}

	// Test out of bounds (negative)
	if val := s.At(-4); val != nil {
		t.Errorf("Expected At(-4) to be nil, but got %v", val)
	}

    // Test empty slice
    s_empty := Slice[int]{}
    if val := s_empty.At(0); val != nil {
        t.Errorf("Expected At(0) on empty slice to be nil, but got %v", val)
    }
}
