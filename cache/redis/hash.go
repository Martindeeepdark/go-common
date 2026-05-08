package redis

import (
	"context"
	"encoding/json"
	"fmt"
)

func (c *Client) HGet(ctx context.Context, key, field string, dest interface{}) error {
	val, err := c.rdb.HGet(ctx, key, field).Result()
	if err != nil {
		return fmt.Errorf("cache hget %s.%s: %w", key, field, err)
	}
	return unmarshal(val, dest)
}

func (c *Client) HSet(ctx context.Context, key, field string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("cache hset %s.%s: %w", key, field, err)
	}
	return c.rdb.HSet(ctx, key, field, data).Err()
}

func (c *Client) HDel(ctx context.Context, key, field string) error {
	return c.rdb.HDel(ctx, key, field).Err()
}

func (c *Client) HGetAll(ctx context.Context, key string) (map[string]interface{}, error) {
	result, err := c.rdb.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("cache hgetall %s: %w", key, err)
	}
	m := make(map[string]interface{}, len(result))
	for k, v := range result {
		var val interface{}
		if err := json.Unmarshal([]byte(v), &val); err != nil {
			m[k] = v
		} else {
			m[k] = val
		}
	}
	return m, nil
}

func (c *Client) HExists(ctx context.Context, key, field string) (bool, error) {
	return c.rdb.HExists(ctx, key, field).Result()
}
