package taskgroup

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	ctx := context.Background()
	g := New(ctx)

	assert.NotNil(t, g)
	assert.NotNil(t, g.Context())
}

func TestGo(t *testing.T) {
	t.Run("successful task", func(t *testing.T) {
		ctx := context.Background()
		g := New(ctx)

		executed := false
		g.Go(func() error {
			executed = true
			return nil
		})

		err := g.Wait()
		assert.NoError(t, err)
		assert.True(t, executed)
	})

	t.Run("task with error", func(t *testing.T) {
		ctx := context.Background()
		g := New(ctx)

		g.Go(func() error {
			return errors.New("task error")
		})

		err := g.Wait()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "task error")
	})

	t.Run("multiple tasks", func(t *testing.T) {
		ctx := context.Background()
		g := New(ctx)

		count := 0
		for i := 0; i < 5; i++ {
			g.Go(func() error {
				count++
				return nil
			})
		}

		err := g.Wait()
		assert.NoError(t, err)
		assert.Equal(t, 5, count)
	})

	t.Run("one task fails", func(t *testing.T) {
		ctx := context.Background()
		g := New(ctx)

		for i := 0; i < 5; i++ {
			idx := i
			g.Go(func() error {
				if idx == 2 {
					return errors.New("failed")
				}
				return nil
			})
		}

		err := g.Wait()
		assert.Error(t, err)
	})
}

func TestGoWithContext(t *testing.T) {
	t.Run("task with context", func(t *testing.T) {
		ctx := context.Background()
		g := New(ctx)

		executed := false
		g.GoWithContext(func(ctx context.Context) error {
			executed = true
			assert.Equal(t, g.ctx, ctx)
			return nil
		})

		err := g.Wait()
		assert.NoError(t, err)
		assert.True(t, executed)
	})

	t.Run("context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		g := New(ctx)

		g.GoWithContext(func(ctx context.Context) error {
			time.Sleep(100 * time.Millisecond)
			return nil
		})

		g.GoWithContext(func(ctx context.Context) error {
			cancel() // Cancel context
			return errors.New("cancelled")
		})

		err := g.Wait()
		assert.Error(t, err)
	})
}

func TestWait(t *testing.T) {
	t.Run("wait returns first error", func(t *testing.T) {
		ctx := context.Background()
		g := New(ctx)

		g.Go(func() error {
			return errors.New("error1")
		})
		g.Go(func() error {
			return errors.New("error2")
		})

		err := g.Wait()
		assert.Error(t, err)
	})

	t.Run("wait returns nil on success", func(t *testing.T) {
		ctx := context.Background()
		g := New(ctx)

		g.Go(func() error {
			return nil
		})

		err := g.Wait()
		assert.NoError(t, err)
	})
}

func TestWaitAll(t *testing.T) {
	t.Run("wait all returns all errors", func(t *testing.T) {
		ctx := context.Background()
		g := New(ctx)

		g.Go(func() error {
			return errors.New("error1")
		})
		g.Go(func() error {
			return errors.New("error2")
		})
		g.Go(func() error {
			return nil
		})

		errs := g.WaitAll()
		assert.Len(t, errs, 2)
	})

	t.Run("wait all returns empty on success", func(t *testing.T) {
		ctx := context.Background()
		g := New(ctx)

		g.Go(func() error {
			return nil
		})

		errs := g.WaitAll()
		assert.Empty(t, errs)
	})
}

func TestCancel(t *testing.T) {
	t.Run("cancel task group", func(t *testing.T) {
		ctx := context.Background()
		g := New(ctx)

		g.Go(func() error {
			time.Sleep(200 * time.Millisecond)
			return nil
		})

		time.Sleep(50 * time.Millisecond)
		g.Cancel()

		err := g.Wait()
		assert.NoError(t, err) // No error, just cancelled
	})
}

func TestErrors(t *testing.T) {
	t.Run("get errors", func(t *testing.T) {
		ctx := context.Background()
		g := New(ctx)

		g.Go(func() error {
			return errors.New("error1")
		})
		g.Go(func() error {
			return errors.New("error2")
		})

		g.Wait()

		errs := g.Errors()
		assert.Len(t, errs, 2)
	})
}

func TestHasError(t *testing.T) {
	t.Run("has error when task fails", func(t *testing.T) {
		ctx := context.Background()
		g := New(ctx)

		g.Go(func() error {
			return errors.New("error")
		})

		g.Wait()
		assert.True(t, g.HasError())
	})

	t.Run("no error when all succeed", func(t *testing.T) {
		ctx := context.Background()
		g := New(ctx)

		g.Go(func() error {
			return nil
		})

		g.Wait()
		assert.False(t, g.HasError())
	})
}

func TestParallel(t *testing.T) {
	t.Run("parallel execution", func(t *testing.T) {
		ctx := context.Background()
		count := 0

		err := Parallel(ctx,
			func() error {
				count++
				return nil
			},
			func() error {
				count++
				return nil
			},
			func() error {
				count++
				return nil
			},
		)

		assert.NoError(t, err)
		assert.Equal(t, 3, count)
	})

	t.Run("parallel with error", func(t *testing.T) {
		ctx := context.Background()

		err := Parallel(ctx,
			func() error {
				return nil
			},
			func() error {
				return errors.New("failed")
			},
		)

		assert.Error(t, err)
	})
}

func TestParallelWithContext(t *testing.T) {
	t.Run("parallel with context", func(t *testing.T) {
		ctx := context.Background()
		count := 0

		err := ParallelWithContext(ctx,
			func(ctx context.Context) error {
				count++
				return nil
			},
			func(ctx context.Context) error {
				count++
				return nil
			},
		)

		assert.NoError(t, err)
		assert.Equal(t, 2, count)
	})
}

func TestSerial(t *testing.T) {
	t.Run("serial execution", func(t *testing.T) {
		ctx := context.Background()
		order := []int{}

		err := Serial(ctx,
			func() error {
				order = append(order, 1)
				return nil
			},
			func() error {
				order = append(order, 2)
				return nil
			},
			func() error {
				order = append(order, 3)
				return nil
			},
		)

		assert.NoError(t, err)
		assert.Equal(t, []int{1, 2, 3}, order)
	})

	t.Run("serial stops on error", func(t *testing.T) {
		ctx := context.Background()
		count := 0

		err := Serial(ctx,
			func() error {
				count++
				return nil
			},
			func() error {
				count++
				return errors.New("error")
			},
			func() error {
				count++
				return nil
			},
		)

		assert.Error(t, err)
		assert.Equal(t, 2, count) // First two executed
	})
}

func TestSerialWithContext(t *testing.T) {
	t.Run("serial with context", func(t *testing.T) {
		ctx := context.Background()
		order := []int{}

		err := SerialWithContext(ctx,
			func(ctx context.Context) error {
				order = append(order, 1)
				return nil
			},
			func(ctx context.Context) error {
				order = append(order, 2)
				return nil
			},
		)

		assert.NoError(t, err)
		assert.Equal(t, []int{1, 2}, order)
	})
}

func TestContextCancellation(t *testing.T) {
	t.Run("context cancelled propagates", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		g := New(ctx)

		g.Go(func() error {
			time.Sleep(500 * time.Millisecond)
			return nil
		})

		time.Sleep(50 * time.Millisecond)
		cancel()

		err := g.Wait()
		assert.NoError(t, err) // No explicit error, just context cancelled
	})
}

func TestConcurrency(t *testing.T) {
	t.Run("concurrent tasks", func(t *testing.T) {
		ctx := context.Background()
		g := New(ctx)

		for i := 0; i < 100; i++ {
			g.Go(func() error {
				time.Sleep(10 * time.Millisecond)
				return nil
			})
		}

		err := g.Wait()
		assert.NoError(t, err)
	})
}
