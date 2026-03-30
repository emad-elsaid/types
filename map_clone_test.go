package types

import "testing"

func TestMapClone(t *testing.T) {
	t.Run("empty map", func(t *testing.T) {
		original := &Map[string, int]{}
		cloned := original.Clone()

		if cloned.Size() != 0 {
			t.Errorf("Expected cloned empty map to have size 0, got %d", cloned.Size())
		}
	})

	t.Run("map with entries", func(t *testing.T) {
		original := &Map[string, int]{}
		original.Store("one", 1)
		original.Store("two", 2)
		original.Store("three", 3)

		cloned := original.Clone()

		// Verify size
		if cloned.Size() != 3 {
			t.Errorf("Expected cloned map to have size 3, got %d", cloned.Size())
		}

		// Verify all entries were copied
		if val, ok := cloned.Load("one"); !ok || val != 1 {
			t.Errorf("Expected cloned map to have key 'one' with value 1, got %v, ok=%v", val, ok)
		}
		if val, ok := cloned.Load("two"); !ok || val != 2 {
			t.Errorf("Expected cloned map to have key 'two' with value 2, got %v, ok=%v", val, ok)
		}
		if val, ok := cloned.Load("three"); !ok || val != 3 {
			t.Errorf("Expected cloned map to have key 'three' with value 3, got %v, ok=%v", val, ok)
		}
	})

	t.Run("independence - modifying clone doesn't affect original", func(t *testing.T) {
		original := &Map[string, int]{}
		original.Store("key", 100)

		cloned := original.Clone()
		cloned.Store("key", 200)
		cloned.Store("new", 300)

		// Original should remain unchanged
		if val, ok := original.Load("key"); !ok || val != 100 {
			t.Errorf("Expected original map key to remain 100, got %v", val)
		}
		if original.Has("new") {
			t.Error("Expected original map not to have 'new' key")
		}

		// Clone should have new values
		if val, ok := cloned.Load("key"); !ok || val != 200 {
			t.Errorf("Expected cloned map key to be 200, got %v", val)
		}
		if val, ok := cloned.Load("new"); !ok || val != 300 {
			t.Errorf("Expected cloned map to have 'new' key with value 300, got %v", val)
		}
	})

	t.Run("independence - modifying original doesn't affect clone", func(t *testing.T) {
		original := &Map[string, int]{}
		original.Store("key", 100)

		cloned := original.Clone()
		
		// Modify original after cloning
		original.Store("key", 500)
		original.Store("another", 600)

		// Clone should retain original values
		if val, ok := cloned.Load("key"); !ok || val != 100 {
			t.Errorf("Expected cloned map key to remain 100, got %v", val)
		}
		if cloned.Has("another") {
			t.Error("Expected cloned map not to have 'another' key")
		}
	})

	t.Run("independence - deleting from clone doesn't affect original", func(t *testing.T) {
		original := &Map[string, int]{}
		original.Store("key1", 1)
		original.Store("key2", 2)

		cloned := original.Clone()
		cloned.Delete("key1")

		// Original should still have both keys
		if !original.Has("key1") {
			t.Error("Expected original map to still have 'key1'")
		}
		if !original.Has("key2") {
			t.Error("Expected original map to still have 'key2'")
		}

		// Clone should only have key2
		if cloned.Has("key1") {
			t.Error("Expected cloned map not to have 'key1'")
		}
		if !cloned.Has("key2") {
			t.Error("Expected cloned map to have 'key2'")
		}
	})

	t.Run("with different types", func(t *testing.T) {
		original := &Map[int, string]{}
		original.Store(1, "one")
		original.Store(2, "two")

		cloned := original.Clone()

		if val, ok := cloned.Load(1); !ok || val != "one" {
			t.Errorf("Expected cloned map to have key 1 with value 'one', got %v", val)
		}
		if val, ok := cloned.Load(2); !ok || val != "two" {
			t.Errorf("Expected cloned map to have key 2 with value 'two', got %v", val)
		}
	})

	t.Run("clear doesn't affect clone", func(t *testing.T) {
		original := &Map[string, int]{}
		original.Store("key", 100)

		cloned := original.Clone()
		original.Clear()

		// Original should be empty
		if original.Size() != 0 {
			t.Errorf("Expected original map to be empty after Clear, got size %d", original.Size())
		}

		// Clone should be unaffected
		if cloned.Size() != 1 {
			t.Errorf("Expected cloned map to still have size 1, got %d", cloned.Size())
		}
		if val, ok := cloned.Load("key"); !ok || val != 100 {
			t.Errorf("Expected cloned map to still have key with value 100, got %v", val)
		}
	})
}
