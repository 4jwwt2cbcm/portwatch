package watch

import (
	"sync"
	"time"
)

// RateLimiter enforces a minimum interval between allowed events using a
// token-bucket-style approach with a fixed refill period.
type RateLimiter struct {
	mu       sync.Mutex
	rate     int
	period   time.Duration
	tokens   int
	lastFill time.Time
	now      func() time.Time
}

// RateLimitPolicy holds configuration for a RateLimiter.
type RateLimitPolicy struct {
	Rate   int
	Period time.Duration
}

// DefaultRateLimitPolicy returns a sensible default: 10 events per minute.
func DefaultRateLimitPolicy() RateLimitPolicy {
	return RateLimitPolicy{
		Rate:   10,
		Period: time.Minute,
	}
}

// NewRateLimiter creates a RateLimiter from the given policy.
func NewRateLimiter(p RateLimitPolicy) *RateLimiter {
	now := time.Now()
	return &RateLimiter{
		rate:     p.Rate,
		period:   p.Period,
		tokens:   p.Rate,
		lastFill: now,
		now:      time.Now,
	}
}

// Allow returns true if the event is permitted under the current rate limit.
// It refills tokens proportionally to elapsed time since the last fill.
func (r *RateLimiter) Allow() bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := r.now()
	elapsed := now.Sub(r.lastFill)

	if elapsed >= r.period {
		periods := int(elapsed / r.period)
		r.tokens += periods * r.rate
		if r.tokens > r.rate {
			r.tokens = r.rate
		}
		r.lastFill = r.lastFill.Add(time.Duration(periods) * r.period)
	}

	if r.tokens <= 0 {
		return false
	}
	r.tokens--
	return true
}

// Reset restores the token bucket to its full capacity.
func (r *RateLimiter) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tokens = r.rate
	r.lastFill = r.now()
}
