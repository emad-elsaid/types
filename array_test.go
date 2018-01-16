package types

import "testing"

func TestArrayAt(t *testing.T) {
	a := Array{1, 2, 3, 4, 5, 6, 7}
	tcs := map[int64]interface{}{
		1:  2,
		6:  7,
		-1: 7,
		-7: 1,
		-8: nil,
		8:  nil,
	}

	for i, v := range tcs {
		result := a.At(i)
		if result != v {
			t.Errorf("With %d expected %d but found %d", i, v, result)
		}
	}
}

func TestArrayCount(t *testing.T) {
	a := Array{1, 2, 3, 4}
	result := a.Count()
	if result != 4 {
		t.Errorf("Expected %d but found %d", 4, result)
	}
}

func TestArrayCountElement(t *testing.T) {
	a := Array{1, 1, 1, 2, 2, 2, 2, 3, 3, 3, 3, 3}
	tcs := map[int]int64{
		1: 3,
		2: 4,
		3: 5,
		4: 0,
	}

	for i, v := range tcs {
		result := a.CountElement(i)
		if result != v {
			t.Errorf("With %d expected %d but found %d", i, v, result)
		}
	}
}

func TestArrayCountBy(t *testing.T) {
	a := Array{1, 1, 1, 2, 2, 2, 2, 3, 3, 3, 3, 3}
	tcs := []struct {
		name   string
		f      func(e interface{}) bool
		result int64
	}{
		{
			name:   "e == 1",
			f:      func(e interface{}) bool { return e == 1 },
			result: 3,
		},
		{
			name:   "e < 1",
			f:      func(e interface{}) bool { return e.(int) < 1 },
			result: 0,
		},
		{
			name:   "e > 1",
			f:      func(e interface{}) bool { return e.(int) > 1 },
			result: 9,
		},
	}

	for _, tc := range tcs {
		result := a.CountBy(tc.f)
		if result != tc.result {
			t.Errorf("With %s expected %d but found %d", tc.name, tc.result, result)
		}
	}
}

func TestArrayCycle(t *testing.T) {
	a := Array{1, 2}

	elements := []int{1, 2, 1, 2, 1, 2}
	a.Cycle(3, func(e interface{}) {
		result := e.(int)
		if result != elements[0] {
			t.Errorf("Expected %d but found %d", elements[0], result)
		}

		elements = elements[1:]
	})
}

func TestArrayAny(t *testing.T) {
	a := Array{false, false, false}
	identity := func(e interface{}) bool {
		return e.(bool)
	}
	if a.Any(identity) {
		t.Error("Expected false but got true")
	}

	a = Array{false, true, false}
	if !a.Any(identity) {
		t.Error("Expected true but got false")
	}
}

func TestArrayAll(t *testing.T) {
	a := Array{true, true, true}
	identity := func(e interface{}) bool {
		return e.(bool)
	}
	if !a.All(identity) {
		t.Error("Expected true but got false")
	}

	a = Array{false, true, false}
	if a.All(identity) {
		t.Error("Expected false but got true")
	}
}
