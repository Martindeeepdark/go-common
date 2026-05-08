package redis

import (
	"context"
	"encoding/json"
	"fmt"
)

func marshalBatch(values []interface{}) ([]interface{}, error) {
	args := make([]interface{}, len(values))
	for i, v := range values {
		data, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		args[i] = data
	}
	return args, nil
}

func (c *Client) LPush(ctx context.Context, key string, values ...interface{}) error {
	args, err := marshalBatch(values)
	if err != nil {
		return fmt.Errorf("cache lpush %s: %w", key, err)
	}
	return c.rdb.LPush(ctx, key, args...).Err()
}

func (c *Client) RPush(ctx context.Context, key string, values ...interface{}) error {
	args, err := marshalBatch(values)
	if err != nil {
		return fmt.Errorf("cache rpush %s: %w", key, err)
	}
	return c.rdb.RPush(ctx, key, args...).Err()
}

func (c *Client) LPop(ctx context.Context, key string, dest interface{}) error {
	val, err := c.rdb.LPop(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("cache lpop %s: %w", key, err)
	}
	return unmarshal(val, dest)
}

func (c *Client) RPop(ctx context.Context, key string, dest interface{}) error {
	val, err := c.rdb.RPop(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("cache rpop %s: %w", key, err)
	}
	return unmarshal(val, dest)
}

func (c *Client) LLen(ctx context.Context, key string) (int64, error) {
	return c.rdb.LLen(ctx, key).Result()
}

func (c *Client) LRange(ctx context.Context, key string, start, stop int64, dest interface{}) error {
	result, err := c.rdb.LRange(ctx, key, start, stop).Result()
	if err != nil {
		return fmt.Errorf("cache lrange %s: %w", key, err)
	}
	data, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("cache lrange %s: %w", key, err)
	}
	return json.Unmarshal(data, dest)
}
