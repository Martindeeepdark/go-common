package sql

import (
	"context"
	"errors"
	"fmt"

	"github.com/Martindeeepdark/go-common/database/defs"
)

var (
	// ErrTxNotFound is returned when transaction is not found in context
	ErrTxNotFound = errors.New("transaction not found in context")
)

// Manager manages database transactions
type Manager struct{}

// TxManager is the global transaction manager
var TxManager = &Manager{}

// txKey is the context key for storing transaction
type ctxKey struct{}

var txKeyVal = ctxKey{}

// WithContext adds transaction to context
func (m *Manager) WithContext(ctx context.Context, tx defs.Transaction) context.Context {
	return context.WithValue(ctx, txKeyVal, tx)
}

// FromContext retrieves transaction from context
func (m *Manager) FromContext(ctx context.Context) (defs.Transaction, bool) {
	tx, ok := ctx.Value(txKeyVal).(defs.Transaction)
	return tx, ok
}

// MustFromContext retrieves transaction from context or panics
func (m *Manager) MustFromContext(ctx context.Context) defs.Transaction {
	tx, ok := m.FromContext(ctx)
	if !ok {
		panic(ErrTxNotFound)
	}
	return tx
}

// ExecuteInTx executes a function within a transaction
// If the function returns an error, the transaction is rolled back
// Otherwise, the transaction is committed
func (m *Manager) ExecuteInTx(ctx context.Context, db defs.SQLDatabase, fn func(context.Context) error) error {
	// Begin transaction
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Add transaction to context
	ctx = m.WithContext(ctx, tx)

	// Execute function
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p) // re-throw panic after rollback
		}
	}()

	if err := fn(ctx); err != nil {
		// Rollback on error
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("failed to rollback: %v (original error: %w)", rbErr, err)
		}
		return err
	}

	// Commit on success
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// ExecuteInTxWithOpts executes a function within a transaction with options
func (m *Manager) ExecuteInTxWithOpts(ctx context.Context, db defs.SQLDatabase, opts *defs.TxOptions, fn func(context.Context) error) error {
	tx, err := db.BeginTx(ctx, opts)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	ctx = m.WithContext(ctx, tx)

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(ctx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("failed to rollback: %v (original error: %w)", rbErr, err)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// IsInTx checks if there's an active transaction in context
func (m *Manager) IsInTx(ctx context.Context) bool {
	_, ok := m.FromContext(ctx)
	return ok
}
