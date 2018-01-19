package types

import "testing"

func AssertArraysEquals(t *testing.T, a1 Array, a2 Array) {
	if len(a1) != len(a2) {
		t.Errorf("Arrays doesn't have the same length %d, %d", len(a1), len(a2))
		return
	}

	for i := range a1 {
		if a1[i] != a2[i] {
			t.Errorf("Expected (%d) = %d but got %d", i, a1[i], a2[i])
		}
	}
}

func TestArrayAt(t *testing.T) {
	a := Array{1, 2, 3, 4, 5, 6, 7}
	tcs := map[int]Element{
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

func TestArrayLen(t *testing.T) {
	a := Array{1, 2, 3, 4}
	result := a.Len()
	if result != 4 {
		t.Errorf("Expected %d but found %d", 4, result)
	}
}

func TestArrayCountElement(t *testing.T) {
	a := Array{1, 1, 1, 2, 2, 2, 2, 3, 3, 3, 3, 3}
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

func TestArrayCountBy(t *testing.T) {
	a := Array{1, 1, 1, 2, 2, 2, 2, 3, 3, 3, 3, 3}
	tcs := []struct {
		name   string
		f      func(e Element) bool
		result int
	}{
		{
			name:   "e == 1",
			f:      func(e Element) bool { return e == 1 },
			result: 3,
		},
		{
			name:   "e < 1",
			f:      func(e Element) bool { return e.(int) < 1 },
			result: 0,
		},
		{
			name:   "e > 1",
			f:      func(e Element) bool { return e.(int) > 1 },
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
	a.Cycle(3, func(e Element) {
		result := e.(int)
		if result != elements[0] {
			t.Errorf("Expected %d but found %d", elements[0], result)
		}

		elements = elements[1:]
	})
}

func TestArrayAny(t *testing.T) {
	a := Array{false, false, false}
	identity := func(e Element) bool {
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
	identity := func(e Element) bool {
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

func TestArrayCompact(t *testing.T) {
	a := Array{1, nil, 3}
	a = a.Compact()

	result := Array{1, 3}
	AssertArraysEquals(t, result, a)
}

func TestArrayDelete(t *testing.T) {
	a := Array{1, 2, 3, 4, 1, 2, 3, 4}
	a = a.Delete(1)
	result := Array{2, 3, 4, 2, 3, 4}
	AssertArraysEquals(t, result, a)
}

func TestArrayDeleteAt(t *testing.T) {
	a := Array{1, 2, 3, 4}
	a = a.DeleteAt(1)
	result := Array{1, 3, 4}
	AssertArraysEquals(t, result, a)
}

func TestArrayDeleteIf(t *testing.T) {
	a := Array{1, 2, 3, 4, 5, 6}
	a = a.DeleteIf(func(e Element) bool {
		return e.(int) > 1
	})
	result := Array{1}
	AssertArraysEquals(t, result, a)
}

func TestArrayDrop(t *testing.T) {
	a := Array{1, 2, 3, 4, 5, 6, 7, 8, 9}
	a = a.Drop(5)
	result := Array{6, 7, 8, 9}
	AssertArraysEquals(t, result, a)
}

func TestArrayEach(t *testing.T) {
	a := Array{1, 2, 3}
	sum := 0
	summer := func(e Element) { sum += e.(int) }
	a.Each(summer)

	if sum != 6 {
		t.Errorf("Expected sum to be 6 but found %d", sum)
	}
}

func TestArrayEachIndex(t *testing.T) {
	a := Array{1, 2, 3}
	var sum int
	summer := func(i int) { sum += i }
	a.EachIndex(summer)

	if sum != 3 {
		t.Errorf("Expected sum to be 3 but found %d", sum)
	}
}

func TestArrayIsEmpty(t *testing.T) {
	a := Array{}
	if !a.IsEmpty() {
		t.Error("Expected to be empty but found not empty")
	}

	a = Array{1, 2, 3}
	if a.IsEmpty() {
		t.Error("Expected to be not empty but found empty")
	}
}

func TestArrayIsEq(t *testing.T) {
	a := Array{1, 2, 3, 4}
	b := Array{1, 2, 3, 4}

	if !a.IsEq(b) {
		t.Error("Expected arrat a to equal b but found otherwise")
	}

	b = Array{"a", "b", "c"}
	if a.IsEq(b) {
		t.Error("Expected array a to not equal b but found otherwise")
	}
}

func TestArrayFetch(t *testing.T) {
	a := Array{1, 2}

	result := a.Fetch(0, "default")
	if result != 1 {
		t.Errorf("Expected 1 but got %s", result)
	}

	result = a.Fetch(-1, "default")
	if result != 2 {
		t.Errorf("Expected 2 but bot %s", result)
	}

	result = a.Fetch(3, "default")
	if result != "default" {
		t.Errorf("Expecte default value but got %s", result)
	}
}

func TestArrayFill(t *testing.T) {
	a := Array{1, 2, 3, 4, 5, 6}
	result := Array{1, 2, 1, 1, 1, 6}
	a.Fill(1, 2, 3)
	AssertArraysEquals(t, result, a)
}

func TestArrayFillWith(t *testing.T) {
	a := Array{1, 2, 3, 4, 5, 6}
	result := Array{1, 2, 200, 300, 400, 6}
	a.FillWith(2, 3, func(i int) Element {
		return i * 100
	})
	AssertArraysEquals(t, result, a)
}

func TestArrayIndex(t *testing.T) {
	a := Array{1, 2, 3, 4, 5, 6}
	if a.Index(1) != 0 {
		t.Errorf("Expected 1 to have index of 0 but got %d", a.Index(1))
	}

	if a.Index(7) != -1 {
		t.Errorf("Expected 7 to have index of -1 gut git %d", a.Index(7))
	}
}

func TestArrayIndexBy(t *testing.T) {
	a := Array{1, 2, 3, 4, 5, 6}
	index := a.IndexBy(func(element Element) bool {
		return element.(int) > 2
	})
	if index != 2 {
		t.Errorf("Expected element 3 index to be 2 got %d instead", index)
	}

	index = a.IndexBy(func(element Element) bool {
		return element.(int) == -1
	})
	if index != -1 {
		t.Errorf("Expected element -1 index to be -1 go %d instead", index)
	}
}

func TestArrayFirst(t *testing.T) {
	a := Array{1, 2, 3, 4}
	if a.First().(int) != 1 {
		t.Errorf("Expected first element to be 1 got %d", a.First().(int))
	}

	a = Array{}
	if a.First() != nil {
		t.Errorf("Expected first element to be nil got %d", a.First())
	}
}

func TestArrayLast(t *testing.T) {
	a := Array{1, 2, 3, 4}
	if a.Last().(int) != 4 {
		t.Errorf("Expected last element to be 4 got %d", a.Last().(int))
	}

	a = Array{}
	if a.Last() != nil {
		t.Errorf("Expected last element to be nil got %d", a.Last())
	}
}

func TestArrayFirsts(t *testing.T) {
	a := Array{1, 2, 3, 4, 5, 6, 7, 8, 9}
	result := Array{1, 2, 3}
	AssertArraysEquals(t, result, a.Firsts(3))
}

func TestArrayLasts(t *testing.T) {
	a := Array{1, 2, 3, 4, 5, 6, 7, 8, 9}
	result := Array{7, 8, 9}
	AssertArraysEquals(t, result, a.Lasts(3))
}

func TestArrayFlatten(t *testing.T) {
	a := Array{1, 2, 3, Array{4, 5, 6, 7, 8, 9}}
	result := Array{1, 2, 3, 4, 5, 6, 7, 8, 9}
	AssertArraysEquals(t, result, a.Flatten())
}

func TestArrayInclude(t *testing.T) {
	a := Array{1, 2, 3, 4}
	if !a.Include(1) {
		t.Error("Expected 1 to be found but didn't find it")
	}

	if a.Include("not found") {
		t.Error("Expected the string not to be found but it was found!")
	}
}

func TestArrayInsert(t *testing.T) {
	a := Array{1, 2, 3, 4}
	result := Array{1, 2, "here", "is", "inserted", 3, 4}
	b := a.Insert(2, "here", "is", "inserted")
	AssertArraysEquals(t, result, b)

	result = Array{1, 2, 3, 4, "here", "is", "inserted"}
	c := a.Insert(4, "here", "is", "inserted")
	AssertArraysEquals(t, result, c)

	result = Array{"here", "is", "inserted", 1, 2, 3, 4}
	d := a.Insert(0, "here", "is", "inserted")
	AssertArraysEquals(t, result, d)
}

func TestArrayKeepIf(t *testing.T) {
	a := Array{1, 2, 3, 4, 5, 6}
	a = a.KeepIf(func(e Element) bool {
		return e.(int) > 3
	})
	result := Array{4, 5, 6}
	AssertArraysEquals(t, result, a)
}

func TestArrayMap(t *testing.T) {
	a := Array{1, 2, 3, 4, 5}
	inc := func(e Element) Element {
		return e.(int) + 100
	}
	result := Array{101, 102, 103, 104, 105}
	AssertArraysEquals(t, result, a.Map(inc))
}

func TestArrayMax(t *testing.T) {
	a := Array{1, 2, 3, 4}
	identity := func(e Element) int {
		return e.(int)
	}

	result := a.Max(identity)
	if result != 4 {
		t.Errorf("Expected max to be 4 found %d", result)
	}

	a = Array{}
	result = a.Max(identity)
	if result != nil {
		t.Errorf("Expected max of empty array to be nil got %d", result)
	}
}

func TestArrayMin(t *testing.T) {
	a := Array{4, 3, 2, 1}
	identity := func(e Element) int {
		return e.(int)
	}

	result := a.Min(identity)
	if result != 1 {
		t.Errorf("Expected min to be 4 found %d", result)
	}

	a = Array{}
	result = a.Min(identity)
	if result != nil {
		t.Errorf("Expected min of empty array to be %d", result)
	}
}

func TestArrayPush(t *testing.T) {
	a := Array{1, 2}
	a = a.Push(3)
	result := Array{1, 2, 3}
	AssertArraysEquals(t, result, a)
}

func TestArrayPop(t *testing.T) {
	a := Array{1, 2, 3}
	a, e := a.Pop()
	result := Array{1, 2}
	if e != 3 {
		t.Errorf("Expected element to be 3 got %d", e)
	}
	AssertArraysEquals(t, result, a)
}

func TestArrayUnshift(t *testing.T) {
	a := Array{1, 2, 3}
	a = a.Unshift(4)
	result := Array{4, 1, 2, 3}
	AssertArraysEquals(t, result, a)
}

func TestArrayShift(t *testing.T) {
	a := Array{1, 2, 3}
	e, a := a.Shift()
	result := Array{2, 3}
	AssertArraysEquals(t, result, a)
	if e != 1 {
		t.Errorf("Expected element to be 1 got %d", e)
	}

	a = Array{}
	e, a = a.Shift()
	if e != nil {
		t.Errorf("Expected element to be nil got %d", e)
	}
}

func TestArrayReverse(t *testing.T) {
	a := Array{1, 2, 3}
	a = a.Reverse()
	result := Array{3, 2, 1}
	AssertArraysEquals(t, result, a)
}

func TestArrayShuffle(t *testing.T) {
	a := Array{1, 2, 3, 4}
	a = a.Shuffle()
	notResult := Array{1, 2, 3, 4}
	if a.IsEq(notResult) {
		t.Error("Expected arrays not to equal after shuffle but it was the same")
	}
}
