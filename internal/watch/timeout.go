package watch

import (
	"context"
	"errors"
	"time"
)

// ErrTimeout is returned when an operation exceeds its deadline.
var ErrTimeout = errors.New("operation timed out")

// TimeoutPolicy holds configuration for timeout behaviour.
type TimeoutPolicy struct {
	Duration time.Duration
}

// DefaultTimeoutPolicy returns a sensible default timeout policy.
func DefaultTimeoutPolicy() TimeoutPolicy {
	return TimeoutPolicy{
		Duration: 30 * time.Second,
	}
}

// Timeout wraps a function call with a deadline derived from the policy.
type Timeout struct {
	policy TimeoutPolicy
}

// NewTimeout creates a Timeout with the given policy.
func NewTimeout(p TimeoutPolicy) *Timeout {
	if p.Duration <= 0 {
		p = DefaultTimeoutPolicy()
	}
	return &Timeout{policy: p}
}

// Run executes fn within the configured deadline.
// Returns ErrTimeout if the deadline is exceeded.
func (t *Timeout) Run(ctx context.Context, fn func(ctx context.Context) error) error {
	ctx, cancel := context.WithTimeout(ctx, t.policy.Duration)
	defer cancel()

	ch := make(chan error, 1)
	go func() {
		ch <- fn(ctx)
	}()

	select {
	case err := <-ch:
		return err
	case <-ctx.Done():
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return ErrTimeout
		}
		return ctx.Err()
	}
}
