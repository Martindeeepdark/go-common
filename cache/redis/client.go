package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"common/cache/defs"

	goredis "github.com/redis/go-redis/v9"
)

type Client struct {
	rdb    *goredis.Client
	tokens sync.Map
}

var (
	_ defs.Cache          = (*Client)(nil)
	_ defs.HashCache      = (*Client)(nil)
	_ defs.ListCache      = (*Client)(nil)
	_ defs.SetCache       = (*Client)(nil)
	_ defs.SortedSetCache = (*Client)(nil)
	_ defs.Lock           = (*Client)(nil)
	_ defs.StatsProvider  = (*Client)(nil)
)

func New(addr string, opts ...Option) *Client {
	cfg := defaultConfig()
	cfg.Addr = addr
	for _, opt := range opts {
		opt(cfg)
	}
	return &Client{rdb: newGoRedisClient(cfg)}
}

func FromClient(rdb *goredis.Client) *Client {
	return &Client{rdb: rdb}
}

func (c *Client) Close() error {
	return c.rdb.Close()
}

func (c *Client) Redis() *goredis.Client {
	return c.rdb
}

func (c *Client) Get(ctx context.Context, key string, dest interface{}) error {
	val, err := c.rdb.Get(ctx, key).Bytes()
	if err != nil {
		return fmt.Errorf("cache get %s: %w", key, err)
	}
	return json.Unmarshal(val, dest)
}

func (c *Client) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("cache set %s: %w", key, err)
	}
	return c.rdb.Set(ctx, key, data, ttl).Err()
}

func (c *Client) Delete(ctx context.Context, key string) error {
	return c.rdb.Del(ctx, key).Err()
}

func (c *Client) Exists(ctx context.Context, key string) (bool, error) {
	n, err := c.rdb.Exists(ctx, key).Result()
	return n > 0, err
}

func (c *Client) Clear(ctx context.Context) error {
	return c.rdb.FlushDB(ctx).Err()
}

func marshal(v interface{}) (string, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func unmarshal(s string, dest interface{}) error {
	return json.Unmarshal([]byte(s), dest)
}
