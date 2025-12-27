package slices

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContains(t *testing.T) {
	t.Run("element exists", func(t *testing.T) {
		assert.True(t, Contains([]int{1, 2, 3}, 2))
		assert.True(t, Contains([]string{"a", "b", "c"}, "b"))
	})

	t.Run("element does not exist", func(t *testing.T) {
		assert.False(t, Contains([]int{1, 2, 3}, 4))
		assert.False(t, Contains([]string{"a", "b", "c"}, "d"))
	})

	t.Run("empty slice", func(t *testing.T) {
		assert.False(t, Contains([]int{}, 1))
		assert.False(t, Contains([]string(nil), "a"))
	})
}

func TestContainsAny(t *testing.T) {
	t.Run("one of the values exists", func(t *testing.T) {
		assert.True(t, ContainsAny([]int{1, 2, 3}, 2, 4))
		assert.True(t, ContainsAny([]string{"a", "b", "c"}, "d", "b"))
	})

	t.Run("none of the values exist", func(t *testing.T) {
		assert.False(t, ContainsAny([]int{1, 2, 3}, 4, 5))
		assert.False(t, ContainsAny([]string{"a", "b", "c"}, "d", "e"))
	})

	t.Run("all values exist", func(t *testing.T) {
		assert.True(t, ContainsAny([]int{1, 2, 3}, 1, 2))
	})

	t.Run("no values to check", func(t *testing.T) {
		assert.False(t, ContainsAny([]int{1, 2, 3}))
	})
}

func TestContainsAll(t *testing.T) {
	t.Run("all values exist", func(t *testing.T) {
		assert.True(t, ContainsAll([]int{1, 2, 3}, 1, 2))
		assert.True(t, ContainsAll([]string{"a", "b", "c"}, "a", "b", "c"))
	})

	t.Run("not all values exist", func(t *testing.T) {
		assert.False(t, ContainsAll([]int{1, 2, 3}, 1, 4))
		assert.False(t, ContainsAll([]string{"a", "b", "c"}, "a", "d"))
	})

	t.Run("no values to check", func(t *testing.T) {
		assert.True(t, ContainsAll([]int{1, 2, 3}))
	})
}

func TestFilter(t *testing.T) {
	t.Run("filter even numbers", func(t *testing.T) {
		result := Filter([]int{1, 2, 3, 4, 5}, func(n int) bool {
			return n%2 == 0
		})
		assert.Equal(t, []int{2, 4}, result)
	})

	t.Run("filter positive numbers", func(t *testing.T) {
		result := Filter([]int{-1, 0, 1, 2, -3}, func(n int) bool {
			return n >= 0
		})
		assert.Equal(t, []int{0, 1, 2}, result)
	})

	t.Run("no matches", func(t *testing.T) {
		result := Filter([]int{1, 2, 3}, func(n int) bool {
			return n > 10
		})
		assert.Empty(t, result)
	})

	t.Run("empty slice", func(t *testing.T) {
		result := Filter([]int{}, func(n int) bool {
			return true
		})
		assert.Empty(t, result)
	})
}

func TestMap(t *testing.T) {
	t.Run("double numbers", func(t *testing.T) {
		result := Map([]int{1, 2, 3}, func(n int) int {
			return n * 2
		})
		assert.Equal(t, []int{2, 4, 6}, result)
	})

	t.Run("convert to string", func(t *testing.T) {
		result := Map([]int{1, 2, 3}, func(n int) string {
			return "num"
		})
		assert.Equal(t, []string{"num", "num", "num"}, result)
	})

	t.Run("empty slice", func(t *testing.T) {
		result := Map([]int{}, func(n int) int {
			return n * 2
		})
		assert.Empty(t, result)
	})
}

func TestReduce(t *testing.T) {
	t.Run("sum numbers", func(t *testing.T) {
		result := Reduce([]int{1, 2, 3, 4}, 0, func(acc, n int) int {
			return acc + n
		})
		assert.Equal(t, 10, result)
	})

	t.Run("concatenate strings", func(t *testing.T) {
		result := Reduce([]string{"a", "b", "c"}, "", func(acc, s string) string {
			return acc + s
		})
		assert.Equal(t, "abc", result)
	})

	t.Run("empty slice returns initial", func(t *testing.T) {
		result := Reduce([]int{}, 5, func(acc, n int) int {
			return acc + n
		})
		assert.Equal(t, 5, result)
	})
}

func TestFind(t *testing.T) {
	t.Run("element found", func(t *testing.T) {
		result, found := Find([]int{1, 2, 3, 4}, func(n int) bool {
			return n%2 == 0
		})
		assert.True(t, found)
		assert.Equal(t, 2, result)
	})

	t.Run("element not found", func(t *testing.T) {
		result, found := Find([]int{1, 2, 3}, func(n int) bool {
			return n > 10
		})
		assert.False(t, found)
		assert.Equal(t, 0, result)
	})

	t.Run("empty slice", func(t *testing.T) {
		result, found := Find([]int{}, func(n int) bool {
			return true
		})
		assert.False(t, found)
		assert.Equal(t, 0, result)
	})
}

func TestFindIndex(t *testing.T) {
	t.Run("index found", func(t *testing.T) {
		idx := FindIndex([]int{1, 2, 3, 4}, func(n int) bool {
			return n%2 == 0
		})
		assert.Equal(t, 1, idx)
	})

	t.Run("index not found", func(t *testing.T) {
		idx := FindIndex([]int{1, 2, 3}, func(n int) bool {
			return n > 10
		})
		assert.Equal(t, -1, idx)
	})
}

func TestIndexOf(t *testing.T) {
	t.Run("element exists", func(t *testing.T) {
		assert.Equal(t, 1, IndexOf([]int{1, 2, 3}, 2))
		assert.Equal(t, 0, IndexOf([]string{"a", "b", "c"}, "a"))
	})

	t.Run("element does not exist", func(t *testing.T) {
		assert.Equal(t, -1, IndexOf([]int{1, 2, 3}, 4))
		assert.Equal(t, -1, IndexOf([]string{"a", "b", "c"}, "d"))
	})

	t.Run("duplicate elements returns first", func(t *testing.T) {
		assert.Equal(t, 1, IndexOf([]int{1, 2, 2, 3}, 2))
	})
}

func TestLastIndexOf(t *testing.T) {
	t.Run("element exists", func(t *testing.T) {
		assert.Equal(t, 2, LastIndexOf([]int{1, 2, 2, 3}, 2))
		assert.Equal(t, 2, LastIndexOf([]string{"a", "b", "a", "c"}, "a"))
	})

	t.Run("element does not exist", func(t *testing.T) {
		assert.Equal(t, -1, LastIndexOf([]int{1, 2, 3}, 4))
	})
}

func TestUnique(t *testing.T) {
	t.Run("remove duplicates", func(t *testing.T) {
		result := Unique([]int{1, 2, 2, 3, 3, 3, 4})
		assert.Equal(t, []int{1, 2, 3, 4}, result)
	})

	t.Run("already unique", func(t *testing.T) {
		result := Unique([]int{1, 2, 3, 4})
		assert.Equal(t, []int{1, 2, 3, 4}, result)
	})

	t.Run("empty slice", func(t *testing.T) {
		result := Unique([]int{})
		assert.Empty(t, result)
	})

	t.Run("strings", func(t *testing.T) {
		result := Unique([]string{"a", "b", "a", "c", "b"})
		assert.Equal(t, []string{"a", "b", "c"}, result)
	})
}

func TestReverse(t *testing.T) {
	t.Run("reverse slice", func(t *testing.T) {
		result := Reverse([]int{1, 2, 3, 4})
		assert.Equal(t, []int{4, 3, 2, 1}, result)
	})

	t.Run("single element", func(t *testing.T) {
		result := Reverse([]int{1})
		assert.Equal(t, []int{1}, result)
	})

	t.Run("empty slice", func(t *testing.T) {
		result := Reverse([]int{})
		assert.Empty(t, result)
	})
}

func TestChunk(t *testing.T) {
	t.Run("chunk into groups", func(t *testing.T) {
		result := Chunk([]int{1, 2, 3, 4, 5}, 2)
		assert.Equal(t, [][]int{{1, 2}, {3, 4}, {5}}, result)
	})

	t.Run("chunk size 1", func(t *testing.T) {
		result := Chunk([]int{1, 2, 3}, 1)
		assert.Equal(t, [][]int{{1}, {2}, {3}}, result)
	})

	t.Run("chunk larger than slice", func(t *testing.T) {
		result := Chunk([]int{1, 2, 3}, 10)
		assert.Equal(t, [][]int{{1, 2, 3}}, result)
	})

	t.Run("zero or negative size", func(t *testing.T) {
		result := Chunk([]int{1, 2, 3}, 0)
		assert.Empty(t, result)

		result = Chunk([]int{1, 2, 3}, -1)
		assert.Empty(t, result)
	})

	t.Run("empty slice", func(t *testing.T) {
		result := Chunk([]int{}, 2)
		assert.Empty(t, result)
	})
}

func TestFlatten(t *testing.T) {
	t.Run("flatten 2D slice", func(t *testing.T) {
		result := Flatten([][]int{{1, 2}, {3, 4}, {5}})
		assert.Equal(t, []int{1, 2, 3, 4, 5}, result)
	})

	t.Run("empty slices", func(t *testing.T) {
		result := Flatten([][]int{{}, {1, 2}, {}, {3}, {}})
		assert.Equal(t, []int{1, 2, 3}, result)
	})

	t.Run("all empty", func(t *testing.T) {
		result := Flatten([][]int{{}, {}, {}})
		assert.Empty(t, result)
	})
}

func TestShuffle(t *testing.T) {
	t.Run("shuffles elements", func(t *testing.T) {
		original := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		result := Shuffle(original)

		// Same length
		assert.Len(t, result, len(original))

		// Same elements (count)
		for _, v := range original {
			count := 0
			for _, r := range result {
				if r == v {
					count++
				}
			}
			assert.Equal(t, 1, count)
		}

		// Not the same order (very unlikely to be the same)
		same := true
		for i := range original {
			if original[i] != result[i] {
				same = false
				break
			}
		}
		assert.False(t, same, "shuffle should change order")
	})

	t.Run("empty slice", func(t *testing.T) {
		result := Shuffle([]int{})
		assert.Empty(t, result)
	})

	t.Run("single element", func(t *testing.T) {
		result := Shuffle([]int{1})
		assert.Equal(t, []int{1}, result)
	})
}

func TestToMap(t *testing.T) {
	t.Run("convert slice to map", func(t *testing.T) {
		type Person struct {
			ID   int
			Name string
		}
		people := []Person{{1, "Alice"}, {2, "Bob"}}
		result := ToMap(people, func(p Person) int {
			return p.ID
		})

		assert.Len(t, result, 2)
		assert.Equal(t, people[0], result[1])
		assert.Equal(t, people[1], result[2])
	})
}

func TestGroupBy(t *testing.T) {
	t.Run("group by key", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}
		people := []Person{
			{"Alice", 25},
			{"Bob", 30},
			{"Charlie", 25},
		}
		result := GroupBy(people, func(p Person) int {
			return p.Age
		})

		assert.Len(t, result, 2)
		assert.Len(t, result[25], 2)
		assert.Len(t, result[30], 1)
		assert.Equal(t, "Alice", result[25][0].Name)
		assert.Equal(t, "Charlie", result[25][1].Name)
		assert.Equal(t, "Bob", result[30][0].Name)
	})
}

func TestDelete(t *testing.T) {
	t.Run("delete single element", func(t *testing.T) {
		result := Delete([]int{1, 2, 3, 4}, 1)
		assert.Equal(t, []int{1, 3, 4}, result)
	})

	t.Run("delete multiple elements", func(t *testing.T) {
		result := Delete([]int{1, 2, 3, 4, 5}, 1, 3)
		assert.Equal(t, []int{1, 3, 5}, result)
	})

	t.Run("delete out of range indices", func(t *testing.T) {
		result := Delete([]int{1, 2, 3}, -1, 5)
		assert.Equal(t, []int{1, 2, 3}, result)
	})

	t.Run("delete all", func(t *testing.T) {
		result := Delete([]int{1, 2, 3}, 0, 1, 2)
		assert.Empty(t, result)
	})
}

func TestInsert(t *testing.T) {
	t.Run("insert at beginning", func(t *testing.T) {
		result := Insert([]int{2, 3}, 0, 1)
		assert.Equal(t, []int{1, 2, 3}, result)
	})

	t.Run("insert at end", func(t *testing.T) {
		result := Insert([]int{1, 2}, 2, 3)
		assert.Equal(t, []int{1, 2, 3}, result)
	})

	t.Run("insert in middle", func(t *testing.T) {
		result := Insert([]int{1, 4}, 1, 2, 3)
		assert.Equal(t, []int{1, 2, 3, 4}, result)
	})

	t.Run("insert negative index", func(t *testing.T) {
		result := Insert([]int{1, 2, 3}, -1, 0)
		assert.Equal(t, []int{0, 1, 2, 3}, result)
	})

	t.Run("insert beyond length", func(t *testing.T) {
		result := Insert([]int{1, 2}, 10, 3)
		assert.Equal(t, []int{1, 2, 3}, result)
	})
}

func TestConcat(t *testing.T) {
	t.Run("concat multiple slices", func(t *testing.T) {
		result := Concat([]int{1, 2}, []int{3}, []int{4, 5})
		assert.Equal(t, []int{1, 2, 3, 4, 5}, result)
	})

	t.Run("concat with empty slices", func(t *testing.T) {
		result := Concat([]int{}, []int{1, 2}, []int{})
		assert.Equal(t, []int{1, 2}, result)
	})

	t.Run("concat all empty", func(t *testing.T) {
		result := Concat([]int{}, []int{}, []int{})
		assert.Empty(t, result)
	})
}

func TestIntersect(t *testing.T) {
	t.Run("intersection of two slices", func(t *testing.T) {
		result := Intersect([]int{1, 2, 3, 4}, []int{3, 4, 5, 6})
		assert.Equal(t, []int{3, 4}, result)
	})

	t.Run("no intersection", func(t *testing.T) {
		result := Intersect([]int{1, 2, 3}, []int{4, 5, 6})
		assert.Empty(t, result)
	})

	t.Run("with duplicates", func(t *testing.T) {
		result := Intersect([]int{1, 2, 2, 3}, []int{2, 2, 4})
		assert.Equal(t, []int{2}, result)
	})
}

func TestDifference(t *testing.T) {
	t.Run("elements in a not in b", func(t *testing.T) {
		result := Difference([]int{1, 2, 3, 4}, []int{3, 4, 5, 6})
		assert.Equal(t, []int{1, 2}, result)
	})

	t.Run("all elements in b", func(t *testing.T) {
		result := Difference([]int{1, 2}, []int{1, 2, 3})
		assert.Empty(t, result)
	})

	t.Run("no overlap", func(t *testing.T) {
		result := Difference([]int{1, 2, 3}, []int{4, 5, 6})
		assert.Equal(t, []int{1, 2, 3}, result)
	})
}

func TestEach(t *testing.T) {
	t.Run("iterate over slice", func(t *testing.T) {
		sum := 0
		Each([]int{1, 2, 3, 4}, func(n int) {
			sum += n
		})
		assert.Equal(t, 10, sum)
	})
}

func TestEachIndex(t *testing.T) {
	t.Run("iterate with index", func(t *testing.T) {
		indices := []int{}
		EachIndex([]int{10, 20, 30}, func(i int, n int) {
			indices = append(indices, i)
		})
		assert.Equal(t, []int{0, 1, 2}, indices)
	})
}

func TestAny(t *testing.T) {
	t.Run("any element satisfies", func(t *testing.T) {
		assert.True(t, Any([]int{1, 2, 3, 4}, func(n int) bool {
			return n > 3
		}))
	})

	t.Run("no element satisfies", func(t *testing.T) {
		assert.False(t, Any([]int{1, 2, 3}, func(n int) bool {
			return n > 10
		}))
	})
}

func TestAll(t *testing.T) {
	t.Run("all elements satisfy", func(t *testing.T) {
		assert.True(t, All([]int{2, 4, 6, 8}, func(n int) bool {
			return n%2 == 0
		}))
	})

	t.Run("not all elements satisfy", func(t *testing.T) {
		assert.False(t, All([]int{1, 2, 3, 4}, func(n int) bool {
			return n%2 == 0
		}))
	})
}

func TestNone(t *testing.T) {
	t.Run("no element satisfies", func(t *testing.T) {
		assert.True(t, None([]int{1, 3, 5}, func(n int) bool {
			return n%2 == 0
		}))
	})

	t.Run("some elements satisfy", func(t *testing.T) {
		assert.False(t, None([]int{1, 2, 3}, func(n int) bool {
			return n%2 == 0
		}))
	})
}

func TestCount(t *testing.T) {
	t.Run("count matching elements", func(t *testing.T) {
		count := Count([]int{1, 2, 3, 4, 5, 6}, func(n int) bool {
			return n%2 == 0
		})
		assert.Equal(t, 3, count)
	})

	t.Run("no matches", func(t *testing.T) {
		count := Count([]int{1, 2, 3}, func(n int) bool {
			return n > 10
		})
		assert.Equal(t, 0, count)
	})
}

func TestMax(t *testing.T) {
	t.Run("find max", func(t *testing.T) {
		assert.Equal(t, 9, Max([]int{3, 9, 1, 7}))
		assert.Equal(t, 9.5, Max([]float64{3.2, 9.5, 1.1}))
	})

	t.Run("empty slice returns zero", func(t *testing.T) {
		assert.Equal(t, 0, Max([]int{}))
		assert.Equal(t, 0.0, Max([]float64{}))
	})
}

func TestMin(t *testing.T) {
	t.Run("find min", func(t *testing.T) {
		assert.Equal(t, 1, Min([]int{3, 9, 1, 7}))
		assert.Equal(t, 1.1, Min([]float64{3.2, 9.5, 1.1}))
	})

	t.Run("empty slice returns zero", func(t *testing.T) {
		assert.Equal(t, 0, Min([]int{}))
		assert.Equal(t, 0.0, Min([]float64{}))
	})
}

func TestSum(t *testing.T) {
	t.Run("sum integers", func(t *testing.T) {
		assert.Equal(t, 10, Sum([]int{1, 2, 3, 4}))
	})

	t.Run("sum floats", func(t *testing.T) {
		assert.Equal(t, 10.0, Sum([]float64{1.5, 2.5, 3.0, 3.0}))
	})

	t.Run("empty slice returns zero", func(t *testing.T) {
		assert.Equal(t, 0, Sum([]int{}))
	})
}

func TestAverage(t *testing.T) {
	t.Run("average of numbers", func(t *testing.T) {
		assert.Equal(t, 2.5, Average([]int{1, 2, 3, 4}))
		assert.Equal(t, 2.5, Average([]float64{1.0, 2.0, 3.0, 4.0}))
	})

	t.Run("empty slice returns zero", func(t *testing.T) {
		assert.Equal(t, 0.0, Average([]int{}))
	})
}
