package redis

import (
	"context"
	"encoding/json"
	"fmt"
)

func (c *Client) SAdd(ctx context.Context, key string, members ...interface{}) error {
	args, err := marshalBatch(members)
	if err != nil {
		return fmt.Errorf("cache sadd %s: %w", key, err)
	}
	return c.rdb.SAdd(ctx, key, args...).Err()
}

func (c *Client) SRem(ctx context.Context, key string, members ...interface{}) error {
	args, err := marshalBatch(members)
	if err != nil {
		return fmt.Errorf("cache srem %s: %w", key, err)
	}
	return c.rdb.SRem(ctx, key, args...).Err()
}

func (c *Client) SMembers(ctx context.Context, key string, dest interface{}) error {
	result, err := c.rdb.SMembers(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("cache smembers %s: %w", key, err)
	}
	data, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("cache smembers %s: %w", key, err)
	}
	return json.Unmarshal(data, dest)
}

func (c *Client) SIsMember(ctx context.Context, key string, member interface{}) (bool, error) {
	data, err := json.Marshal(member)
	if err != nil {
		return false, fmt.Errorf("cache sismember %s: %w", key, err)
	}
	return c.rdb.SIsMember(ctx, key, data).Result()
}

func (c *Client) SCard(ctx context.Context, key string) (int64, error) {
	return c.rdb.SCard(ctx, key).Result()
}
