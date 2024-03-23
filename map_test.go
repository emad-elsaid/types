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
}
