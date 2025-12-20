# Bug Report: Types Library Code Review

This document outlines potential bugs and unintended behaviors found in the types library implementation.

## Critical Bugs (High Priority)

### 1. **Slice.DeleteAt() - Memory Safety Issue** ⚠️ CRITICAL
**Location:** `slice.go:93-95`

**Issue:** The `DeleteAt` method modifies the underlying array of the original slice, causing data corruption.

```go
func (a Slice[T]) DeleteAt(index int) Slice[T] {
    return append(a[:index], a[index+1:]...)
}
```

**Problem:** The returned slice shares the underlying array with the original slice. The append operation may overwrite elements in the original slice.

**Example of bug:**
```go
a := Slice[int]{1, 2, 3, 4, 5}
b := a.DeleteAt(2)  // Expected: b = [1, 2, 4, 5]
// But a is now corrupted: a = [1, 2, 4, 5, 5]
```

**Fix:** Create a new slice and copy elements:
```go
func (a Slice[T]) DeleteAt(index int) Slice[T] {
    result := make(Slice[T], 0, len(a)-1)
    result = append(result, a[:index]...)
    result = append(result, a[index+1:]...)
    return result
}
```

### 2. **Slice.Partition() - O(n²) Performance Bug** ⚠️ CRITICAL
**Location:** `slice.go:389-402`

**Issue:** Uses `Push()` method in a loop, creating a new slice on each iteration.

```go
for _, item := range s {
    if predicate(item) {
        trueSet = trueSet.Push(item)  // Creates new slice each time!
    } else {
        falseSet = falseSet.Push(item)
    }
}
```

**Problem:** Each `Push()` call returns a new slice, making this O(n²) instead of O(n).

**Fix:** Use append directly:
```go
for _, item := range s {
    if predicate(item) {
        trueSet = append(trueSet, item)
    } else {
        falseSet = append(falseSet, item)
    }
}
```

## High Priority Bugs

### 3. **Slice Methods - Missing Bounds Checking**
**Locations:** Multiple methods lack bounds checking and will panic:

- `Drop(count int)` - slice.go:110-112
  - Panics if `count > len(a)`

- `Firsts(count int)` - slice.go:220-222
  - Panics if `count > len(a)`

- `Lasts(count int)` - slice.go:226-228
  - Panics if `count > len(a)`

- `Pop()` - slice.go:333-335
  - Panics if slice is empty

- `Fill()` - slice.go:169-174
  - Panics if `start + length > len(a)`

- `FillWith()` - slice.go:179-184
  - Panics if `start + length > len(a)`

**Fix Example for Drop:**
```go
func (a Slice[T]) Drop(count int) Slice[T] {
    if count <= 0 {
        return a
    }
    if count >= len(a) {
        return Slice[T]{}
    }
    return a[count:]
}
```

### 4. **Slice.Reverse() and Slice.Shuffle() - Unexpected Mutation**
**Locations:** `slice.go:353-358` (Reverse), `slice.go:363-368` (Shuffle)

**Issue:** These methods modify the slice in-place, which is inconsistent with the library's general pattern of returning new slices.

```go
func (a Slice[T]) Reverse() Slice[T] {
    for i := len(a)/2 - 1; i >= 0; i-- {
        opp := len(a) - 1 - i
        a[i], a[opp] = a[opp], a[i]  // Mutates original!
    }
    return a
}
```

**Problem:** Users expect immutability based on other methods like `Map`, `Filter`, etc.

**Example of bug:**
```go
a := Slice[int]{1, 2, 3, 4, 5}
b := a.Reverse()
// Both a and b are now [5, 4, 3, 2, 1] - unexpected!
```

**Fix:** Create a copy first:
```go
func (a Slice[T]) Reverse() Slice[T] {
    result := make(Slice[T], len(a))
    copy(result, a)
    for i := len(result)/2 - 1; i >= 0; i-- {
        opp := len(result) - 1 - i
        result[i], result[opp] = result[opp], result[i]
    }
    return result
}
```

### 5. **Set.ToSlice() - Exposes Internal State** ⚠️
**Location:** `set.go:81-85`

**Issue:** Returns internal slice directly without copying.

```go
func (s *Set[T]) ToSlice() []T {
    return s.order  // Returns internal slice!
}
```

**Problem:** Caller can modify the returned slice and corrupt the set's internal state.

**Example of bug:**
```go
s := NewSet(1, 2, 3)
slice := s.ToSlice()
slice[0] = 999  // Corrupts the set's internal order!
```

**Fix:** Return a copy:
```go
func (s *Set[T]) ToSlice() []T {
    result := make([]T, len(s.order))
    copy(result, s.order)
    return result
}
```

## Medium Priority Issues

### 6. **Set - Nil Pointer Safety**
**Location:** Multiple methods in `set.go`

**Issue:** Methods don't check if the set is nil or if internal maps are uninitialized.

**Example:**
```go
var s Set[int]  // Not initialized with NewSet
s.Add(1)        // PANIC: assignment to entry in nil map
```

**Fix:** Either document that users must use `NewSet()`, or add nil checks to methods.

### 7. **Set.Each() - Incorrect Documentation**
**Location:** `set.go:156-162`

**Issue:** Comment says "The order of iteration is not guaranteed" but the code iterates in insertion order.

```go
// Each iterates over all elements in the set and calls the provided function for each element.
// The order of iteration is not guaranteed.  // <-- This is wrong!
func (s *Set[T]) Each(fn func(T)) {
    for _, item := range s.order {  // Actually iterates in order
        fn(item)
    }
}
```

**Fix:** Update the comment:
```go
// Each iterates over all elements in the set in insertion order and calls the provided function for each element.
```

### 8. **Set.SetMap() and Set.Filter() - Inefficient Pre-allocation**
**Locations:** `set.go:166-176` (SetMap), `set.go:178-191` (Filter)

**Issue:** Pre-allocates slices/maps then calls `Add()`, which is redundant.

```go
result := NewSet[U]()
result.order = make([]U, 0, len(s.order))  // Unnecessary
result.items = make(map[U]struct{}, len(s.order))  // Unnecessary

for _, item := range s.order {
    result.Add(fn(item))  // Add allocates again
}
```

**Fix:** Remove manual pre-allocation since `NewSet()` already initializes these.

### 9. **Slice.Fill() and Slice.FillWith() - Mutation Inconsistency**
**Locations:** `slice.go:169-174`, `slice.go:179-184`

**Issue:** Comments claim "will return same array object" but this contradicts the functional programming pattern of other methods.

**Recommendation:** Either:
1. Make these methods return a new slice (consistent with other methods)
2. Clearly document in README.md that these methods mutate in place

### 10. **Slice.Reduce() - Misleading Name**
**Location:** `slice.go:273-276`

**Issue:** Named `Reduce` but it's actually an alias for `KeepIf` (filter operation).

```go
// Reduce is an alias for KeepIf
func (a Slice[T]) Reduce(block func(T) bool) Slice[T] {
    return a.KeepIf(block)
}
```

**Problem:** In functional programming, "reduce" typically means folding/accumulating values (like `SliceReduce`), not filtering.

**Recommendation:** Consider deprecating this method or renaming it to avoid confusion.

## Summary

### Critical Issues:
1. ✅ **DeleteAt** corrupts original slice
2. ✅ **Partition** has O(n²) performance
3. ✅ **ToSlice** exposes internal state

### High Priority:
4. ✅ Missing bounds checking in 6+ methods
5. ✅ **Reverse** and **Shuffle** mutate unexpectedly

### Medium Priority:
6. ✅ Nil pointer safety
7. ✅ Documentation errors
8. ✅ Minor inefficiencies
9. ✅ Naming confusion

### Recommendations:
1. Add comprehensive bounds checking to all slice operations
2. Establish clear mutability contract (document which methods mutate)
3. Add safety checks for nil/uninitialized structs
4. Consider adding a test suite specifically for edge cases
5. Review all methods for consistency in mutability behavior
