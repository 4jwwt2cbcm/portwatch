package watch

import (
	"context"
)

// TimeoutRunner wraps an inner function with a Timeout, ensuring it
// completes within the configured deadline on each invocation.
type TimeoutRunner struct {
	timeout *Timeout
	fn      func(ctx context.Context) error
}

// NewTimeoutRunner creates a TimeoutRunner using the given policy and function.
// If policy duration is zero, DefaultTimeoutPolicy is used.
func NewTimeoutRunner(p TimeoutPolicy, fn func(ctx context.Context) error) *TimeoutRunner {
	return &TimeoutRunner{
		timeout: NewTimeout(p),
		fn:      fn,
	}
}

// Run executes the wrapped function with the configured timeout.
func (r *TimeoutRunner) Run(ctx context.Context) error {
	return r.timeout.Run(ctx, r.fn)
}
