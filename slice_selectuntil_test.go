package types

import (
	"testing"
)

// TestSliceSelectUntil_Found tests SelectUntil when condition is met
func TestSliceSelectUntil_Found(t *testing.T) {
	a := Slice[int]{1, 2, 3, 4, 5, 6}
	result := a.SelectUntil(func(e int) bool {
		return e == 3
	})
	expected := Slice[int]{1, 2}
	AssertSlicesEquals(t, expected, result)
}

// TestSliceSelectUntil_NotFound tests SelectUntil when condition is never met
func TestSliceSelectUntil_NotFound(t *testing.T) {
	a := Slice[int]{1, 2, 3, 4, 5}
	result := a.SelectUntil(func(e int) bool {
		return e == 10
	})
	// Should return the entire slice when condition is never met
	AssertSlicesEquals(t, a, result)
}

// TestSliceSelectUntil_FirstElement tests SelectUntil when first element matches
func TestSliceSelectUntil_FirstElement(t *testing.T) {
	a := Slice[int]{1, 2, 3, 4}
	result := a.SelectUntil(func(e int) bool {
		return e == 1
	})
	expected := Slice[int]{}
	AssertSlicesEquals(t, expected, result)
}

// TestSliceSelectUntil_LastElement tests SelectUntil when last element matches
func TestSliceSelectUntil_LastElement(t *testing.T) {
	a := Slice[int]{1, 2, 3, 4, 5}
	result := a.SelectUntil(func(e int) bool {
		return e == 5
	})
	expected := Slice[int]{1, 2, 3, 4}
	AssertSlicesEquals(t, expected, result)
}

// TestSliceSelectUntil_Empty tests SelectUntil on empty slice
func TestSliceSelectUntil_Empty(t *testing.T) {
	a := Slice[int]{}
	result := a.SelectUntil(func(e int) bool {
		return e == 1
	})
	expected := Slice[int]{}
	AssertSlicesEquals(t, expected, result)
}

// TestSliceSelectUntil_Strings tests SelectUntil with string type
func TestSliceSelectUntil_Strings(t *testing.T) {
	a := Slice[string]{"apple", "banana", "cherry", "date"}
	result := a.SelectUntil(func(e string) bool {
		return e == "cherry"
	})
	expected := Slice[string]{"apple", "banana"}
	AssertSlicesEquals(t, expected, result)
}

// TestSliceSelectUntil_ComplexCondition tests SelectUntil with complex predicate
func TestSliceSelectUntil_ComplexCondition(t *testing.T) {
	a := Slice[int]{1, 3, 5, 8, 10, 12}
	result := a.SelectUntil(func(e int) bool {
		return e%2 == 0 // First even number
	})
	expected := Slice[int]{1, 3, 5}
	AssertSlicesEquals(t, expected, result)
}

// TestSliceSelectUntil_Chaining tests SelectUntil with method chaining
func TestSliceSelectUntil_Chaining(t *testing.T) {
	a := Slice[int]{1, 2, 3, 4, 5, 6, 7, 8}
	result := a.SelectUntil(func(e int) bool {
		return e > 5
	}).Map(func(e int) int {
		return e * 2
	})
	expected := Slice[int]{2, 4, 6, 8, 10}
	AssertSlicesEquals(t, expected, result)
}
