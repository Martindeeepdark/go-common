package redis

import (
	"context"
	cryptorand "crypto/rand"
	"encoding/hex"
	"fmt"
	mathrand "math/rand/v2"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

var (
	unlockScript = goredis.NewScript(unlockScriptSrc)
	extendScript = goredis.NewScript(extendScriptSrc)
)

func generateToken() string {
	b := make([]byte, 16)
	_, _ = cryptorand.Read(b)
	return hex.EncodeToString(b)
}

func (c *Client) TryLock(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	lockKey := "lock:" + key
	token := generateToken()

	ok, err := c.rdb.SetNX(ctx, lockKey, token, ttl).Result()
	if err != nil {
		return false, fmt.Errorf("cache trylock %s: %w", key, err)
	}

	if ok {
		c.tokens.Store(lockKey, token)
	}
	return ok, nil
}

func (c *Client) Lock(ctx context.Context, key string, ttl time.Duration) error {
	lockKey := "lock:" + key
	token := generateToken()

	baseDelay := 50 * time.Millisecond
	maxDelay := 500 * time.Millisecond
	delay := baseDelay

	for {
		ok, err := c.rdb.SetNX(ctx, lockKey, token, ttl).Result()
		if err != nil {
			return fmt.Errorf("cache lock %s: %w", key, err)
		}
		if ok {
			c.tokens.Store(lockKey, token)
			return nil
		}

		jitter := time.Duration(mathrand.Int64N(int64(delay / 2)))
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay + jitter):
		}

		delay *= 2
		if delay > maxDelay {
			delay = maxDelay
		}
	}
}

func (c *Client) Unlock(ctx context.Context, key string) error {
	lockKey := "lock:" + key

	v, loaded := c.tokens.LoadAndDelete(lockKey)
	if !loaded {
		return fmt.Errorf("cache unlock %s: not lock owner", key)
	}
	token := v.(string)

	result, err := unlockScript.Run(ctx, c.rdb, []string{lockKey}, token).Int64()
	if err != nil {
		return fmt.Errorf("cache unlock %s: %w", key, err)
	}
	if result == 0 {
		return fmt.Errorf("cache unlock %s: lock expired or token mismatch", key)
	}
	return nil
}

func (c *Client) Extend(ctx context.Context, key string, ttl time.Duration) error {
	lockKey := "lock:" + key

	v, ok := c.tokens.Load(lockKey)
	if !ok {
		return fmt.Errorf("cache extend %s: not lock owner", key)
	}
	token := v.(string)

	result, err := extendScript.Run(ctx, c.rdb, []string{lockKey}, token, int(ttl.Milliseconds())).Int64()
	if err != nil {
		return fmt.Errorf("cache extend %s: %w", key, err)
	}
	if result == 0 {
		return fmt.Errorf("cache extend %s: lock expired or token mismatch", key)
	}
	return nil
}
