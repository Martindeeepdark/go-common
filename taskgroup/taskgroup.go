package taskgroup

import (
	"context"
	"sync"
)

// TaskGroup manages a group of goroutines with error handling and context cancellation
type TaskGroup struct {
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
	mu     sync.Mutex
	errors []error
}

// New creates a new task group with a context
func New(ctx context.Context) *TaskGroup {
	ctx, cancel := context.WithCancel(ctx)
	return &TaskGroup{
		ctx:    ctx,
		cancel: cancel,
		errors: make([]error, 0),
	}
}

// Go runs a function in a goroutine
func (g *TaskGroup) Go(fn func() error) {
	g.wg.Add(1)

	go func() {
		defer g.wg.Done()

		if err := fn(); err != nil {
			g.mu.Lock()
			g.errors = append(g.errors, err)
			g.mu.Unlock()

			// Cancel context on error
			g.cancel()
		}
	}()
}

// GoWithContext runs a function with context in a goroutine
func (g *TaskGroup) GoWithContext(fn func(context.Context) error) {
	g.wg.Add(1)

	go func() {
		defer g.wg.Done()

		if err := fn(g.ctx); err != nil {
			g.mu.Lock()
			g.errors = append(g.errors, err)
			g.mu.Unlock()

			// Cancel context on error
			g.cancel()
		}
	}()
}

// Wait waits for all goroutines to complete and returns any errors
func (g *TaskGroup) Wait() error {
	g.wg.Wait()

	if len(g.errors) > 0 {
		return g.errors[0]
	}

	return nil
}

// WaitAll waits for all goroutines and returns all errors
func (g *TaskGroup) WaitAll() []error {
	g.wg.Wait()
	return g.errors
}

// Cancel cancels the task group context
func (g *TaskGroup) Cancel() {
	g.cancel()
}

// Context returns the task group's context
func (g *TaskGroup) Context() context.Context {
	return g.ctx
}

// Errors returns all collected errors
func (g *TaskGroup) Errors() []error {
	g.mu.Lock()
	defer g.mu.Unlock()

	errors := make([]error, len(g.errors))
	copy(errors, g.errors)
	return errors
}

// HasError returns true if any error occurred
func (g *TaskGroup) HasError() bool {
	g.mu.Lock()
	defer g.mu.Unlock()

	return len(g.errors) > 0
}

// Parallel runs multiple functions in parallel and waits for completion
func Parallel(ctx context.Context, fns ...func() error) error {
	g := New(ctx)
	for _, fn := range fns {
		g.Go(fn)
	}
	return g.Wait()
}

// ParallelWithContext runs multiple functions with context in parallel
func ParallelWithContext(ctx context.Context, fns ...func(context.Context) error) error {
	g := New(ctx)
	for _, fn := range fns {
		g.GoWithContext(fn)
	}
	return g.Wait()
}

// Serial runs functions serially, stopping on first error
func Serial(ctx context.Context, fns ...func() error) error {
	for _, fn := range fns {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err := fn(); err != nil {
				return err
			}
		}
	}
	return nil
}

// SerialWithContext runs functions with context serially
func SerialWithContext(ctx context.Context, fns ...func(context.Context) error) error {
	for _, fn := range fns {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err := fn(ctx); err != nil {
				return err
			}
		}
	}
	return nil
}
