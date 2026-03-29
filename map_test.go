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

	t.Run("Any", func(t *testing.T) {
		m := Map[string, int]{}
		m.Store("a", 1)
		m.Store("b", 2)
		m.Store("c", 3)

		// Check if any value is greater than 2
		if !m.Any(func(k string, v int) bool { return v > 2 }) {
			t.Error("Expected Any to return true for predicate v > 2")
		}

		// Check if any value is greater than 10
		if m.Any(func(k string, v int) bool { return v > 10 }) {
			t.Error("Expected Any to return false for predicate v > 10")
		}

		// Check with key condition
		if !m.Any(func(k string, v int) bool { return k == "b" }) {
			t.Error("Expected Any to return true for key 'b'")
		}
	})

	t.Run("Any empty map", func(t *testing.T) {
		m := Map[string, int]{}

		if m.Any(func(k string, v int) bool { return true }) {
			t.Error("Expected Any to return false for empty map")
		}
	})

	t.Run("All", func(t *testing.T) {
		m := Map[string, int]{}
		m.Store("a", 2)
		m.Store("b", 4)
		m.Store("c", 6)

		// Check if all values are even
		if !m.All(func(k string, v int) bool { return v%2 == 0 }) {
			t.Error("Expected All to return true for all even values")
		}

		m.Store("d", 5)
		// Now check with an odd value
		if m.All(func(k string, v int) bool { return v%2 == 0 }) {
			t.Error("Expected All to return false after adding odd value")
		}

		// Check key condition - all keys should be single chars
		if !m.All(func(k string, v int) bool { return len(k) == 1 }) {
			t.Error("Expected All to return true for single-char keys")
		}
	})

	t.Run("All empty map", func(t *testing.T) {
		m := Map[string, int]{}

		// All should return true for empty map (vacuous truth)
		if !m.All(func(k string, v int) bool { return false }) {
			t.Error("Expected All to return true for empty map")
		}
	})

	t.Run("None", func(t *testing.T) {
		m := Map[string, int]{}
		m.Store("a", 1)
		m.Store("b", 2)
		m.Store("c", 3)

		// Check if no value is greater than 10
		if !m.None(func(k string, v int) bool { return v > 10 }) {
			t.Error("Expected None to return true for predicate v > 10")
		}

		// Check if no value is equal to 2
		if m.None(func(k string, v int) bool { return v == 2 }) {
			t.Error("Expected None to return false for predicate v == 2")
		}

		// Check with key condition
		if m.None(func(k string, v int) bool { return k == "a" }) {
			t.Error("Expected None to return false when key 'a' exists")
		}
	})

	t.Run("None empty map", func(t *testing.T) {
		m := Map[string, int]{}

		// None should return true for empty map
		if !m.None(func(k string, v int) bool { return true }) {
			t.Error("Expected None to return true for empty map")
		}
	})

	t.Run("Any/All/None consistency", func(t *testing.T) {
		m := Map[string, int]{}
		m.Store("a", 1)
		m.Store("b", 2)
		m.Store("c", 3)

		predicate := func(k string, v int) bool { return v > 1 }

		// If Any is true and All is false, None must be false
		if m.Any(predicate) && !m.All(predicate) && m.None(predicate) {
			t.Error("Inconsistent behavior: Any is true but None is also true")
		}

		// None should be the opposite of Any
		if m.None(predicate) == m.Any(predicate) {
			t.Error("None and Any returned the same value")
		}
	})

	t.Run("Partition", func(t *testing.T) {
		m := Map[string, int]{}
		m.Store("a", 1)
		m.Store("b", 2)
		m.Store("c", 3)
		m.Store("d", 4)

		// Partition by even/odd
		evens, odds := m.Partition(func(k string, v int) bool {
			return v%2 == 0
		})

		// Check evens map
		if evens.Size() != 2 {
			t.Errorf("Expected 2 even entries, got %d", evens.Size())
		}
		if !evens.Has("b") || !evens.Has("d") {
			t.Error("Expected evens to contain 'b' and 'd'")
		}
		if v, _ := evens.Load("b"); v != 2 {
			t.Errorf("Expected b=2 in evens, got %d", v)
		}
		if v, _ := evens.Load("d"); v != 4 {
			t.Errorf("Expected d=4 in evens, got %d", v)
		}

		// Check odds map
		if odds.Size() != 2 {
			t.Errorf("Expected 2 odd entries, got %d", odds.Size())
		}
		if !odds.Has("a") || !odds.Has("c") {
			t.Error("Expected odds to contain 'a' and 'c'")
		}
		if v, _ := odds.Load("a"); v != 1 {
			t.Errorf("Expected a=1 in odds, got %d", v)
		}
		if v, _ := odds.Load("c"); v != 3 {
			t.Errorf("Expected c=3 in odds, got %d", v)
		}

		// Original map should be unchanged
		if m.Size() != 4 {
			t.Errorf("Expected original map size 4, got %d", m.Size())
		}
	})

	t.Run("Partition all true", func(t *testing.T) {
		m := Map[string, int]{}
		m.Store("a", 2)
		m.Store("b", 4)
		m.Store("c", 6)

		// All values satisfy predicate
		trueMap, falseMap := m.Partition(func(k string, v int) bool {
			return v%2 == 0
		})

		if trueMap.Size() != 3 {
			t.Errorf("Expected trueMap size 3, got %d", trueMap.Size())
		}
		if falseMap.Size() != 0 {
			t.Errorf("Expected falseMap size 0, got %d", falseMap.Size())
		}
	})

	t.Run("Partition all false", func(t *testing.T) {
		m := Map[string, int]{}
		m.Store("a", 1)
		m.Store("b", 3)
		m.Store("c", 5)

		// No values satisfy predicate
		trueMap, falseMap := m.Partition(func(k string, v int) bool {
			return v%2 == 0
		})

		if trueMap.Size() != 0 {
			t.Errorf("Expected trueMap size 0, got %d", trueMap.Size())
		}
		if falseMap.Size() != 3 {
			t.Errorf("Expected falseMap size 3, got %d", falseMap.Size())
		}
	})

	t.Run("Partition empty map", func(t *testing.T) {
		m := Map[string, int]{}

		trueMap, falseMap := m.Partition(func(k string, v int) bool {
			return v > 0
		})

		if trueMap.Size() != 0 {
			t.Errorf("Expected trueMap size 0 for empty map, got %d", trueMap.Size())
		}
		if falseMap.Size() != 0 {
			t.Errorf("Expected falseMap size 0 for empty map, got %d", falseMap.Size())
		}
	})

	t.Run("Partition with key-based predicate", func(t *testing.T) {
		m := Map[string, string]{}
		m.Store("alice", "engineer")
		m.Store("bob", "designer")
		m.Store("charlie", "engineer")
		m.Store("diana", "manager")

		// Partition by names starting with vowels
		vowels, consonants := m.Partition(func(k string, v string) bool {
			first := k[0]
			return first == 'a' || first == 'e' || first == 'i' || first == 'o' || first == 'u'
		})

		if vowels.Size() != 1 {
			t.Errorf("Expected 1 vowel name, got %d", vowels.Size())
		}
		if !vowels.Has("alice") {
			t.Error("Expected vowels to contain 'alice'")
		}

		if consonants.Size() != 3 {
			t.Errorf("Expected 3 consonant names, got %d", consonants.Size())
		}
	})

	t.Run("Partition preserves all entries", func(t *testing.T) {
		m := Map[int, string]{}
		m.Store(1, "one")
		m.Store(2, "two")
		m.Store(3, "three")
		m.Store(4, "four")
		m.Store(5, "five")

		// Partition by value length
		longNames, shortNames := m.Partition(func(k int, v string) bool {
			return len(v) > 3
		})

		// Verify total count
		totalSize := longNames.Size() + shortNames.Size()
		if totalSize != m.Size() {
			t.Errorf("Expected total size %d, got %d", m.Size(), totalSize)
		}

		// Verify no overlap
		for _, k := range longNames.Keys() {
			if shortNames.Has(k) {
				t.Errorf("Key %d exists in both partitions", k)
			}
		}
	})

	t.Run("Merge basic", func(t *testing.T) {
		m1 := Map[string, int]{}
		m1.Store("a", 1)
		m1.Store("b", 2)

		m2 := Map[string, int]{}
		m2.Store("c", 3)
		m2.Store("d", 4)

		merged := m1.Merge(&m2)

		if merged.Size() != 4 {
			t.Errorf("Expected size 4, got %d", merged.Size())
		}

		// Verify all keys exist
		if v, ok := merged.Load("a"); !ok || v != 1 {
			t.Errorf("Expected a=1, got %d", v)
		}
		if v, ok := merged.Load("b"); !ok || v != 2 {
			t.Errorf("Expected b=2, got %d", v)
		}
		if v, ok := merged.Load("c"); !ok || v != 3 {
			t.Errorf("Expected c=3, got %d", v)
		}
		if v, ok := merged.Load("d"); !ok || v != 4 {
			t.Errorf("Expected d=4, got %d", v)
		}
	})

	t.Run("Merge with overlapping keys", func(t *testing.T) {
		m1 := Map[string, int]{}
		m1.Store("a", 1)
		m1.Store("b", 2)
		m1.Store("c", 3)

		m2 := Map[string, int]{}
		m2.Store("b", 20) // Override
		m2.Store("c", 30) // Override
		m2.Store("d", 4)  // New

		merged := m1.Merge(&m2)

		if merged.Size() != 4 {
			t.Errorf("Expected size 4, got %d", merged.Size())
		}

		// Verify m2 values take precedence
		if v, ok := merged.Load("a"); !ok || v != 1 {
			t.Errorf("Expected a=1, got %d", v)
		}
		if v, ok := merged.Load("b"); !ok || v != 20 {
			t.Errorf("Expected b=20 (from m2), got %d", v)
		}
		if v, ok := merged.Load("c"); !ok || v != 30 {
			t.Errorf("Expected c=30 (from m2), got %d", v)
		}
		if v, ok := merged.Load("d"); !ok || v != 4 {
			t.Errorf("Expected d=4, got %d", v)
		}
	})

	t.Run("Merge does not mutate originals", func(t *testing.T) {
		m1 := Map[string, string]{}
		m1.Store("key1", "value1")

		m2 := Map[string, string]{}
		m2.Store("key2", "value2")

		merged := m1.Merge(&m2)

		// Original maps should be unchanged
		if m1.Size() != 1 {
			t.Errorf("m1 should still have size 1, got %d", m1.Size())
		}
		if m2.Size() != 1 {
			t.Errorf("m2 should still have size 1, got %d", m2.Size())
		}
		if merged.Size() != 2 {
			t.Errorf("merged should have size 2, got %d", merged.Size())
		}

		// Verify m1 doesn't have m2's keys
		if m1.Has("key2") {
			t.Error("m1 should not have key2")
		}
		// Verify m2 doesn't have m1's keys
		if m2.Has("key1") {
			t.Error("m2 should not have key1")
		}
	})

	t.Run("Merge empty maps", func(t *testing.T) {
		m1 := Map[string, int]{}
		m2 := Map[string, int]{}

		merged := m1.Merge(&m2)

		if merged.Size() != 0 {
			t.Errorf("Expected size 0, got %d", merged.Size())
		}
	})

	t.Run("Merge with one empty map", func(t *testing.T) {
		m1 := Map[string, int]{}
		m1.Store("a", 1)
		m1.Store("b", 2)

		m2 := Map[string, int]{}

		merged := m1.Merge(&m2)

		if merged.Size() != 2 {
			t.Errorf("Expected size 2, got %d", merged.Size())
		}
		if v, ok := merged.Load("a"); !ok || v != 1 {
			t.Errorf("Expected a=1, got %d", v)
		}
	})
}
