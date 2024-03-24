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
}
