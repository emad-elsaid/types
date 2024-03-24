package types

import "sync"

// Map is a wrapper for sync.Map with same methods (except for CompareAnd*)
// Check: https://pkg.go.dev/sync#Map
type Map[K comparable, V any] struct {
	store sync.Map
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
