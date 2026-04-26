package watch

import (
	"sync"
	"time"
)

// FencePolicy configures the behaviour of a Fence.
type FencePolicy struct {
	// MaxCrossings is the maximum number of times the fence may be crossed
	// within the window before it blocks. Zero defaults to 1.
	MaxCrossings int
	// Window is the duration of the observation window. Zero defaults to 1s.
	Window time.Duration
}

// DefaultFencePolicy returns a FencePolicy with sensible defaults.
func DefaultFencePolicy() FencePolicy {
	return FencePolicy{
		MaxCrossings: 1,
		Window:       time.Second,
	}
}

// Fence is a concurrency primitive that limits how many times a boundary may
// be crossed within a sliding time window. Unlike a rate-limiter it is
// stateful per named boundary, making it useful for guarding re-entrant
// transitions in a state machine or event loop.
type Fence struct {
	mu       sync.Mutex
	policy   FencePolicy
	crossings map[string][]time.Time
}

// NewFence creates a Fence with the given policy. Zero-value fields are
// replaced with their defaults.
func NewFence(p FencePolicy) *Fence {
	if p.MaxCrossings <= 0 {
		p.MaxCrossings = DefaultFencePolicy().MaxCrossings
	}
	if p.Window <= 0 {
		p.Window = DefaultFencePolicy().Window
	}
	return &Fence{
		policy:    p,
		crossings: make(map[string][]time.Time),
	}
}

// Cross attempts to cross the named boundary. It returns true when the
// crossing is permitted (i.e. the number of crossings within the current
// window has not yet reached MaxCrossings). The crossing is recorded
// regardless of the return value so that callers can observe the full
// pressure on the fence.
func (f *Fence) Cross(name string) bool {
	f.mu.Lock()
	defer f.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-f.policy.Window)

	// Evict expired entries.
	prev := f.crossings[name]
	valid := prev[:0]
	for _, t := range prev {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}

	allowed := len(valid) < f.policy.MaxCrossings
	valid = append(valid, now)
	f.crossings[name] = valid
	return allowed
}

// Count returns the number of crossings recorded for the named boundary
// within the current window.
func (f *Fence) Count(name string) int {
	f.mu.Lock()
	defer f.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-f.policy.Window)
	count := 0
	for _, t := range f.crossings[name] {
		if t.After(cutoff) {
			count++
		}
	}
	return count
}

// Reset clears all recorded crossings for the named boundary.
func (f *Fence) Reset(name string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.crossings, name)
}

// ResetAll clears recorded crossings for every boundary.
func (f *Fence) ResetAll() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.crossings = make(map[string][]time.Time)
}
