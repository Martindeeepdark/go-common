package twocache

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// KeyEncoder converts a generic key to a string for BigCache/Redis.
type KeyEncoder[K comparable] func(key K) string

// ValueCodec handles serialization of cache values.
type ValueCodec[V any] interface {
	Encode(v V) ([]byte, error)
	Decode(data []byte) (V, error)
}

// JSONCodec is the default JSON-based codec.
type JSONCodec[V any] struct{}

func (JSONCodec[V]) Encode(v V) ([]byte, error) { return json.Marshal(v) }
func (JSONCodec[V]) Decode(data []byte) (V, error) {
	var v V
	err := json.Unmarshal(data, &v)
	return v, err
}

// Option configures a TwoLevelCache.
type Option[K comparable, V any] func(*config[K, V])

type config[K comparable, V any] struct {
	keyEncoder KeyEncoder[K]
	codec      ValueCodec[V]
	redisClient *redis.Client
	redisTTL    time.Duration
	localTTL    time.Duration
	nullTTL     time.Duration
	shards      int
	maxEntries  int
	logger      *zap.Logger
}

func newDefaultConfig[K comparable, V any]() config[K, V] {
	return config[K, V]{
		keyEncoder: func(key K) string { return fmt.Sprintf("%v", key) },
		codec:      JSONCodec[V]{},
		redisTTL:   30 * time.Minute,
		localTTL:   5 * time.Minute,
		nullTTL:    60 * time.Second,
		shards:     64,
		maxEntries: 1024,
		logger:     zap.NewNop(),
	}
}

func WithKeyEncoder[K comparable, V any](enc KeyEncoder[K]) Option[K, V] {
	return func(c *config[K, V]) { c.keyEncoder = enc }
}

func WithValueCodec[K comparable, V any](codec ValueCodec[V]) Option[K, V] {
	return func(c *config[K, V]) { c.codec = codec }
}

func WithRedisClient[K comparable, V any](client *redis.Client) Option[K, V] {
	return func(c *config[K, V]) { c.redisClient = client }
}

func WithRedisTTL[K comparable, V any](ttl time.Duration) Option[K, V] {
	return func(c *config[K, V]) { c.redisTTL = ttl }
}

func WithLocalTTL[K comparable, V any](ttl time.Duration) Option[K, V] {
	return func(c *config[K, V]) { c.localTTL = ttl }
}

func WithNullTTL[K comparable, V any](ttl time.Duration) Option[K, V] {
	return func(c *config[K, V]) { c.nullTTL = ttl }
}

func WithShards[K comparable, V any](n int) Option[K, V] {
	return func(c *config[K, V]) { c.shards = n }
}

func WithMaxEntries[K comparable, V any](n int) Option[K, V] {
	return func(c *config[K, V]) { c.maxEntries = n }
}

func WithLogger[K comparable, V any](logger *zap.Logger) Option[K, V] {
	return func(c *config[K, V]) { c.logger = logger }
}
