package twocache

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisStore struct {
	client *redis.Client
	ttl    time.Duration
}

func newRedisStore(client *redis.Client, ttl time.Duration) *redisStore {
	if client == nil {
		return nil
	}
	return &redisStore{client: client, ttl: ttl}
}

func (s *redisStore) Get(ctx context.Context, key string) ([]byte, error) {
	if s == nil || s.client == nil {
		return nil, errNotFound
	}
	data, err := s.client.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, errNotFound
		}
		return nil, err
	}
	return data, nil
}

func (s *redisStore) Set(ctx context.Context, key string, data []byte, ttl time.Duration) error {
	if s == nil || s.client == nil {
		return nil
	}
	return s.client.Set(ctx, key, data, ttl).Err()
}

func (s *redisStore) Delete(ctx context.Context, key string) error {
	if s == nil || s.client == nil {
		return nil
	}
	return s.client.Del(ctx, key).Err()
}
