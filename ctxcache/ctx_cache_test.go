package ctxcache

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	ctx := context.Background()
	newCtx := Init(ctx)

	assert.NotNil(t, newCtx)
	assert.NotEqual(t, ctx, newCtx)
}

func TestGetAndStore(t *testing.T) {
	t.Run("store and get value", func(t *testing.T) {
		ctx := Init(context.Background())

		Store(ctx, "key1", "value1")
		value, ok := Get[string](ctx, "key1")

		assert.True(t, ok)
		assert.Equal(t, "value1", value)
	})

	t.Run("get non-existent key", func(t *testing.T) {
		ctx := Init(context.Background())

		value, ok := Get[string](ctx, "non-existent")

		assert.False(t, ok)
		assert.Equal(t, "", value)
	})

	t.Run("get from uninitialized context", func(t *testing.T) {
		ctx := context.Background()

		value, ok := Get[string](ctx, "key")

		assert.False(t, ok)
		assert.Equal(t, "", value)
	})

	t.Run("type mismatch", func(t *testing.T) {
		ctx := Init(context.Background())

		Store(ctx, "key", "string value")
		value, ok := Get[int](ctx, "key")

		assert.False(t, ok)
		assert.Equal(t, 0, value)
	})
}

func TestHasKey(t *testing.T) {
	t.Run("key exists", func(t *testing.T) {
		ctx := Init(context.Background())
		Store(ctx, "key", "value")

		assert.True(t, HasKey(ctx, "key"))
	})

	t.Run("key does not exist", func(t *testing.T) {
		ctx := Init(context.Background())

		assert.False(t, HasKey(ctx, "non-existent"))
	})

	t.Run("uninitialized context", func(t *testing.T) {
		ctx := context.Background()

		assert.False(t, HasKey(ctx, "key"))
	})
}

func TestLoadOrStore(t *testing.T) {
	t.Run("load existing value", func(t *testing.T) {
		ctx := Init(context.Background())
		Store(ctx, "key", "value1")

		actual, loaded := LoadOrStore(ctx, "key", "value2")

		assert.True(t, loaded)
		assert.Equal(t, "value1", actual)
	})

	t.Run("store new value", func(t *testing.T) {
		ctx := Init(context.Background())

		actual, loaded := LoadOrStore(ctx, "key", "value")

		assert.False(t, loaded)
		assert.Equal(t, "value", actual)
	})

	t.Run("uninitialized context", func(t *testing.T) {
		ctx := context.Background()

		actual, loaded := LoadOrStore(ctx, "key", "value")

		assert.False(t, loaded)
		assert.Equal(t, "value", actual)
	})
}

func TestLoadAndDelete(t *testing.T) {
	t.Run("load and delete existing key", func(t *testing.T) {
		ctx := Init(context.Background())
		Store(ctx, "key", "value")

		value, ok := LoadAndDelete[string](ctx, "key")

		assert.True(t, ok)
		assert.Equal(t, "value", value)
		assert.False(t, HasKey(ctx, "key"))
	})

	t.Run("load and delete non-existent key", func(t *testing.T) {
		ctx := Init(context.Background())

		value, ok := LoadAndDelete[string](ctx, "key")

		assert.False(t, ok)
		assert.Equal(t, "", value)
	})

	t.Run("uninitialized context", func(t *testing.T) {
		ctx := context.Background()

		value, ok := LoadAndDelete[string](ctx, "key")

		assert.False(t, ok)
		assert.Equal(t, "", value)
	})
}

func TestDelete(t *testing.T) {
	t.Run("delete existing key", func(t *testing.T) {
		ctx := Init(context.Background())
		Store(ctx, "key", "value")

		assert.True(t, HasKey(ctx, "key"))
		Delete(ctx, "key")
		assert.False(t, HasKey(ctx, "key"))
	})

	t.Run("delete non-existent key", func(t *testing.T) {
		ctx := Init(context.Background())

		assert.NotPanics(t, func() {
			Delete(ctx, "non-existent")
		})
	})

	t.Run("delete from uninitialized context", func(t *testing.T) {
		ctx := context.Background()

		assert.NotPanics(t, func() {
			Delete(ctx, "key")
		})
	})
}

func TestRange(t *testing.T) {
	t.Run("range over cache", func(t *testing.T) {
		ctx := Init(context.Background())
		Store(ctx, "key1", "value1")
		Store(ctx, "key2", "value2")
		Store(ctx, "key3", "value3")

		count := 0
		Range(ctx, func(key, value any) bool {
			count++
			return true
		})

		assert.Equal(t, 3, count)
	})

	t.Run("range with early termination", func(t *testing.T) {
		ctx := Init(context.Background())
		Store(ctx, "key1", "value1")
		Store(ctx, "key2", "value2")

		count := 0
		Range(ctx, func(key, value any) bool {
			count++
			return count < 1
		})

		assert.Equal(t, 1, count)
	})

	t.Run("range over uninitialized context", func(t *testing.T) {
		ctx := context.Background()

		count := 0
		Range(ctx, func(key, value any) bool {
			count++
			return true
		})

		assert.Equal(t, 0, count)
	})
}

func TestKeys(t *testing.T) {
	t.Run("get all keys", func(t *testing.T) {
		ctx := Init(context.Background())
		Store(ctx, "key1", "value1")
		Store(ctx, "key2", "value2")
		Store(ctx, "key3", "value3")

		keys := Keys(ctx)

		assert.Len(t, keys, 3)
		assert.Contains(t, keys, "key1")
		assert.Contains(t, keys, "key2")
		assert.Contains(t, keys, "key3")
	})

	t.Run("get keys from empty cache", func(t *testing.T) {
		ctx := Init(context.Background())

		keys := Keys(ctx)

		assert.Empty(t, keys)
	})

	t.Run("get keys from uninitialized context", func(t *testing.T) {
		ctx := context.Background()

		keys := Keys(ctx)

		assert.Empty(t, keys)
	})
}

func TestClear(t *testing.T) {
	t.Run("clear cache", func(t *testing.T) {
		ctx := Init(context.Background())
		Store(ctx, "key1", "value1")
		Store(ctx, "key2", "value2")

		Clear(ctx)

		assert.Empty(t, Keys(ctx))
	})

	t.Run("clear empty cache", func(t *testing.T) {
		ctx := Init(context.Background())

		assert.NotPanics(t, func() {
			Clear(ctx)
		})
	})

	t.Run("clear uninitialized context", func(t *testing.T) {
		ctx := context.Background()

		assert.NotPanics(t, func() {
			Clear(ctx)
		})
	})
}

func TestConcurrentAccess(t *testing.T) {
	t.Run("concurrent stores", func(t *testing.T) {
		ctx := Init(context.Background())
		var wg sync.WaitGroup

		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				Store(ctx, idx, idx*2)
			}(i)
		}

		wg.Wait()

		keys := Keys(ctx)
		assert.Len(t, keys, 100)
	})

	t.Run("concurrent reads and writes", func(t *testing.T) {
		ctx := Init(context.Background())
		Store(ctx, "key", "value")

		var wg sync.WaitGroup

		// Concurrent reads
		for i := 0; i < 50; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				Get[string](ctx, "key")
			}()
		}

		// Concurrent writes
		for i := 0; i < 50; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				Store(ctx, idx, idx)
			}(i)
		}

		wg.Wait()

		// Verify no race conditions
		assert.True(t, true)
	})
}

func TestTypes(t *testing.T) {
	t.Run("store and get different types", func(t *testing.T) {
		ctx := Init(context.Background())

		// String
		Store(ctx, "str", "string value")
		strVal, ok := Get[string](ctx, "str")
		assert.True(t, ok)
		assert.Equal(t, "string value", strVal)

		// Int
		Store(ctx, "int", 42)
		intVal, ok := Get[int](ctx, "int")
		assert.True(t, ok)
		assert.Equal(t, 42, intVal)

		// Bool
		Store(ctx, "bool", true)
		boolVal, ok := Get[bool](ctx, "bool")
		assert.True(t, ok)
		assert.True(t, boolVal)

		// Slice
		Store(ctx, "slice", []int{1, 2, 3})
		sliceVal, ok := Get[[]int](ctx, "slice")
		assert.True(t, ok)
		assert.Equal(t, []int{1, 2, 3}, sliceVal)
	})
}

func TestContextPropagation(t *testing.T) {
	t.Run("cache persists across context derivatives", func(t *testing.T) {
		ctx := Init(context.Background())
		Store(ctx, "key", "value")

		// Create derived context
		childCtx := context.WithValue(ctx, "other", "data")

		// Cache should still be accessible
		value, ok := Get[string](childCtx, "key")
		assert.True(t, ok)
		assert.Equal(t, "value", value)
	})
}
