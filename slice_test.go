package types

import "testing"

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

func TestSliceAt(t *testing.T) {
	a := Slice[int]{1, 2, 3, 4, 5, 6, 7}
	tcs := map[int]int{
		1:  2,
		6:  7,
		-1: 7,
		-7: 1,
	}

	for i, v := range tcs {
		result := a.At(i)
		if *result != v {
			t.Errorf("With %d expected %d but found %d", i, v, result)
		}
	}

	result := a.At(8)
	if  result != nil {
		t.Errorf("With %d expected %v but found %d", 8, nil, result)
	}
}

func TestSliceLen(t *testing.T) {
	a := Slice[int]{1, 2, 3, 4}
	result := a.Len()
	if result != 4 {
		t.Errorf("Expected %d but found %d", 4, result)
	}
}

func TestSliceCountElement(t *testing.T) {
	a := Slice[int]{1, 1, 1, 2, 2, 2, 2, 3, 3, 3, 3, 3}
	tcs := map[int]int{
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

func TestSliceCountBy(t *testing.T) {
	a := Slice[int]{1, 1, 1, 2, 2, 2, 2, 3, 3, 3, 3, 3}
	tcs := []struct {
		name   string
		f      func(e int) bool
		result int
	}{
		{
			name:   "e == 1",
			f:      func(e int) bool { return e == 1 },
			result: 3,
		},
		{
			name:   "e < 1",
			f:      func(e int) bool { return e < 1 },
			result: 0,
		},
		{
			name:   "e > 1",
			f:      func(e int) bool { return e > 1 },
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

func TestSliceCycle(t *testing.T) {
	a := Slice[int]{1, 2}

	elements := []int{1, 2, 1, 2, 1, 2}
	a.Cycle(3, func(e int) {
		result := e
		if result != elements[0] {
			t.Errorf("Expected %d but found %d", elements[0], result)
		}

		elements = elements[1:]
	})
}

func TestSliceAny(t *testing.T) {
	a := Slice[bool]{false, false, false}
	identity := func(e bool) bool {
		return e
	}
	if a.Any(identity) {
		t.Error("Expected false but got true")
	}

	a = Slice[bool]{false, true, false}
	if !a.Any(identity) {
		t.Error("Expected true but got false")
	}
}

func TestSliceAll(t *testing.T) {
	a := Slice[bool]{true, true, true}
	identity := func(e bool) bool {
		return e
	}
	if !a.All(identity) {
		t.Error("Expected true but got false")
	}

	a = Slice[bool]{false, true, false}
	if a.All(identity) {
		t.Error("Expected false but got true")
	}
}

func TestSliceDelete(t *testing.T) {
	a := Slice[int]{1, 2, 3, 4, 1, 2, 3, 4}
	a = a.Delete(1)
	result := Slice[int]{2, 3, 4, 2, 3, 4}
	AssertSlicesEquals(t, result, a)
}

func TestSliceDeleteAt(t *testing.T) {
	a := Slice[int]{1, 2, 3, 4}
	a = a.DeleteAt(1)
	result := Slice[int]{1, 3, 4}
	AssertSlicesEquals(t, result, a)
}

func TestSliceDeleteIf(t *testing.T) {
	a := Slice[int]{1, 2, 3, 4, 5, 6}
	a = a.DeleteIf(func(e int) bool {
		return e > 1
	})
	result := Slice[int]{1}
	AssertSlicesEquals(t, result, a)
}

func TestSliceDrop(t *testing.T) {
	a := Slice[int]{1, 2, 3, 4, 5, 6, 7, 8, 9}
	a = a.Drop(5)
	result := Slice[int]{6, 7, 8, 9}
	AssertSlicesEquals(t, result, a)
}

func TestSliceEach(t *testing.T) {
	a := Slice[int]{1, 2, 3}
	sum := 0
	summer := func(e int) { sum += e }
	a.Each(summer)

	if sum != 6 {
		t.Errorf("Expected sum to be 6 but found %d", sum)
	}
}

func TestSliceEachIndex(t *testing.T) {
	a := Slice[int]{1, 2, 3}
	var sum int
	summer := func(i int) { sum += i }
	a.EachIndex(summer)

	if sum != 3 {
		t.Errorf("Expected sum to be 3 but found %d", sum)
	}
}

func TestSliceIsEmpty(t *testing.T) {
	a := Slice[int]{}
	if !a.IsEmpty() {
		t.Error("Expected to be empty but found not empty")
	}

	a = Slice[int]{1, 2, 3}
	if a.IsEmpty() {
		t.Error("Expected to be not empty but found empty")
	}
}

func TestSliceIsEq(t *testing.T) {
	a := Slice[int]{1, 2, 3, 4}
	b := Slice[int]{1, 2, 3, 4}

	if !a.IsEq(b) {
		t.Error("Expected arrat a to equal b but found otherwise")
	}
}

func TestSliceFetch(t *testing.T) {
	a := Slice[int]{1, 2}

	result := a.Fetch(0, -1)
	if result != 1 {
		t.Errorf("Expected 1 but got %d", result)
	}

	result = a.Fetch(-1,-1)
	if result != 2 {
		t.Errorf("Expected 2 but bot %d", result)
	}

	result = a.Fetch(3, -1)
	if result != -1 {
		t.Errorf("Expecte default value but got %d", result)
	}
}

func TestSliceFill(t *testing.T) {
	a := Slice[int]{1, 2, 3, 4, 5, 6}
	result := Slice[int]{1, 2, 1, 1, 1, 6}
	a.Fill(1, 2, 3)
	AssertSlicesEquals(t, result, a)
}

func TestSliceFillWith(t *testing.T) {
	a := Slice[int]{1, 2, 3, 4, 5, 6}
	result := Slice[int]{1, 2, 200, 300, 400, 6}
	a.FillWith(2, 3, func(i int) int {
		return i * 100
	})
	AssertSlicesEquals(t, result, a)
}

func TestSliceIndex(t *testing.T) {
	a := Slice[int]{1, 2, 3, 4, 5, 6}
	if a.Index(1) != 0 {
		t.Errorf("Expected 1 to have index of 0 but got %d", a.Index(1))
	}

	if a.Index(7) != -1 {
		t.Errorf("Expected 7 to have index of -1 gut git %d", a.Index(7))
	}
}

func TestSliceIndexBy(t *testing.T) {
	a := Slice[int]{1, 2, 3, 4, 5, 6}
	index := a.IndexBy(func(element int) bool {
		return element > 2
	})
	if index != 2 {
		t.Errorf("Expected element 3 index to be 2 got %d instead", index)
	}

	index = a.IndexBy(func(element int) bool {
		return element == -1
	})
	if index != -1 {
		t.Errorf("Expected element -1 index to be -1 go %d instead", index)
	}
}

func TestSliceFirst(t *testing.T) {
	a := Slice[int]{1, 2, 3, 4}
	if *a.First() != 1 {
		t.Errorf("Expected first element to be 1 got %d", a.First())
	}

	a = Slice[int]{}
	if a.First() != nil {
		t.Errorf("Expected first element to be nil got %d", a.First())
	}
}

func TestSliceLast(t *testing.T) {
	a := Slice[int]{1, 2, 3, 4}
	if *a.Last() != 4 {
		t.Errorf("Expected last element to be 4 got %d", a.Last())
	}

	a = Slice[int]{}
	if a.Last() != nil {
		t.Errorf("Expected last element to be nil got %d", a.Last())
	}
}

func TestSliceFirsts(t *testing.T) {
	a := Slice[int]{1, 2, 3, 4, 5, 6, 7, 8, 9}
	result := Slice[int]{1, 2, 3}
	AssertSlicesEquals(t, result, a.Firsts(3))
}

func TestSliceLasts(t *testing.T) {
	a := Slice[int]{1, 2, 3, 4, 5, 6, 7, 8, 9}
	result := Slice[int]{7, 8, 9}
	AssertSlicesEquals(t, result, a.Lasts(3))
}

func TestSliceInclude(t *testing.T) {
	a := Slice[int]{1, 2, 3, 4}
	if !a.Include(1) {
		t.Error("Expected 1 to be found but didn't find it")
	}

	if a.Include(-1) {
		t.Error("Expected the string not to be found but it was found!")
	}
}

func TestSliceInsert(t *testing.T) {
	a := Slice[int]{1, 2, 3, 4}
	result := Slice[int]{1, 2, 0, 3, 4}
	b := a.Insert(2, 0)
	AssertSlicesEquals(t, result, b)

	result = Slice[int]{1, 2, 3, 4, 0}
	c := a.Insert(4, 0)
	AssertSlicesEquals(t, result, c)

	result = Slice[int]{0, 1, 2, 3, 4}
	d := a.Insert(0, 0)
	AssertSlicesEquals(t, result, d)
}

func TestSliceKeepIf(t *testing.T) {
	a := Slice[int]{1, 2, 3, 4, 5, 6}
	a = a.KeepIf(func(e int) bool {
		return e > 3
	})
	result := Slice[int]{4, 5, 6}
	AssertSlicesEquals(t, result, a)
}

func TestSliceSelect(t *testing.T) {
	a := Slice[int]{1, 2, 3, 4, 5, 6}
	a = a.Select(func(e int) bool {
		return e > 3
	})
	result := Slice[int]{4, 5, 6}
	AssertSlicesEquals(t, result, a)
}

func TestSliceReduce(t *testing.T) {
	a := Slice[int]{1, 2, 3, 4, 5, 6}
	a = a.Reduce(func(e int) bool {
		return e > 3
	})
	result := Slice[int]{4, 5, 6}
	AssertSlicesEquals(t, result, a)
}

func TestSliceMap(t *testing.T) {
	a := Slice[int]{1, 2, 3, 4, 5}
	inc := func(e int) int {
		return e + 100
	}
	result := Slice[int]{101, 102, 103, 104, 105}
	AssertSlicesEquals(t, result, a.Map(inc))
}

func TestSliceMax(t *testing.T) {
	a := Slice[int]{1, 2, 3, 4}
	identity := func(e int) int {
		return e
	}

	result := a.Max(identity)
	if result != 4 {
		t.Errorf("Expected max to be 4 found %d", result)
	}

	a = Slice[int]{}
	result = a.Max(identity)
	if result != 0 {
		t.Errorf("Expected max of empty array to be nil got %d", result)
	}
}

func TestSliceMin(t *testing.T) {
	a := Slice[int]{4, 3, 2, 1}
	identity := func(e int) int {
		return e
	}

	result := a.Min(identity)
	if result != 1 {
		t.Errorf("Expected min to be 4 found %d", result)
	}

	a = Slice[int]{}
	result = a.Min(identity)
	if result != 0 {
		t.Errorf("Expected min of empty array to be %d", result)
	}
}

func TestSlicePush(t *testing.T) {
	a := Slice[int]{1, 2}
	a = a.Push(3)
	result := Slice[int]{1, 2, 3}
	AssertSlicesEquals(t, result, a)
}

func TestSlicePop(t *testing.T) {
	a := Slice[int]{1, 2, 3}
	a, e := a.Pop()
	result := Slice[int]{1, 2}
	if e != 3 {
		t.Errorf("Expected element to be 3 got %d", e)
	}
	AssertSlicesEquals(t, result, a)
}

func TestSliceUnshift(t *testing.T) {
	a := Slice[int]{1, 2, 3}
	a = a.Unshift(4)
	result := Slice[int]{4, 1, 2, 3}
	AssertSlicesEquals(t, result, a)
}

func TestSliceShift(t *testing.T) {
	a := Slice[int]{1, 2, 3}
	e, a := a.Shift()
	result := Slice[int]{2, 3}
	AssertSlicesEquals(t, result, a)
	if e != 1 {
		t.Errorf("Expected element to be 1 got %d", e)
	}

	a = Slice[int]{}
	e, a = a.Shift()
	if e != 0 {
		t.Errorf("Expected element to be nil got %d", e)
	}
}

func TestSliceReverse(t *testing.T) {
	a := Slice[int]{1, 2, 3}
	a = a.Reverse()
	result := Slice[int]{3, 2, 1}
	AssertSlicesEquals(t, result, a)
}

func TestSliceShuffle(t *testing.T) {
	a := Slice[int]{1, 2, 3, 4}
	a = a.Shuffle()
	notResult := Slice[int]{1, 2, 3, 4}
	if a.IsEq(notResult) {
		t.Error("Expected arrays not to equal after shuffle but it was the same")
	}
}
