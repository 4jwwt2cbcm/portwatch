package watch

import (
	"sync"
	"time"
)

// DebouncePolicy holds configuration for a Debounce.
type DebouncePolicy struct {
	Wait time.Duration
}

// DefaultDebouncePolicy returns sensible defaults.
func DefaultDebouncePolicy() DebouncePolicy {
	return DebouncePolicy{
		Wait: 500 * time.Millisecond,
	}
}

// Debounce delays execution of a function until a quiet period has elapsed
// since the last call. Concurrent calls within the wait window are collapsed.
type Debounce struct {
	policy  DebouncePolicy
	mu      sync.Mutex
	timer   *time.Timer
	pending bool
}

// NewDebounce creates a Debounce with the given policy.
// Zero-value Wait is replaced with the default.
func NewDebounce(p DebouncePolicy) *Debounce {
	if p.Wait <= 0 {
		p.Wait = DefaultDebouncePolicy().Wait
	}
	return &Debounce{policy: p}
}

// Trigger schedules fn to run after the wait period.
// If Trigger is called again before the timer fires, the timer resets.
func (d *Debounce) Trigger(fn func()) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.timer != nil {
		d.timer.Stop()
	}
	d.pending = true
	d.timer = time.AfterFunc(d.policy.Wait, func() {
		d.mu.Lock()
		d.pending = false
		d.mu.Unlock()
		fn()
	})
}

// Pending reports whether a call is waiting to fire.
func (d *Debounce) Pending() bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.pending
}

// Cancel stops any pending call.
func (d *Debounce) Cancel() {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.timer != nil {
		d.timer.Stop()
		d.pending = false
	}
}
