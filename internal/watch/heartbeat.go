package watch

import (
	"context"
	"sync"
	"time"
)

// HeartbeatPolicy configures heartbeat behaviour.
type HeartbeatPolicy struct {
	Interval time.Duration
	Timeout  time.Duration
}

// DefaultHeartbeatPolicy returns sensible defaults.
func DefaultHeartbeatPolicy() HeartbeatPolicy {
	return HeartbeatPolicy{
		Interval: 30 * time.Second,
		Timeout:  5 * time.Second,
	}
}

// Heartbeat tracks periodic liveness beats and exposes whether the
// monitored component is still considered alive.
type Heartbeat struct {
	policy HeartbeatPolicy
	mu     sync.Mutex
	last   time.Time
}

// NewHeartbeat creates a Heartbeat with the given policy.
// A zero-value policy falls back to DefaultHeartbeatPolicy.
func NewHeartbeat(p HeartbeatPolicy) *Heartbeat {
	if p.Interval == 0 {
		p = DefaultHeartbeatPolicy()
	}
	return &Heartbeat{policy: p}
}

// Beat records the current time as the latest heartbeat.
func (h *Heartbeat) Beat() {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.last = time.Now()
}

// Alive reports whether a beat was recorded within the configured timeout.
func (h *Heartbeat) Alive() bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.last.IsZero() {
		return false
	}
	return time.Since(h.last) <= h.policy.Timeout
}

// Last returns the time of the most recent beat (zero if none).
func (h *Heartbeat) Last() time.Time {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.last
}

// Run emits beats on the configured interval until ctx is cancelled.
func (h *Heartbeat) Run(ctx context.Context) error {
	ticker := time.NewTicker(h.policy.Interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			h.Beat()
		}
	}
}
