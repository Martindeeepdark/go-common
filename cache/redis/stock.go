package redis

import (
	"context"
	"fmt"

	goredis "github.com/redis/go-redis/v9"
)

var (
	deductStockScript = goredis.NewScript(deductStockSrc)
	checkLimitScript  = goredis.NewScript(checkLimitSrc)
)

const (
	stockResultKeyMissing   int64 = -1
	stockResultInsufficient int64 = -2
)

// DeductStock atomically checks and deducts stock for a given key.
// Returns the remaining stock after deduction, or an error if key missing or stock insufficient.
func (c *Client) DeductStock(ctx context.Context, key string, quantity int64) (int64, error) {
	result, err := deductStockScript.Run(ctx, c.rdb, []string{key}, quantity).Int64()
	if err != nil {
		return 0, fmt.Errorf("cache deduct_stock %s: %w", key, err)
	}
	switch result {
	case stockResultKeyMissing:
		return 0, fmt.Errorf("cache deduct_stock %s: key not found", key)
	case stockResultInsufficient:
		return 0, fmt.Errorf("cache deduct_stock %s: insufficient stock", key)
	}
	return result, nil
}

// RestoreStock adds quantity back to the stock key.
// Uses a single INCRBY which is already atomic — no Lua needed.
func (c *Client) RestoreStock(ctx context.Context, key string, quantity int64) error {
	if err := c.rdb.IncrBy(ctx, key, quantity).Err(); err != nil {
		return fmt.Errorf("cache restore_stock %s: %w", key, err)
	}
	return nil
}

// CheckLimit atomically checks if a user's purchase would exceed the limit and increments the counter.
// Key does not expire — in seckill scenarios Redis is the source of truth, expiry would cause DB penetration.
// Returns the updated purchase count, or an error if the limit would be exceeded.
func (c *Client) CheckLimit(ctx context.Context, key string, quantity, limit int64) (int64, error) {
	result, err := checkLimitScript.Run(ctx, c.rdb, []string{key}, quantity, limit).Int64()
	if err != nil {
		return 0, fmt.Errorf("cache check_limit %s: %w", key, err)
	}
	if result == -1 {
		return 0, fmt.Errorf("cache check_limit %s: would exceed limit %d", key, limit)
	}
	return result, nil
}

// RollbackLimit decrements the user's purchase count by quantity.
// Used when stock deduction fails after limit check passed.
func (c *Client) RollbackLimit(ctx context.Context, key string, quantity int64) error {
	if err := c.rdb.IncrBy(ctx, key, -quantity).Err(); err != nil {
		return fmt.Errorf("cache rollback_limit %s: %w", key, err)
	}
	return nil
}
