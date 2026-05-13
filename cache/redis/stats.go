package redis

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/Martindeeepdark/go-common/cache/defs"
)

func (c *Client) Stats(ctx context.Context) (*defs.CacheStats, error) {
	stats := &defs.CacheStats{}

	n, err := c.rdb.DBSize(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("cache stats: %w", err)
	}
	stats.KeyCount = n

	info, err := c.rdb.Info(ctx, "stats", "memory").Result()
	if err != nil {
		return stats, nil
	}

	for _, line := range strings.Split(info, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "#") || !strings.Contains(line, ":") {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		k, v := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])

		switch k {
		case "keyspace_hits":
			stats.Hits, _ = strconv.ParseInt(v, 10, 64)
		case "keyspace_misses":
			stats.Misses, _ = strconv.ParseInt(v, 10, 64)
		case "used_memory":
			stats.MemoryUsed, _ = strconv.ParseInt(v, 10, 64)
		}
	}

	total := stats.Hits + stats.Misses
	if total > 0 {
		stats.HitRatio = float64(stats.Hits) / float64(total)
	}

	return stats, nil
}

func (c *Client) Keys(ctx context.Context, pattern string) ([]string, error) {
	return c.rdb.Keys(ctx, pattern).Result()
}
