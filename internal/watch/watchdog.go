package watch

import (
	"context"
	"sync"
	"time"
)

// DefaultWatchdogPolicy returns a WatchdogPolicy with sensible defaults.
func DefaultWatchdogPolicy() WatchdogPolicy {
	return WatchdogPolicy{
		Timeout:  30 * time.Second,
		Interval: 5 * time.Second,
	}
}

// WatchdogPolicy controls watchdog behaviour.
type WatchdogPolicy struct {
	Timeout  time.Duration // how long without a kick before firing
	Interval time.Duration // how often to check the deadline
}

// Watchdog fires a callback when it has not been kicked within the timeout.
type Watchdog struct {
	policy  WatchdogPolicy
	mu      sync.Mutex
	last    time.Time
	onFire  func()
}

// NewWatchdog creates a Watchdog with the given policy and fire callback.
// If policy fields are zero, defaults are applied. onFire must not be nil.
func NewWatchdog(policy WatchdogPolicy, onFire func()) *Watchdog {
	if policy.Timeout <= 0 {
		policy.Timeout = DefaultWatchdogPolicy().Timeout
	}
	if policy.Interval <= 0 {
		policy.Interval = DefaultWatchdogPolicy().Interval
	}
	if onFire == nil {
		onFire = func() {}
	}
	return &Watchdog{
		policy: policy,
		last:   time.Now(),
		onFire: onFire,
	}
}

// Kick resets the watchdog deadline.
func (w *Watchdog) Kick() {
	w.mu.Lock()
	w.last = time.Now()
	w.mu.Unlock()
}

// LastKick returns the time of the most recent kick.
func (w *Watchdog) LastKick() time.Time {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.last
}

// Run starts the watchdog loop, blocking until ctx is cancelled.
func (w *Watchdog) Run(ctx context.Context) error {
	ticker := time.NewTicker(w.policy.Interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			w.mu.Lock()
			expired := time.Since(w.last) > w.policy.Timeout
			w.mu.Unlock()
			if expired {
				w.onFire()
			}
		}
	}
}
