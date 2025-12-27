package maps

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeys(t *testing.T) {
	t.Run("get all keys", func(t *testing.T) {
		m := map[int]string{1: "a", 2: "b", 3: "c"}
		keys := Keys(m)

		assert.Len(t, keys, 3)
		assert.Contains(t, keys, 1)
		assert.Contains(t, keys, 2)
		assert.Contains(t, keys, 3)
	})

	t.Run("empty map", func(t *testing.T) {
		keys := Keys(map[int]string{})
		assert.Empty(t, keys)

		keys = Keys[int, string](nil)
		assert.Empty(t, keys)
	})
}

func TestValues(t *testing.T) {
	t.Run("get all values", func(t *testing.T) {
		m := map[int]string{1: "a", 2: "b", 3: "c"}
		values := Values(m)

		assert.Len(t, values, 3)
		assert.Contains(t, values, "a")
		assert.Contains(t, values, "b")
		assert.Contains(t, values, "c")
	})

	t.Run("empty map", func(t *testing.T) {
		values := Values(map[int]string{})
		assert.Empty(t, values)
	})
}

func TestEntries(t *testing.T) {
	t.Run("get all entries", func(t *testing.T) {
		m := map[int]string{1: "a", 2: "b"}
		entries := Entries(m)

		assert.Len(t, entries, 2)

		entryMap := FromEntries(entries)
		assert.Equal(t, m, entryMap)
	})

	t.Run("empty map", func(t *testing.T) {
		entries := Entries(map[int]string{})
		assert.Empty(t, entries)
	})
}

func TestFromEntries(t *testing.T) {
	t.Run("create map from entries", func(t *testing.T) {
		entries := []Entry[int, string]{
			{Key: 1, Value: "a"},
			{Key: 2, Value: "b"},
		}

		m := FromEntries(entries)
		assert.Equal(t, "a", m[1])
		assert.Equal(t, "b", m[2])
	})

	t.Run("empty entries", func(t *testing.T) {
		m := FromEntries([]Entry[int, string]{})
		assert.Empty(t, m)
	})
}

func TestClone(t *testing.T) {
	t.Run("clone map", func(t *testing.T) {
		original := map[int]string{1: "a", 2: "b"}
		cloned := Clone(original)

		assert.Equal(t, original, cloned)

		// Modify original, clone should not be affected
		original[1] = "x"
		assert.Equal(t, "a", cloned[1])
	})

	t.Run("clone empty map", func(t *testing.T) {
		cloned := Clone(map[int]string{})
		assert.Empty(t, cloned)
	})
}

func TestMerge(t *testing.T) {
	t.Run("merge two maps", func(t *testing.T) {
		m1 := map[int]string{1: "a", 2: "b"}
		m2 := map[int]string{2: "x", 3: "c"}

		result := Merge(m1, m2)

		assert.Equal(t, "a", result[1])
		assert.Equal(t, "x", result[2]) // m2 overrides
		assert.Equal(t, "c", result[3])
	})

	t.Run("merge three maps", func(t *testing.T) {
		m1 := map[int]string{1: "a"}
		m2 := map[int]string{2: "b"}
		m3 := map[int]string{3: "c"}

		result := Merge(m1, m2, m3)

		assert.Equal(t, "a", result[1])
		assert.Equal(t, "b", result[2])
		assert.Equal(t, "c", result[3])
	})

	t.Run("merge empty maps", func(t *testing.T) {
		result := Merge[int, string]()
		assert.Empty(t, result)
	})
}

func TestMergeInto(t *testing.T) {
	t.Run("merge src into dest", func(t *testing.T) {
		dest := map[int]string{1: "a", 2: "b"}
		src := map[int]string{2: "x", 3: "c"}

		MergeInto(dest, src)

		assert.Equal(t, "a", dest[1])
		assert.Equal(t, "x", dest[2]) // overridden
		assert.Equal(t, "c", dest[3])
	})

	t.Run("merge empty src", func(t *testing.T) {
		dest := map[int]string{1: "a"}
		src := map[int]string{}

		MergeInto(dest, src)

		assert.Equal(t, "a", dest[1])
	})
}

func TestHasKey(t *testing.T) {
	t.Run("key exists", func(t *testing.T) {
		m := map[int]string{1: "a", 2: "b"}
		assert.True(t, HasKey(m, 1))
		assert.True(t, HasKey(m, 2))
	})

	t.Run("key does not exist", func(t *testing.T) {
		m := map[int]string{1: "a"}
		assert.False(t, HasKey(m, 2))
	})

	t.Run("nil map", func(t *testing.T) {
		assert.False(t, HasKey[int, string](nil, 1))
	})
}

func TestGetOrDefault(t *testing.T) {
	t.Run("key exists", func(t *testing.T) {
		m := map[int]string{1: "a", 2: "b"}
		assert.Equal(t, "a", GetOrDefault(m, 1, "default"))
	})

	t.Run("key does not exist", func(t *testing.T) {
		m := map[int]string{1: "a"}
		assert.Equal(t, "default", GetOrDefault(m, 2, "default"))
	})

	t.Run("nil map", func(t *testing.T) {
		assert.Equal(t, "default", GetOrDefault[int, string](nil, 1, "default"))
	})
}

func TestGetOrSet(t *testing.T) {
	t.Run("key exists", func(t *testing.T) {
		m := map[int]string{1: "a"}
		result := GetOrSet(m, 1, "default")

		assert.Equal(t, "a", result)
		assert.Equal(t, "a", m[1])
	})

	t.Run("key does not exist", func(t *testing.T) {
		m := map[int]string{1: "a"}
		result := GetOrSet(m, 2, "default")

		assert.Equal(t, "default", result)
		assert.Equal(t, "default", m[2])
	})
}

func TestDeleteKeys(t *testing.T) {
	t.Run("delete multiple keys", func(t *testing.T) {
		m := map[int]string{1: "a", 2: "b", 3: "c", 4: "d"}
		DeleteKeys(m, 1, 3)

		_, has1 := m[1]
		_, has2 := m[2]
		_, has3 := m[3]
		_, has4 := m[4]

		assert.False(t, has1)
		assert.True(t, has2)
		assert.False(t, has3)
		assert.True(t, has4)
	})

	t.Run("delete non-existent keys", func(t *testing.T) {
		m := map[int]string{1: "a"}
		DeleteKeys(m, 2, 3)

		_, has1 := m[1]
		assert.True(t, has1)
	})

	t.Run("delete from nil map", func(t *testing.T) {
		var m map[int]string
		assert.NotPanics(t, func() {
			DeleteKeys(m, 1, 2)
		})
	})
}

func TestKeepOnly(t *testing.T) {
	t.Run("keep only specified keys", func(t *testing.T) {
		m := map[int]string{1: "a", 2: "b", 3: "c", 4: "d"}
		KeepOnly(m, 1, 3)

		assert.Len(t, m, 2)
		assert.Equal(t, "a", m[1])
		assert.Equal(t, "c", m[3])
	})

	t.Run("keep non-existent keys", func(t *testing.T) {
		m := map[int]string{1: "a", 2: "b"}
		KeepOnly(m, 3, 4)

		assert.Empty(t, m)
	})
}

func TestClear(t *testing.T) {
	t.Run("clear map", func(t *testing.T) {
		m := map[int]string{1: "a", 2: "b", 3: "c"}
		Clear(m)

		assert.Empty(t, m)
	})

	t.Run("clear nil map", func(t *testing.T) {
		var m map[int]string
		assert.NotPanics(t, func() {
			Clear(m)
		})
	})
}

func TestIsEmpty(t *testing.T) {
	t.Run("empty map", func(t *testing.T) {
		assert.True(t, IsEmpty(map[int]string{}))
		assert.True(t, IsEmpty[int, string](nil))
	})

	t.Run("non-empty map", func(t *testing.T) {
		assert.False(t, IsEmpty(map[int]string{1: "a"}))
	})
}

func TestSize(t *testing.T) {
	t.Run("get size", func(t *testing.T) {
		assert.Equal(t, 0, Size(map[int]string{}))
		assert.Equal(t, 3, Size(map[int]string{1: "a", 2: "b", 3: "c"}))
	})
}

func TestFilter(t *testing.T) {
	t.Run("filter entries", func(t *testing.T) {
		m := map[int]string{1: "a", 2: "b", 3: "c", 4: "d"}
		result := Filter(m, func(k int, v string) bool {
			return k%2 == 0
		})

		assert.Len(t, result, 2)
		assert.Equal(t, "b", result[2])
		assert.Equal(t, "d", result[4])
	})

	t.Run("no matches", func(t *testing.T) {
		m := map[int]string{1: "a", 2: "b"}
		result := Filter(m, func(k int, v string) bool {
			return k > 10
		})
		assert.Empty(t, result)
	})
}

func TestMapValues(t *testing.T) {
	t.Run("transform values", func(t *testing.T) {
		m := map[int]int{1: 2, 2: 3, 3: 4}
		result := MapValues(m, func(k int, v int) int {
			return v * 2
		})

		assert.Equal(t, 4, result[1])
		assert.Equal(t, 6, result[2])
		assert.Equal(t, 8, result[3])
	})
}

func TestMapKeys(t *testing.T) {
	t.Run("transform keys", func(t *testing.T) {
		m := map[int]string{1: "a", 2: "b"}
		result := MapKeys(m, func(k int, v string) string {
			return v
		})

		assert.Equal(t, "a", result["a"])
		assert.Equal(t, "b", result["b"])
	})
}

func TestInvert(t *testing.T) {
	t.Run("invert map", func(t *testing.T) {
		m := map[string]int{"a": 1, "b": 2}
		result := Invert(m)

		assert.Equal(t, "a", result[1])
		assert.Equal(t, "b", result[2])
	})

	t.Run("with duplicate values", func(t *testing.T) {
		m := map[string]int{"a": 1, "b": 1}
		result := Invert(m)

		// Only one key will remain (undefined which)
		assert.Len(t, result, 1)
	})
}

func TestEach(t *testing.T) {
	t.Run("iterate over entries", func(t *testing.T) {
		m := map[int]string{1: "a", 2: "b", 3: "c"}
		sum := 0
		Each(m, func(k int, v string) {
			sum += k
		})

		assert.Equal(t, 6, sum)
	})
}

func TestAny(t *testing.T) {
	t.Run("any entry satisfies", func(t *testing.T) {
		m := map[int]int{1: 10, 2: 20, 3: 30}
		assert.True(t, Any(m, func(k int, v int) bool {
			return v > 25
		}))
	})

	t.Run("no entry satisfies", func(t *testing.T) {
		m := map[int]int{1: 10, 2: 20}
		assert.False(t, Any(m, func(k int, v int) bool {
			return v > 100
		}))
	})
}

func TestAll(t *testing.T) {
	t.Run("all entries satisfy", func(t *testing.T) {
		m := map[int]int{1: 10, 2: 20, 3: 30}
		assert.True(t, All(m, func(k int, v int) bool {
			return v > 5
		}))
	})

	t.Run("not all entries satisfy", func(t *testing.T) {
		m := map[int]int{1: 10, 2: 20}
		assert.False(t, All(m, func(k int, v int) bool {
			return v > 15
		}))
	})
}

func TestEqual(t *testing.T) {
	t.Run("maps are equal", func(t *testing.T) {
		m1 := map[int]string{1: "a", 2: "b"}
		m2 := map[int]string{1: "a", 2: "b"}
		assert.True(t, Equal(m1, m2))
	})

	t.Run("maps have different sizes", func(t *testing.T) {
		m1 := map[int]string{1: "a"}
		m2 := map[int]string{1: "a", 2: "b"}
		assert.False(t, Equal(m1, m2))
	})

	t.Run("maps have different values", func(t *testing.T) {
		m1 := map[int]string{1: "a", 2: "b"}
		m2 := map[int]string{1: "a", 2: "x"}
		assert.False(t, Equal(m1, m2))
	})

	t.Run("maps have different keys", func(t *testing.T) {
		m1 := map[int]string{1: "a", 2: "b"}
		m2 := map[int]string{1: "a", 3: "b"}
		assert.False(t, Equal(m1, m2))
	})

	t.Run("both empty", func(t *testing.T) {
		assert.True(t, Equal(map[int]string{}, map[int]string{}))
	})
}

func TestDiff(t *testing.T) {
	t.Run("entries in a not in b", func(t *testing.T) {
		a := map[int]string{1: "a", 2: "b", 3: "c"}
		b := map[int]string{2: "b", 3: "x"}
		result := Diff(a, b)

		assert.Len(t, result, 2)
		assert.Equal(t, "a", result[1])
		assert.Equal(t, "c", result[3])
	})
}

func TestIntersect(t *testing.T) {
	t.Run("entries in both maps", func(t *testing.T) {
		a := map[int]string{1: "a", 2: "b", 3: "c"}
		b := map[int]string{2: "b", 3: "x", 4: "d"}
		result := Intersect(a, b)

		assert.Len(t, result, 1)
		assert.Equal(t, "b", result[2])
	})
}

func TestUnion(t *testing.T) {
	t.Run("union of two maps", func(t *testing.T) {
		a := map[int]string{1: "a", 2: "b"}
		b := map[int]string{2: "x", 3: "c"}
		result := Union(a, b)

		assert.Len(t, result, 3)
		assert.Equal(t, "a", result[1])
		assert.Equal(t, "x", result[2]) // b overrides
		assert.Equal(t, "c", result[3])
	})
}
