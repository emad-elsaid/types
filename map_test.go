package types

import "testing"

func TestMap(t *testing.T) {
	t.Run("Store/Load", func(t *testing.T) {
		m := Map[string, string]{}
		m.Store("Key1", "Value1")
		m.Store("Key2", "Value2")

		if v, ok := m.Load("Key1"); v != "Value1" || !ok {
			t.Errorf("Expected (%s, %t) but got (%s,%t)", v, ok, v, ok)
		}

		if v, ok := m.Load("Key2"); v != "Value2" || !ok {
			t.Errorf("Expected (%s, %t) but got (%s,%t)", v, ok, v, ok)
		}

		if v, ok := m.Load("Key3"); v != "" || ok {
			t.Errorf("Expected (%s, %t) but got (%s,%t)", v, ok, v, ok)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		m := Map[string, string]{}
		m.Store("Key1", "Value1")
		m.Delete("Key1")

		if v, ok := m.Load("Key1"); ok {
			t.Errorf("Key shouldn't exist but it was found: %v", v)
		}
	})

	t.Run("Range", func(t *testing.T) {
		m := Map[string, string]{}
		m.Store("Key1", "Value1")
		m.Store("Key2", "Value2")

		cpy := Map[string, string]{}
		m.Range(func(k string, v string) bool {
			cpy.Store(k, v)
			return true
		})

		if v, ok := cpy.Load("Key1"); v != "Value1" || !ok {
			t.Errorf("Expected: Value1, got: %s", v)
		}

		if v, ok := cpy.Load("Key2"); v != "Value2" || !ok {
			t.Errorf("Expected: Value2, got: %s", v)
		}
	})

	t.Run("LoadAndDelete", func(t *testing.T) {
		m := Map[string, string]{}
		m.Store("Key1", "Value1")

		if v, loaded := m.LoadAndDelete("Key1"); v != "Value1" || !loaded {
			t.Errorf("Expected: Value1 and true, got: %s, %t", v, loaded)
		}

		if v, loaded := m.LoadAndDelete("Key2"); v != "" || loaded {
			t.Errorf("Expected: empty and false, got: %s, %t", v, loaded)
		}
	})

	t.Run("LoadOrStore", func(t *testing.T) {
		m := Map[string, string]{}
		m.Store("Key1", "Value1")

		if v, loaded := m.LoadOrStore("Key1", "Value1/1"); v != "Value1" || !loaded {
			t.Errorf("Expected Value1 and loaded, got: %s, %t", v, loaded)
		}

		if v, loaded := m.LoadOrStore("Key2", "Value2"); v != "Value2" || loaded {
			t.Errorf("Expected Value2 and  not loaded, got: %s, %t", v, loaded)
		}
	})

	t.Run("Swap", func(t *testing.T) {
		m := Map[string, string]{}
		m.Store("Key1", "Value1")

		if v, loaded := m.Swap("Key1", "Value1/1"); v != "Value1" || !loaded {
			t.Errorf("Expected Value1 and loaded, got: %s, %t", v, loaded)
		}

		if v, loaded := m.Swap("Key2", "Value2"); v != "" || loaded {
			t.Errorf("Expected Value2 and not loaded, got: %s, %t", v, loaded)
		}
	})

	t.Run("Keys", func(t *testing.T) {
		m := Map[string, int]{}
		m.Store("a", 1)
		m.Store("b", 2)
		m.Store("c", 3)

		keys := m.Keys()
		if len(keys) != 3 {
			t.Errorf("Expected 3 keys, got %d", len(keys))
		}

		// Verify all expected keys are present
		keyMap := make(map[string]bool)
		for _, k := range keys {
			keyMap[k] = true
		}
		if !keyMap["a"] || !keyMap["b"] || !keyMap["c"] {
			t.Errorf("Missing expected keys, got: %v", keys)
		}
	})

	t.Run("Values", func(t *testing.T) {
		m := Map[string, int]{}
		m.Store("a", 1)
		m.Store("b", 2)
		m.Store("c", 3)

		values := m.Values()
		if len(values) != 3 {
			t.Errorf("Expected 3 values, got %d", len(values))
		}

		// Verify all expected values are present
		valueMap := make(map[int]bool)
		for _, v := range values {
			valueMap[v] = true
		}
		if !valueMap[1] || !valueMap[2] || !valueMap[3] {
			t.Errorf("Missing expected values, got: %v", values)
		}
	})

	t.Run("Size", func(t *testing.T) {
		m := Map[string, int]{}
		if m.Size() != 0 {
			t.Errorf("Expected size 0 for empty map, got %d", m.Size())
		}

		m.Store("a", 1)
		if m.Size() != 1 {
			t.Errorf("Expected size 1, got %d", m.Size())
		}

		m.Store("b", 2)
		m.Store("c", 3)
		if m.Size() != 3 {
			t.Errorf("Expected size 3, got %d", m.Size())
		}

		m.Delete("a")
		if m.Size() != 2 {
			t.Errorf("Expected size 2 after delete, got %d", m.Size())
		}
	})

	t.Run("Clear", func(t *testing.T) {
		m := Map[string, int]{}
		m.Store("a", 1)
		m.Store("b", 2)
		m.Store("c", 3)

		m.Clear()

		if m.Size() != 0 {
			t.Errorf("Expected size 0 after Clear, got %d", m.Size())
		}

		if _, ok := m.Load("a"); ok {
			t.Error("Key 'a' should not exist after Clear")
		}
	})

	t.Run("Has", func(t *testing.T) {
		m := Map[string, int]{}
		m.Store("a", 1)

		if !m.Has("a") {
			t.Error("Expected Has('a') to return true")
		}

		if m.Has("b") {
			t.Error("Expected Has('b') to return false")
		}

		m.Delete("a")
		if m.Has("a") {
			t.Error("Expected Has('a') to return false after delete")
		}
	})

	t.Run("ForEach", func(t *testing.T) {
		m := Map[string, int]{}
		m.Store("a", 1)
		m.Store("b", 2)
		m.Store("c", 3)

		sum := 0
		m.ForEach(func(k string, v int) {
			sum += v
		})

		if sum != 6 {
			t.Errorf("Expected sum 6, got %d", sum)
		}
	})

	t.Run("Filter", func(t *testing.T) {
		m := Map[string, int]{}
		m.Store("a", 1)
		m.Store("b", 2)
		m.Store("c", 3)
		m.Store("d", 4)

		// Filter for even values
		filtered := m.Filter(func(k string, v int) bool {
			return v%2 == 0
		})

		if filtered.Size() != 2 {
			t.Errorf("Expected filtered size 2, got %d", filtered.Size())
		}

		if !filtered.Has("b") || !filtered.Has("d") {
			t.Error("Expected filtered map to contain 'b' and 'd'")
		}

		if filtered.Has("a") || filtered.Has("c") {
			t.Error("Expected filtered map not to contain 'a' or 'c'")
		}

		// Original map should be unchanged
		if m.Size() != 4 {
			t.Errorf("Expected original map size 4, got %d", m.Size())
		}
	})

	t.Run("Transform", func(t *testing.T) {
		m := Map[string, int]{}
		m.Store("a", 1)
		m.Store("b", 2)
		m.Store("c", 3)

		// Transform by doubling values
		transformed := m.Transform(func(k string, v int) int {
			return v * 2
		})

		if transformed.Size() != 3 {
			t.Errorf("Expected transformed size 3, got %d", transformed.Size())
		}

		// Check transformed values
		if v, ok := transformed.Load("a"); !ok || v != 2 {
			t.Errorf("Expected a=2, got %d", v)
		}
		if v, ok := transformed.Load("b"); !ok || v != 4 {
			t.Errorf("Expected b=4, got %d", v)
		}
		if v, ok := transformed.Load("c"); !ok || v != 6 {
			t.Errorf("Expected c=6, got %d", v)
		}

		// Original map should be unchanged
		if v, ok := m.Load("a"); !ok || v != 1 {
			t.Errorf("Expected original a=1, got %d", v)
		}
		if m.Size() != 3 {
			t.Errorf("Expected original map size 3, got %d", m.Size())
		}
	})

	t.Run("Transform with key dependency", func(t *testing.T) {
		m := Map[string, string]{}
		m.Store("name", "john")
		m.Store("city", "berlin")
		m.Store("country", "germany")

		// Transform by capitalizing and prefixing with key
		transformed := m.Transform(func(k string, v string) string {
			return k + ":" + v
		})

		if v, ok := transformed.Load("name"); !ok || v != "name:john" {
			t.Errorf("Expected 'name:john', got %s", v)
		}
		if v, ok := transformed.Load("city"); !ok || v != "city:berlin" {
			t.Errorf("Expected 'city:berlin', got %s", v)
		}
		if v, ok := transformed.Load("country"); !ok || v != "country:germany" {
			t.Errorf("Expected 'country:germany', got %s", v)
		}
	})

	t.Run("Transform empty map", func(t *testing.T) {
		m := Map[string, int]{}

		transformed := m.Transform(func(k string, v int) int {
			return v * 10
		})

		if transformed.Size() != 0 {
			t.Errorf("Expected transformed empty map, got size %d", transformed.Size())
		}
	})

	t.Run("NewMapFrom", func(t *testing.T) {
		source := map[string]int{
			"a": 1,
			"b": 2,
			"c": 3,
		}

		m := NewMapFrom(source)

		if m.Size() != 3 {
			t.Errorf("Expected size 3, got %d", m.Size())
		}

		if v, ok := m.Load("a"); !ok || v != 1 {
			t.Errorf("Expected a=1, got %d", v)
		}
		if v, ok := m.Load("b"); !ok || v != 2 {
			t.Errorf("Expected b=2, got %d", v)
		}
		if v, ok := m.Load("c"); !ok || v != 3 {
			t.Errorf("Expected c=3, got %d", v)
		}
	})

	t.Run("NewMapFrom empty map", func(t *testing.T) {
		source := map[string]int{}
		m := NewMapFrom(source)

		if m.Size() != 0 {
			t.Errorf("Expected size 0 for empty map, got %d", m.Size())
		}
	})

	t.Run("Entries", func(t *testing.T) {
		m := Map[string, int]{}
		m.Store("a", 1)
		m.Store("b", 2)
		m.Store("c", 3)

		entries := m.Entries()

		if len(entries) != 3 {
			t.Errorf("Expected 3 entries, got %d", len(entries))
		}

		if entries["a"] != 1 {
			t.Errorf("Expected a=1, got %d", entries["a"])
		}
		if entries["b"] != 2 {
			t.Errorf("Expected b=2, got %d", entries["b"])
		}
		if entries["c"] != 3 {
			t.Errorf("Expected c=3, got %d", entries["c"])
		}
	})

	t.Run("Entries empty map", func(t *testing.T) {
		m := Map[string, int]{}
		entries := m.Entries()

		if len(entries) != 0 {
			t.Errorf("Expected 0 entries for empty map, got %d", len(entries))
		}
	})

	t.Run("NewMapFrom and Entries roundtrip", func(t *testing.T) {
		original := map[string]int{
			"x": 10,
			"y": 20,
			"z": 30,
		}

		m := NewMapFrom(original)
		result := m.Entries()

		if len(result) != len(original) {
			t.Errorf("Expected %d entries, got %d", len(original), len(result))
		}

		for k, v := range original {
			if result[k] != v {
				t.Errorf("Expected %s=%d, got %d", k, v, result[k])
			}
		}
	})
}
