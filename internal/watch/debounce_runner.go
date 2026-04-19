package watch

import (
	"context"
	"time"
)

// DebounceRunner wraps a function so that rapid successive invocations
// are collapsed into a single execution after the debounce wait period.
type DebounceRunner struct {
	debounce *Debounce
	fn       func(ctx context.Context) error
	resultCh chan error
}

// NewDebounceRunner creates a DebounceRunner.
// If d is nil a default Debounce is used.
func NewDebounceRunner(d *Debounce, fn func(ctx context.Context) error) *DebounceRunner {
	if d == nil {
		d = NewDebounce(DefaultDebouncePolicy())
	}
	if fn == nil {
		fn = func(_ context.Context) error { return nil }
	}
	return &DebounceRunner{
		debounce: d,
		fn:       fn,
		resultCh: make(chan error, 1),
	}
}

// Schedule enqueues a debounced call. The result of the eventual fn
// invocation is sent to the returned channel (buffered, size 1).
func (r *DebounceRunner) Schedule(ctx context.Context) <-chan error {
	ch := make(chan error, 1)
	r.debounce.Trigger(func() {
		ch <- r.fn(ctx)
	})
	return ch
}

// Wait blocks until the debounce fires or the context is cancelled.
func (r *DebounceRunner) Wait(ctx context.Context, wait time.Duration) error {
	ch := r.Schedule(ctx)
	select {
	case err := <-ch:
		return err
	case <-ctx.Done():
		r.debounce.Cancel()
		return ctx.Err()
	}
}
