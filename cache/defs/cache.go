package defs

import (
	"context"
	"time"
)

// Cache defines the cache interface
type Cache interface {
	// Get retrieves a value from cache
	Get(ctx context.Context, key string, dest interface{}) error
	// Set stores a value in cache
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	// Delete removes a value from cache
	Delete(ctx context.Context, key string) error
	// Exists checks if a key exists in cache
	Exists(ctx context.Context, key string) (bool, error)
	// Clear clears all cache entries
	Clear(ctx context.Context) error
}

// HashCache defines hash operations (like Redis HGET/HSET)
type HashCache interface {
	// HGet retrieves a field value from a hash
	HGet(ctx context.Context, key, field string, dest interface{}) error
	// HSet sets a field value in a hash
	HSet(ctx context.Context, key, field string, value interface{}) error
	// HDel removes a field from a hash
	HDel(ctx context.Context, key, field string) error
	// HGetAll retrieves all fields and values from a hash
	HGetAll(ctx context.Context, key string) (map[string]interface{}, error)
	// HExists checks if a field exists in a hash
	HExists(ctx context.Context, key, field string) (bool, error)
}

// ListCache defines list operations (like Redis LPUSH/RPOP)
type ListCache interface {
	// LPush pushes a value to the left of a list
	LPush(ctx context.Context, key string, values ...interface{}) error
	// RPush pushes a value to the right of a list
	RPush(ctx context.Context, key string, values ...interface{}) error
	// LPop pops a value from the left of a list
	LPop(ctx context.Context, key string, dest interface{}) error
	// RPop pops a value from the right of a list
	RPop(ctx context.Context, key string, dest interface{}) error
	// LLen returns the length of a list
	LLen(ctx context.Context, key string) (int64, error)
	// LRange retrieves a range of elements from a list
	LRange(ctx context.Context, key string, start, stop int64, dest interface{}) error
}

// SetCache defines set operations (like Redis SADD/SMEMBERS)
type SetCache interface {
	// SAdd adds members to a set
	SAdd(ctx context.Context, key string, members ...interface{}) error
	// SRem removes members from a set
	SRem(ctx context.Context, key string, members ...interface{}) error
	// SMembers retrieves all members of a set
	SMembers(ctx context.Context, key string, dest interface{}) error
	// SIsMember checks if a member is in a set
	SIsMember(ctx context.Context, key string, member interface{}) (bool, error)
	// SCard returns the number of members in a set
	SCard(ctx context.Context, key string) (int64, error)
}

// SortedSetCache defines sorted set operations (like Redis ZADD/ZRANGE)
type SortedSetCache interface {
	// ZAdd adds members with scores to a sorted set
	ZAdd(ctx context.Context, key string, members ...*ScoredMember) error
	// ZRem removes members from a sorted set
	ZRem(ctx context.Context, key string, members ...interface{}) error
	// ZRange retrieves members by rank range
	ZRange(ctx context.Context, key string, start, stop int64, dest interface{}) error
	// ZRangeByScore retrieves members by score range
	ZRangeByScore(ctx context.Context, key string, min, max float64, dest interface{}) error
	// ZScore returns the score of a member
	ZScore(ctx context.Context, key string, member interface{}) (float64, error)
}

// ScoredMember represents a member with a score in a sorted set
type ScoredMember struct {
	Member interface{} `json:"member"`
	Score  float64     `json:"score"`
}

// Lock defines the distributed lock interface
type Lock interface {
	// TryLock attempts to acquire a lock
	TryLock(ctx context.Context, key string, ttl time.Duration) (bool, error)
	// Lock acquires a lock (will block until acquired or context cancelled)
	Lock(ctx context.Context, key string, ttl time.Duration) error
	// Unlock releases a lock
	Unlock(ctx context.Context, key string) error
}

// CacheStats represents cache statistics
type CacheStats struct {
	Hits       int64   `json:"hits"`        // Cache hits
	Misses     int64   `json:"misses"`      // Cache misses
	HitRatio   float64 `json:"hit_ratio"`   // Hit ratio (hits / (hits + misses))
	KeyCount   int64   `json:"key_count"`   // Total number of keys
	MemoryUsed int64   `json:"memory_used"` // Memory usage in bytes
}

// StatsProvider defines cache statistics interface
type StatsProvider interface {
	// Stats returns cache statistics
	Stats(ctx context.Context) (*CacheStats, error)
	// Keys returns all keys matching a pattern
	Keys(ctx context.Context, pattern string) ([]string, error)
}
