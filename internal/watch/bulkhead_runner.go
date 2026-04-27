package watch

import "context"

// BulkheadRunner wraps a Bulkhead and executes a function within its
// concurrency and queue constraints.
type BulkheadRunner struct {
	bulkhead *Bulkhead
	fn       func(context.Context) error
}

// NewBulkheadRunner creates a BulkheadRunner. If bulkhead is nil, a default
// Bulkhead is created. If fn is nil, a no-op is used.
func NewBulkheadRunner(bulkhead *Bulkhead, fn func(context.Context) error) *BulkheadRunner {
	if bulkhead == nil {
		bulkhead = NewBulkhead(DefaultBulkheadPolicy())
	}
	if fn == nil {
		fn = func(context.Context) error { return nil }
	}
	return &BulkheadRunner{bulkhead: bulkhead, fn: fn}
}

// Run executes the wrapped function within the bulkhead. It returns
// ErrBulkheadFull if there is no capacity and no queue space, or the
// underlying function's error otherwise.
func (r *BulkheadRunner) Run(ctx context.Context) error {
	return r.bulkhead.Do(ctx, func() error {
		return r.fn(ctx)
	})
}
