package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Martindeeepdark/go-common/cache/defs"

	goredis "github.com/redis/go-redis/v9"
)

func (c *Client) ZAdd(ctx context.Context, key string, members ...*defs.ScoredMember) error {
	zs := make([]goredis.Z, len(members))
	for i, m := range members {
		data, err := json.Marshal(m.Member)
		if err != nil {
			return fmt.Errorf("cache zadd %s: %w", key, err)
		}
		zs[i] = goredis.Z{Score: m.Score, Member: data}
	}
	return c.rdb.ZAdd(ctx, key, zs...).Err()
}

func (c *Client) ZRem(ctx context.Context, key string, members ...interface{}) error {
	args, err := marshalBatch(members)
	if err != nil {
		return fmt.Errorf("cache zrem %s: %w", key, err)
	}
	return c.rdb.ZRem(ctx, key, args).Err()
}

func (c *Client) ZRange(ctx context.Context, key string, start, stop int64, dest interface{}) error {
	result, err := c.rdb.ZRange(ctx, key, start, stop).Result()
	if err != nil {
		return fmt.Errorf("cache zrange %s: %w", key, err)
	}
	data, _ := json.Marshal(result)
	return json.Unmarshal(data, dest)
}

func (c *Client) ZRangeByScore(ctx context.Context, key string, min, max float64, dest interface{}) error {
	result, err := c.rdb.ZRangeByScore(ctx, key, &goredis.ZRangeBy{
		Min: fmt.Sprintf("%f", min),
		Max: fmt.Sprintf("%f", max),
	}).Result()
	if err != nil {
		return fmt.Errorf("cache zrangebyscore %s: %w", key, err)
	}
	data, _ := json.Marshal(result)
	return json.Unmarshal(data, dest)
}

func (c *Client) ZScore(ctx context.Context, key string, member interface{}) (float64, error) {
	data, err := json.Marshal(member)
	if err != nil {
		return 0, fmt.Errorf("cache zscore %s: %w", key, err)
	}
	return c.rdb.ZScore(ctx, key, string(data)).Result()
}
