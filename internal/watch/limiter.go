package watch

import (
	"context"
	"sync"
	"time"
)

// LimiterPolicy configures the sliding window limiter.
type LimiterPolicy struct {
	MaxCalls int
	Window   time.Duration
}

// DefaultLimiterPolicy returns sensible defaults.
func DefaultLimiterPolicy() LimiterPolicy {
	return LimiterPolicy{
		MaxCalls: 10,
		Window:   time.Minute,
	}
}

// Limiter enforces a sliding window call limit.
type Limiter struct {
	policy LimiterPolicy
	mu     sync.Mutex
	calls  []time.Time
	now    func() time.Time
}

// NewLimiter creates a Limiter with the given policy.
func NewLimiter(p LimiterPolicy) *Limiter {
	if p.MaxCalls <= 0 {
		p.MaxCalls = DefaultLimiterPolicy().MaxCalls
	}
	if p.Window <= 0 {
		p.Window = DefaultLimiterPolicy().Window
	}
	return &Limiter{policy: p, now: time.Now}
}

// Allow returns true if the call is within the window limit.
func (l *Limiter) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := l.now()
	cutoff := now.Add(-l.policy.Window)
	valid := l.calls[:0]
	for _, t := range l.calls {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}
	l.calls = valid
	if len(l.calls) >= l.policy.MaxCalls {
		return false
	}
	l.calls = append(l.calls, now)
	return true
}

// Wait blocks until a call is allowed or ctx is done.
func (l *Limiter) Wait(ctx context.Context) error {
	for {
		if l.Allow() {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(50 * time.Millisecond):
		}
	}
}

// Reset clears all recorded calls.
func (l *Limiter) Reset() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.calls = nil
}
