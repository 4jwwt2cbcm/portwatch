package watch

import (
	"sync"
	"time"
)

// TokenBucketPolicy configures a token bucket.
type TokenBucketPolicy struct {
	Capacity  int
	RefillRate int // tokens per second
}

// DefaultTokenBucketPolicy returns sensible defaults.
func DefaultTokenBucketPolicy() TokenBucketPolicy {
	return TokenBucketPolicy{
		Capacity:   10,
		RefillRate: 1,
	}
}

// TokenBucket is a thread-safe token bucket for rate limiting.
type TokenBucket struct {
	policy    TokenBucketPolicy
	tokens    int
	lastRefil time.Time
	mu        sync.Mutex
	now       func() time.Time
}

// NewTokenBucket creates a TokenBucket with the given policy.
func NewTokenBucket(p TokenBucketPolicy) *TokenBucket {
	if p.Capacity <= 0 {
		p.Capacity = DefaultTokenBucketPolicy().Capacity
	}
	if p.RefillRate <= 0 {
		p.RefillRate = DefaultTokenBucketPolicy().RefillRate
	}
	return &TokenBucket{
		policy:    p,
		tokens:    p.Capacity,
		lastRefil: time.Now(),
		now:       time.Now,
	}
}

// refill adds tokens based on elapsed time.
func (tb *TokenBucket) refill() {
	now := tb.now()
	elapsed := now.Sub(tb.lastRefil).Seconds()
	add := int(elapsed * float64(tb.policy.RefillRate))
	if add > 0 {
		tb.tokens += add
		if tb.tokens > tb.policy.Capacity {
			tb.tokens = tb.policy.Capacity
		}
		tb.lastRefil = now
	}
}

// Allow returns true and consumes a token if one is available.
func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	tb.refill()
	if tb.tokens > 0 {
		tb.tokens--
		return true
	}
	return false
}

// Tokens returns the current token count.
func (tb *TokenBucket) Tokens() int {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	tb.refill()
	return tb.tokens
}
