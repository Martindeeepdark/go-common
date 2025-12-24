package ctxcache

import (
	"context"
	"sync"
)

type ctxCacheKey struct{}

// Init initializes a cache map in the context
func Init(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxCacheKey{}, new(sync.Map))
}

// Get retrieves a value from the context cache
func Get[T any](ctx context.Context, key any) (value T, ok bool) {
	var zero T

	cacheMap, valid := ctx.Value(ctxCacheKey{}).(*sync.Map)
	if !valid {
		return zero, false
	}

	loadedValue, exists := cacheMap.Load(key)
	if !exists {
		return zero, false
	}

	if v, match := loadedValue.(T); match {
		return v, true
	}

	return zero, false
}

// Store stores a value in the context cache
func Store(ctx context.Context, key any, obj any) {
	if cacheMap, ok := ctx.Value(ctxCacheKey{}).(*sync.Map); ok {
		cacheMap.Store(key, obj)
	}
}

// HasKey checks if a key exists in the context cache
func HasKey(ctx context.Context, key any) bool {
	if cacheMap, ok := ctx.Value(ctxCacheKey{}).(*sync.Map); ok {
		_, ok := cacheMap.Load(key)
		return ok
	}

	return false
}

// LoadOrStore returns the existing value for the key if present,
// or stores and returns the given value if not present
func LoadOrStore[T any](ctx context.Context, key any, value T) (actual T, loaded bool) {
	cacheMap, ok := ctx.Value(ctxCacheKey{}).(*sync.Map)
	if !ok {
		return value, false
	}

	actualVal, loaded := cacheMap.LoadOrStore(key, value)
	if loaded {
		if actual, ok := actualVal.(T); ok {
			return actual, true
		}
	}
	return value, false
}

// LoadAndDelete deletes the value for a key, returning the previous value if any
func LoadAndDelete[T any](ctx context.Context, key any) (value T, ok bool) {
	var zero T

	cacheMap, ok := ctx.Value(ctxCacheKey{}).(*sync.Map)
	if !ok {
		return zero, false
	}

	loadedValue, exists := cacheMap.LoadAndDelete(key)
	if !exists {
		return zero, false
	}

	if v, match := loadedValue.(T); match {
		return v, true
	}

	return zero, false
}

// Delete deletes the value for a key
func Delete(ctx context.Context, key any) {
	if cacheMap, ok := ctx.Value(ctxCacheKey{}).(*sync.Map); ok {
		cacheMap.Delete(key)
	}
}

// Range calls f sequentially for each key and value present in the context cache
func Range(ctx context.Context, f func(key, value any) bool) {
	if cacheMap, ok := ctx.Value(ctxCacheKey{}).(*sync.Map); ok {
		cacheMap.Range(f)
	}
}

// Keys returns all keys in the context cache
func Keys(ctx context.Context) []any {
	var keys []any
	Range(ctx, func(key, value any) bool {
		keys = append(keys, key)
		return true
	})
	return keys
}

// Clear removes all entries from the context cache
func Clear(ctx context.Context) {
	if cacheMap, ok := ctx.Value(ctxCacheKey{}).(*sync.Map); ok {
		cacheMap.Range(func(key, value any) bool {
			cacheMap.Delete(key)
			return true
		})
	}
}
