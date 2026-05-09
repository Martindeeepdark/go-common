package snowflake

import (
	"fmt"
	"sync"
	"time"

	"github.com/sony/sonyflake"
)

var (
	defaultFlake *Flake
	once         sync.Once
)

type Option func(*config)

type config struct {
	startTime time.Time
	machineID func() (uint16, error)
}

func WithStartTime(t time.Time) Option {
	return func(c *config) { c.startTime = t }
}

func WithMachineID(fn func() (uint16, error)) Option {
	return func(c *config) { c.machineID = fn }
}

func Init(opts ...Option) {
	once.Do(func() {
		cfg := config{
			startTime: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		}
		for _, opt := range opts {
			opt(&cfg)
		}

		sf := sonyflake.NewSonyflake(sonyflake.Settings{
			StartTime: cfg.startTime,
			MachineID: cfg.machineID,
		})
		if sf == nil {
			panic("snowflake: invalid settings")
		}
		defaultFlake = &Flake{sf: sf}
	})
}

func init() {
	Init()
}

func NewID() int64 {
	if defaultFlake == nil {
		panic("snowflake: not initialized")
	}
	return defaultFlake.MustNewID()
}

type Flake struct {
	sf *sonyflake.Sonyflake
}

func New(opts ...Option) (*Flake, error) {
	cfg := config{
		startTime: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	for _, opt := range opts {
		opt(&cfg)
	}

	sf := sonyflake.NewSonyflake(sonyflake.Settings{
		StartTime: cfg.startTime,
		MachineID: cfg.machineID,
	})
	if sf == nil {
		return nil, fmt.Errorf("snowflake: invalid settings")
	}
	return &Flake{sf: sf}, nil
}

func (f *Flake) NewID() (int64, error) {
	id, err := f.sf.NextID()
	if err != nil {
		return 0, fmt.Errorf("snowflake: %w", err)
	}
	return int64(id), nil
}

func (f *Flake) MustNewID() int64 {
	id, err := f.NewID()
	if err != nil {
		panic(err)
	}
	return id
}

func Decompose(id int64) map[string]uint64 {
	return sonyflake.Decompose(uint64(id))
}
