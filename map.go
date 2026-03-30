package types

import "sync"

// Map is a wrapper for sync.Map with same methods (except for CompareAnd*)
// Check: https://pkg.go.dev/sync#Map
type Map[K comparable, V any] struct {
	store sync.Map
}

// NewMapFrom creates a new Map initialized with entries from a regular Go map.
// This provides a convenient way to convert from map[K]V to *Map[K, V].
func NewMapFrom[K comparable, V any](m map[K]V) *Map[K, V] {
	result := &Map[K, V]{}
	for k, v := range m {
		result.Store(k, v)
	}
	return result
}

func (m *Map[K, V]) Store(k K, v V) {
	m.store.Store(k, v)
}

func (m *Map[K, V]) Load(k K) (v V, ok bool) {
	var a any

	a, ok = m.store.Load(k)
	if ok {
		v = a.(V)
	}

	return
}

// Delete see sync.Map#Delete
func (m *Map[K, V]) Delete(k K) {
	m.store.Delete(k)
}

// Range see sync.Map#Range
func (m *Map[K, V]) Range(f func(k K, v V) bool) {
	m.store.Range(func(k, v any) bool {
		return f(k.(K), v.(V))
	})
}

// LoadAndDelete see sync.Map#LoadAndDelete
func (m *Map[K, V]) LoadAndDelete(k K) (v V, loaded bool) {
	var a any
	a, loaded = m.store.LoadAndDelete(k)
	if loaded {
		v = a.(V)
	}

	return
}

// LoadOrStore see sync.Map#LoadOrStore
func (m *Map[K, V]) LoadOrStore(k K, v V) (actual V, loaded bool) {
	var a any
	a, loaded = m.store.LoadOrStore(k, v)
	actual = a.(V)

	return
}

// Swap see sync.Map#Swap
func (m *Map[K, V]) Swap(k K, v V) (previous V, loaded bool) {
	var a any
	a, loaded = m.store.Swap(k, v)
	previous, _ = a.(V)

	return
}

// Keys returns a slice containing all keys in the map.
// The order of keys is not guaranteed.
func (m *Map[K, V]) Keys() []K {
	keys := []K{}
	m.store.Range(func(k, v any) bool {
		keys = append(keys, k.(K))
		return true
	})
	return keys
}

// Values returns a slice containing all values in the map.
// The order of values is not guaranteed.
func (m *Map[K, V]) Values() []V {
	values := []V{}
	m.store.Range(func(k, v any) bool {
		values = append(values, v.(V))
		return true
	})
	return values
}

// Size returns the number of entries in the map.
// Note: This operation requires iterating through all entries.
func (m *Map[K, V]) Size() int {
	size := 0
	m.store.Range(func(k, v any) bool {
		size++
		return true
	})
	return size
}

// Clear removes all entries from the map.
func (m *Map[K, V]) Clear() {
	m.store.Range(func(k, v any) bool {
		m.store.Delete(k)
		return true
	})
}

// Has checks if a key exists in the map.
// This is a convenience method equivalent to checking the ok return value of Load.
func (m *Map[K, V]) Has(k K) bool {
	_, ok := m.store.Load(k)
	return ok
}

// ForEach iterates over all entries in the map, calling the provided function
// for each key-value pair. Unlike Range, this does not support early termination.
func (m *Map[K, V]) ForEach(f func(k K, v V)) {
	m.store.Range(func(k, v any) bool {
		f(k.(K), v.(V))
		return true
	})
}

// Filter returns a new Map containing only the entries that satisfy the predicate.
func (m *Map[K, V]) Filter(predicate func(k K, v V) bool) *Map[K, V] {
	result := &Map[K, V]{}
	m.store.Range(func(k, v any) bool {
		key := k.(K)
		value := v.(V)
		if predicate(key, value) {
			result.Store(key, value)
		}
		return true
	})
	return result
}

// Transform returns a new Map with transformed values while preserving keys.
// This is useful for mapping operations on map values while keeping the same key structure.
func (m *Map[K, V]) Transform(transformer func(k K, v V) V) *Map[K, V] {
	result := &Map[K, V]{}
	m.store.Range(func(k, v any) bool {
		key := k.(K)
		value := v.(V)
		result.Store(key, transformer(key, value))
		return true
	})
	return result
}

// Entries returns a regular Go map containing all key-value pairs from the Map.
// This provides a convenient way to convert from *Map[K, V] to map[K]V.
// The order of entries is not guaranteed.
func (m *Map[K, V]) Entries() map[K]V {
	result := make(map[K]V)
	m.store.Range(func(k, v any) bool {
		result[k.(K)] = v.(V)
		return true
	})
	return result
}

// Any returns true if the predicate returns true for any entry in the map.
// Returns false for empty maps.
func (m *Map[K, V]) Any(predicate func(k K, v V) bool) bool {
	found := false
	m.store.Range(func(k, v any) bool {
		if predicate(k.(K), v.(V)) {
			found = true
			return false // stop iteration
		}
		return true
	})
	return found
}

// All returns true if the predicate returns true for all entries in the map.
// Returns true for empty maps.
func (m *Map[K, V]) All(predicate func(k K, v V) bool) bool {
	allMatch := true
	m.store.Range(func(k, v any) bool {
		if !predicate(k.(K), v.(V)) {
			allMatch = false
			return false // stop iteration
		}
		return true
	})
	return allMatch
}

// None returns true if no entries in the map satisfy the predicate.
// Returns true for empty maps.
func (m *Map[K, V]) None(predicate func(k K, v V) bool) bool {
	return !m.Any(predicate)
}

// Partition divides the map into two new maps based on the predicate function.
// Returns two maps: the first contains entries that satisfy the predicate,
// the second contains entries that do not satisfy the predicate.
// This is inspired by Ruby's Hash#partition method.
func (m *Map[K, V]) Partition(predicate func(k K, v V) bool) (*Map[K, V], *Map[K, V]) {
	trueMap := &Map[K, V]{}
	falseMap := &Map[K, V]{}

	m.store.Range(func(k, v any) bool {
		key := k.(K)
		value := v.(V)
		if predicate(key, value) {
			trueMap.Store(key, value)
		} else {
			falseMap.Store(key, value)
		}
		return true
	})

	return trueMap, falseMap
}

// Merge returns a new Map combining entries from both maps.
// If a key exists in both maps, the value from the other map takes precedence.
// This is inspired by Ruby's Hash#merge method.
func (m *Map[K, V]) Merge(other *Map[K, V]) *Map[K, V] {
	result := &Map[K, V]{}

	// Copy all entries from this map
	m.store.Range(func(k, v any) bool {
		result.Store(k.(K), v.(V))
		return true
	})

	// Copy all entries from other map, overwriting conflicts
	other.store.Range(func(k, v any) bool {
		result.Store(k.(K), v.(V))
		return true
	})

	return result
}

// Clone creates and returns a shallow copy of the map.
// All key-value pairs from the original map are copied to the new map.
// Changes to the cloned map will not affect the original map and vice versa.
func (m *Map[K, V]) Clone() *Map[K, V] {
	result := &Map[K, V]{}
	m.store.Range(func(k, v any) bool {
		result.Store(k.(K), v.(V))
		return true
	})
	return result
}
