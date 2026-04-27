package watch

import (
	"context"
	"errors"
)

// ErrQuotaExceeded is returned when the quota has been exhausted.
var ErrQuotaExceeded = errors.New("quota exceeded")

// QuotaRunner wraps a function and enforces a Quota before each invocation.
type QuotaRunner struct {
	quota *Quota
	fn    func(ctx context.Context) error
}

// NewQuotaRunner creates a QuotaRunner. If quota is nil a default quota is used.
func NewQuotaRunner(q *Quota, fn func(ctx context.Context) error) *QuotaRunner {
	if q == nil {
		q = NewQuota(DefaultQuotaPolicy())
	}
	if fn == nil {
		fn = func(_ context.Context) error { return nil }
	}
	return &QuotaRunner{quota: q, fn: fn}
}

// Run executes the wrapped function if the quota permits, otherwise returns
// ErrQuotaExceeded.
func (r *QuotaRunner) Run(ctx context.Context) error {
	if !r.quota.Allow() {
		return ErrQuotaExceeded
	}
	return r.fn(ctx)
}
