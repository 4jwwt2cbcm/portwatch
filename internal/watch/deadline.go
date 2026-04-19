package watch

import (
	"context"
	"sync"
	"time"
)

// DeadlinePolicy configures Deadline behaviour.
type DeadlinePolicy struct {
	At time.Time
}

// DefaultDeadlinePolicy returns a policy with no deadline set.
func DefaultDeadlinePolicy() DeadlinePolicy {
	return DeadlinePolicy{}
}

// Deadline tracks an absolute point in time after which it is considered
// expired. It is safe for concurrent use.
type Deadline struct {
	mu  sync.RWMutex
	at  time.Time
}

// NewDeadline constructs a Deadline from the given policy.
func NewDeadline(p DeadlinePolicy) *Deadline {
	return &Deadline{at: p.At}
}

// Set updates the deadline to t.
func (d *Deadline) Set(t time.Time) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.at = t
}

// At returns the current deadline.
func (d *Deadline) At() time.Time {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.at
}

// Expired reports whether the deadline has passed relative to now.
func (d *Deadline) Expired() bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	if d.at.IsZero() {
		return false
	}
	return time.Now().After(d.at)
}

// Wait blocks until the deadline is reached or ctx is cancelled.
// Returns context.DeadlineExceeded when the deadline fires, or ctx.Err().
func (d *Deadline) Wait(ctx context.Context) error {
	d.mu.RLock()
	at := d.at
	d.mu.RUnlock()
	if at.IsZero() {
		<-ctx.Done()
		return ctx.Err()
	}
	timer := time.NewTimer(time.Until(at))
	defer timer.Stop()
	select {
	case <-timer.C:
		return context.DeadlineExceeded
	case <-ctx.Done():
		return ctx.Err()
	}
}
