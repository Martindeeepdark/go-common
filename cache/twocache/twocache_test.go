package twocache

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupRedis(t *testing.T) (*redis.Client, func()) {
	t.Helper()
	mr := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	return client, func() {
		client.Close()
		mr.Close()
	}
}

func intEncoder(key int) string { return "key:" + string(rune('0'+key)) }

func newTestCache[V any](t *testing.T, client *redis.Client, opts ...Option[int, V]) *TwoLevelCache[int, V] {
	t.Helper()
	defaultOpts := []Option[int, V]{
		WithKeyEncoder[int, V](intEncoder),
		WithLocalTTL[int, V](5 * time.Minute),
		WithRedisTTL[int, V](30 * time.Minute),
		WithShards[int, V](4),
		WithMaxEntries[int, V](100),
	}
	if client != nil {
		defaultOpts = append(defaultOpts, WithRedisClient[int, V](client))
	}
	allOpts := append(defaultOpts, opts...)
	c, err := New[int, V](allOpts...)
	require.NoError(t, err)
	t.Cleanup(func() { c.Close() })
	return c
}

func TestNew_TwoLevelCache(t *testing.T) {
	client, cleanup := setupRedis(t)
	defer cleanup()

	c := newTestCache[string](t, client)
	assert.NotNil(t, c)
}

func TestL1Hit(t *testing.T) {
	client, cleanup := setupRedis(t)
	defer cleanup()

	c := newTestCache[string](t, client)
	ctx := context.Background()

	var loaderCalls atomic.Int32

	// Populate L1 via Set
	err := c.Set(ctx, 1, "hello")
	require.NoError(t, err)

	// GetOrLoad should hit L1, not call loader
	val, err := c.GetOrLoad(ctx, 1, func(ctx context.Context, key int) (string, error) {
		loaderCalls.Add(1)
		return "from-loader", nil
	})

	assert.NoError(t, err)
	assert.Equal(t, "hello", val)
	assert.Equal(t, int32(0), loaderCalls.Load())
}

func TestL2Hit_BackfillL1(t *testing.T) {
	client, cleanup := setupRedis(t)
	defer cleanup()

	c := newTestCache[string](t, client)
	ctx := context.Background()

	// Write directly to Redis (bypass L1)
	err := client.Set(ctx, "key:2", "\"from-redis\"", 30*time.Minute).Err()
	require.NoError(t, err)

	var loaderCalls atomic.Int32

	val, err := c.GetOrLoad(ctx, 2, func(ctx context.Context, key int) (string, error) {
		loaderCalls.Add(1)
		return "from-loader", nil
	})

	assert.NoError(t, err)
	assert.Equal(t, "from-redis", val)
	assert.Equal(t, int32(0), loaderCalls.Load())

	// Second call should hit L1 (backfilled)
	val2, err := c.GetOrLoad(ctx, 2, func(ctx context.Context, key int) (string, error) {
		loaderCalls.Add(1)
		return "from-loader", nil
	})
	assert.NoError(t, err)
	assert.Equal(t, "from-redis", val2)
	assert.Equal(t, int32(0), loaderCalls.Load())
}

func TestLoaderFallback(t *testing.T) {
	client, cleanup := setupRedis(t)
	defer cleanup()

	c := newTestCache[string](t, client)
	ctx := context.Background()

	var loaderCalls atomic.Int32

	val, err := c.GetOrLoad(ctx, 3, func(ctx context.Context, key int) (string, error) {
		loaderCalls.Add(1)
		return "from-loader", nil
	})

	assert.NoError(t, err)
	assert.Equal(t, "from-loader", val)
	assert.Equal(t, int32(1), loaderCalls.Load())

	// Second call should hit L1 (was backfilled)
	val2, err := c.GetOrLoad(ctx, 3, func(ctx context.Context, key int) (string, error) {
		loaderCalls.Add(1)
		return "unexpected", nil
	})
	assert.NoError(t, err)
	assert.Equal(t, "from-loader", val2)
	assert.Equal(t, int32(1), loaderCalls.Load())
}

func TestNullValueCache(t *testing.T) {
	client, cleanup := setupRedis(t)
	defer cleanup()

	c := newTestCache[*string](t, client) // pointer type, nil = zero value
	ctx := context.Background()

	var loaderCalls atomic.Int32

	// loader returns nil (zero value for *string)
	val, err := c.GetOrLoad(ctx, 4, func(ctx context.Context, key int) (*string, error) {
		loaderCalls.Add(1)
		return nil, nil
	})

	assert.NoError(t, err)
	assert.Nil(t, val)
	assert.Equal(t, int32(1), loaderCalls.Load())

	// Second call should hit null marker cache, not call loader
	val2, err := c.GetOrLoad(ctx, 4, func(ctx context.Context, key int) (*string, error) {
		loaderCalls.Add(1)
		s := "should-not-reach"
		return &s, nil
	})
	assert.NoError(t, err)
	assert.Nil(t, val2)
	assert.Equal(t, int32(1), loaderCalls.Load())
}

func TestSet_SyncWrite(t *testing.T) {
	client, cleanup := setupRedis(t)
	defer cleanup()

	c := newTestCache[string](t, client)
	ctx := context.Background()

	err := c.Set(ctx, 5, "sync-value")
	require.NoError(t, err)

	// Verify L1 (via GetOrLoad without loader)
	val, err := c.GetOrLoad(ctx, 5, func(ctx context.Context, key int) (string, error) {
		return "should-not-reach", errors.New("should not be called")
	})
	assert.Equal(t, "sync-value", val)

	// Verify L2 (via Redis directly)
	data, err := client.Get(ctx, "key:5").Result()
	assert.NoError(t, err)
	assert.Equal(t, "\"sync-value\"", data)
}

func TestDelete_SyncDelete(t *testing.T) {
	client, cleanup := setupRedis(t)
	defer cleanup()

	c := newTestCache[string](t, client)
	ctx := context.Background()

	err := c.Set(ctx, 6, "to-delete")
	require.NoError(t, err)

	err = c.Delete(ctx, 6)
	require.NoError(t, err)

	// Verify L2 deleted
	_, err = client.Get(ctx, "key:6").Result()
	assert.Equal(t, redis.Nil, err)

	// GetOrLoad should fall through to loader
	var loaderCalls atomic.Int32
	val, err := c.GetOrLoad(ctx, 6, func(ctx context.Context, key int) (string, error) {
		loaderCalls.Add(1)
		return "reloaded", nil
	})
	assert.Equal(t, "reloaded", val)
	assert.Equal(t, int32(1), loaderCalls.Load())
}

func TestLocalOnlyMode(t *testing.T) {
	// No Redis client
	c := newTestCache[string](t, nil)
	ctx := context.Background()

	val, err := c.GetOrLoad(ctx, 7, func(ctx context.Context, key int) (string, error) {
		return "local-only", nil
	})
	assert.NoError(t, err)
	assert.Equal(t, "local-only", val)

	// Second call should hit L1
	val2, err := c.GetOrLoad(ctx, 7, func(ctx context.Context, key int) (string, error) {
		return "should-not-reach", nil
	})
	assert.Equal(t, "local-only", val2)
}

func TestConcurrentGetOrLoad(t *testing.T) {
	client, cleanup := setupRedis(t)
	defer cleanup()

	c := newTestCache[string](t, client)
	ctx := context.Background()

	var loaderCalls atomic.Int32
	var wg sync.WaitGroup

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			val, err := c.GetOrLoad(ctx, 8, func(ctx context.Context, key int) (string, error) {
				loaderCalls.Add(1)
				time.Sleep(10 * time.Millisecond) // simulate slow loader
				return "concurrent", nil
			})
			assert.NoError(t, err)
			assert.Equal(t, "concurrent", val)
		}()
	}
	wg.Wait()

	// loader should be called at least once (may be multiple due to race)
	assert.True(t, loaderCalls.Load() >= 1)
}
