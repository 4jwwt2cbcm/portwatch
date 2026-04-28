package watch

import (
	"sync"
	"time"
)

// DefaultHoldOffPolicy returns a HoldOffPolicy with sensible defaults.
func DefaultHoldOffPolicy() HoldOffPolicy {
	return HoldOffPolicy{
		Duration: 5 * time.Second,
	}
}

// HoldOffPolicy configures a HoldOff.
type HoldOffPolicy struct {
	// Duration is the minimum quiet period before the hold-off clears.
	Duration time.Duration
}

// HoldOff suppresses action until a quiet period has elapsed since the
// last activity signal. Useful for debouncing bursts where you want to
// wait until things settle before proceeding.
type HoldOff struct {
	mu       sync.Mutex
	policy   HoldOffPolicy
	lastSeen time.Time
	now      func() time.Time
}

// NewHoldOff creates a HoldOff with the given policy.
// If policy.Duration is zero, DefaultHoldOffPolicy is used.
func NewHoldOff(policy HoldOffPolicy) *HoldOff {
	if policy.Duration <= 0 {
		policy = DefaultHoldOffPolicy()
	}
	return &HoldOff{
		policy: policy,
		now:    time.Now,
	}
}

// Signal records activity, resetting the quiet period.
func (h *HoldOff) Signal() {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.lastSeen = h.now()
}

// Clear returns true if no activity has been signalled within the hold-off
// duration, meaning the system has been quiet long enough to proceed.
func (h *HoldOff) Clear() bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.lastSeen.IsZero() {
		return true
	}
	return h.now().Sub(h.lastSeen) >= h.policy.Duration
}

// Reset clears the last-seen timestamp, making Clear return true immediately.
func (h *HoldOff) Reset() {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.lastSeen = time.Time{}
}

// LastSeen returns the time of the most recent Signal call, or zero if
// Signal has never been called or Reset was called.
func (h *HoldOff) LastSeen() time.Time {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.lastSeen
}
