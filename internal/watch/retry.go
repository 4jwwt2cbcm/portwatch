package watch

import (
	"context"
	"errors"
	"time"
)

// RetryPolicy defines how retries are attempted.
type RetryPolicy struct {
	MaxAttempts int
	Backoff     *Backoff
}

// DefaultRetryPolicy returns a sensible default retry policy.
func DefaultRetryPolicy() RetryPolicy {
	return RetryPolicy{
		MaxAttempts: 5,
		Backoff:     NewBackoff(DefaultBackoffPolicy()),
	}
}

// Retrier executes a function with retry logic.
type Retrier struct {
	policy RetryPolicy
	logger *Logger
}

// NewRetrier creates a Retrier with the given policy and logger.
func NewRetrier(policy RetryPolicy, logger *Logger) *Retrier {
	return &Retrier{policy: policy, logger: logger}
}

// Run executes fn, retrying on error up to MaxAttempts times.
// Returns nil on success, or the last error if all attempts fail.
// Respects context cancellation between attempts.
func (r *Retrier) Run(ctx context.Context, fn func(ctx context.Context) error) error {
	var lastErr error
	for attempt := 1; attempt <= r.policy.MaxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}
		lastErr = fn(ctx)
		if lastErr == nil {
			r.policy.Backoff.Reset()
			return nil
		}
		if errors.Is(lastErr, context.Canceled) || errors.Is(lastErr, context.DeadlineExceeded) {
			return lastErr
		}
		wait := r.policy.Backoff.Next()
		if r.logger != nil {
			r.logger.Error("retrier: attempt failed, retrying", lastErr)
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(wait):
		}
	}
	return lastErr
}
