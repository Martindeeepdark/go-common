package twocache

import (
	"errors"
	"time"

	"github.com/allegro/bigcache/v3"
)

var errNotFound = errors.New("cache: key not found")

type bigcacheStore struct {
	cache *bigcache.BigCache
}

func newBigCacheStore(shards, maxEntries int, ttl time.Duration) (*bigcacheStore, error) {
	if shards <= 0 {
		shards = 64
	}
	if maxEntries <= 0 {
		maxEntries = 1024
	}
	config := bigcache.Config{
		Shards:             shards,
		LifeWindow:         ttl,
		MaxEntriesInWindow: maxEntries,
		MaxEntrySize:       500,
		Verbose:            false,
	}
	bc, err := bigcache.NewBigCache(config)
	if err != nil {
		return nil, err
	}
	return &bigcacheStore{cache: bc}, nil
}

func (s *bigcacheStore) Get(key string) ([]byte, error) {
	data, err := s.cache.Get(key)
	if err != nil {
		if errors.Is(err, bigcache.ErrEntryNotFound) {
			return nil, errNotFound
		}
		return nil, errNotFound
	}
	return data, nil
}

func (s *bigcacheStore) Set(key string, data []byte) error {
	return s.cache.Set(key, data)
}

func (s *bigcacheStore) Delete(key string) error {
	return s.cache.Delete(key)
}

func (s *bigcacheStore) Close() error {
	return s.cache.Close()
}
