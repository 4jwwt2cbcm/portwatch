package watch

import (
	"context"
	"errors"
)

// ErrNoToken is returned when no token is available in the pool.
var ErrNoToken = errors.New("no token available")

// TokenRunner wraps a function and only executes it when a token
// can be taken from the pool. If the pool is empty the call is
// rejected with ErrNoToken.
type TokenRunner struct {
	pool *TokenPool
	fn   func(ctx context.Context) error
}

// NewTokenRunner creates a TokenRunner. If pool is nil a default
// pool (capacity 10, 1 token/s refill) is used. If fn is nil a
// no-op function is used.
func NewTokenRunner(pool *TokenPool, fn func(ctx context.Context) error) *TokenRunner {
	if pool == nil {
		pool = NewTokenPool(DefaultTokenPolicy())
	}
	if fn == nil {
		fn = func(_ context.Context) error { return nil }
	}
	return &TokenRunner{pool: pool, fn: fn}
}

// Run attempts to take a token and, if successful, calls the
// wrapped function. Returns ErrNoToken if the pool is exhausted.
func (r *TokenRunner) Run(ctx context.Context) error {
	if !r.pool.Take() {
		return ErrNoToken
	}
	return r.fn(ctx)
}
