package twocache

import (
	"context"
	"reflect"
	"time"

	"go.uber.org/zap"
)

// nullMarker is a special byte prefix that distinguishes cached null values
// from legitimate data. \x00 avoids collisions with JSON output which never
// starts with a null byte.
const nullMarker = "\x00NULL\x00"

// TwoLevelCache combines a BigCache local store (L1) with an optional Redis
// remote store (L2). GetOrLoad transparently falls through L1 → L2 → loader
// and backfills the layers it traversed.
type TwoLevelCache[K comparable, V any] struct {
	local      *bigcacheStore
	remote     *redisStore
	keyEncoder KeyEncoder[K]
	codec      ValueCodec[V]
	logger     *zap.Logger
	redisTTL   time.Duration
	localTTL   time.Duration
	nullTTL    time.Duration
}

// New creates a TwoLevelCache. When WithRedisClient is not provided, the cache
// operates in local-only mode (L1 only).
func New[K comparable, V any](opts ...Option[K, V]) (*TwoLevelCache[K, V], error) {
	cfg := newDefaultConfig[K, V]()
	for _, opt := range opts {
		opt(&cfg)
	}

	local, err := newBigCacheStore(cfg.shards, cfg.maxEntries, cfg.localTTL)
	if err != nil {
		return nil, err
	}

	c := &TwoLevelCache[K, V]{
		local:      local,
		remote:     newRedisStore(cfg.redisClient, cfg.redisTTL),
		keyEncoder: cfg.keyEncoder,
		codec:      cfg.codec,
		logger:     cfg.logger,
		redisTTL:   cfg.redisTTL,
		localTTL:   cfg.localTTL,
		nullTTL:    cfg.nullTTL,
	}
	if c.logger == nil {
		c.logger = zap.NewNop()
	}
	return c, nil
}

// GetOrLoad returns the cached value for key, loading it via loader on miss.
// L1 → L2 → loader. Each layer hit backfills the preceding layers.
func (c *TwoLevelCache[K, V]) GetOrLoad(
	ctx context.Context,
	key K,
	loader func(ctx context.Context, key K) (V, error),
) (V, error) {
	var zero V
	strKey := c.keyEncoder(key)

	// L1: BigCache
	if data, err := c.local.Get(strKey); err == nil {
		if isNullMarker(data) {
			return zero, nil
		}
		val, decodeErr := c.codec.Decode(data)
		if decodeErr == nil {
			return val, nil
		}
		c.logger.Warn("twocache: L1 decode failed, falling through",
			zap.String("key", strKey), zap.Error(decodeErr))
	}

	// L2: Redis
	if data, err := c.remote.Get(ctx, strKey); err == nil {
		if isNullMarker(data) {
			_ = c.local.Set(strKey, []byte(nullMarker))
			return zero, nil
		}
		val, decodeErr := c.codec.Decode(data)
		if decodeErr == nil {
			if encoded, encErr := c.codec.Encode(val); encErr == nil {
				_ = c.local.Set(strKey, encoded)
			}
			return val, nil
		}
		c.logger.Warn("twocache: L2 decode failed, falling through",
			zap.String("key", strKey), zap.Error(decodeErr))
	}

	// Loader
	val, err := loader(ctx, key)
	if err != nil {
		return zero, err
	}

	if isZero(val) {
		_ = c.local.Set(strKey, []byte(nullMarker))
		_ = c.remote.Set(ctx, strKey, []byte(nullMarker), c.nullTTL)
		return val, nil
	}

	encoded, encErr := c.codec.Encode(val)
	if encErr != nil {
		c.logger.Warn("twocache: encode failed after load",
			zap.String("key", strKey), zap.Error(encErr))
		return val, nil
	}
	_ = c.local.Set(strKey, encoded)
	_ = c.remote.Set(ctx, strKey, encoded, c.redisTTL)
	return val, nil
}

// Set writes the value to both L1 and L2 synchronously.
func (c *TwoLevelCache[K, V]) Set(ctx context.Context, key K, val V) error {
	strKey := c.keyEncoder(key)
	encoded, err := c.codec.Encode(val)
	if err != nil {
		return err
	}
	if err := c.local.Set(strKey, encoded); err != nil {
		return err
	}
	return c.remote.Set(ctx, strKey, encoded, c.redisTTL)
}

// Delete removes the value from both L1 and L2 synchronously.
func (c *TwoLevelCache[K, V]) Delete(ctx context.Context, key K) error {
	strKey := c.keyEncoder(key)
	_ = c.local.Delete(strKey)
	return c.remote.Delete(ctx, strKey)
}

// Close releases BigCache resources. Redis client is not closed (caller owns it).
func (c *TwoLevelCache[K, V]) Close() error {
	return c.local.Close()
}

func isNullMarker(data []byte) bool {
	return string(data) == nullMarker
}

func isZero[V any](v V) bool {
	var zero V
	return any(v) == nil || reflect.DeepEqual(v, zero)
}
