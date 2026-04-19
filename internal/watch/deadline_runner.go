package watch

import (
	"context"
	"time"
)

// DeadlineRunner wraps a function and refuses to execute it after a
// deadline has passed. If the deadline is zero the function always runs.
type DeadlineRunner struct {
	deadline *Deadline
	fn       func(ctx context.Context) error
}

// NewDeadlineRunner returns a DeadlineRunner. If deadline is nil a zero
// Deadline (no expiry) is used. If fn is nil it defaults to a noop.
func NewDeadlineRunner(deadline *Deadline, fn func(ctx context.Context) error) *DeadlineRunner {
	if deadline == nil {
		deadline = NewDeadline(DefaultDeadlinePolicy())
	}
	if fn == nil {
		fn = func(ctx context.Context) error { return nil }
	}
	return &DeadlineRunner{deadline: deadline, fn: fn}
}

// Run executes fn unless the deadline has already expired.
// Returns ErrDeadlineExceeded when skipped due to expiry.
func (r *DeadlineRunner) Run(ctx context.Context) error {
	if r.deadline.Expired() {
		return ErrDeadlineExceeded
	}
	return r.fn(ctx)
}

// ErrDeadlineExceeded is returned when a DeadlineRunner skips execution.
var ErrDeadlineExceeded = deadlineExceededError("deadline exceeded")

type deadlineExceededError string

func (e deadlineExceededError) Error() string { return string(e) }

// RunUntil executes fn repeatedly on interval until the deadline fires or
// ctx is cancelled.
func (r *DeadlineRunner) RunUntil(ctx context.Context, interval time.Duration) error {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := r.Run(ctx); err != nil {
				return err
			}
		}
	}
}
