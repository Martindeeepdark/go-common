package redis

import (
	"time"

	goredis "github.com/redis/go-redis/v9"
)

type Config struct {
	Addr         string
	Password     string
	DB           int
	PoolSize     int
	MinIdleConns int
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type Option func(*Config)

func WithPassword(password string) Option {
	return func(c *Config) { c.Password = password }
}

func WithDB(db int) Option {
	return func(c *Config) { c.DB = db }
}

func WithPoolSize(size int) Option {
	return func(c *Config) { c.PoolSize = size }
}

func WithTimeouts(dial, read, write time.Duration) Option {
	return func(c *Config) {
		c.DialTimeout = dial
		c.ReadTimeout = read
		c.WriteTimeout = write
	}
}

func defaultConfig() *Config {
	return &Config{
		Addr:         "localhost:6379",
		DB:           0,
		PoolSize:     10,
		MinIdleConns: 5,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	}
}

func newGoRedisClient(cfg *Config) *goredis.Client {
	return goredis.NewClient(&goredis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	})
}
