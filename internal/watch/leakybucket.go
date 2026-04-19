package watch

import (
	"sync"
	"time"
)

// DefaultLeakyBucketPolicy returns a policy with sensible defaults.
func DefaultLeakyBucketPolicy() LeakyBucketPolicy {
	return LeakyBucketPolicy{
		Capacity:  10,
		LeakRate:  1, // units per second
		LeakEvery: time.Second,
	}
}

// LeakyBucketPolicy configures a LeakyBucket.
type LeakyBucketPolicy struct {
	Capacity  int
	LeakRate  int
	LeakEvery time.Duration
}

// LeakyBucket is a concurrency-safe leaky bucket rate limiter.
type LeakyBucket struct {
	policy LeakyBucketPolicy
	level  int
	lastAt time.Time
	mu     sync.Mutex
}

// NewLeakyBucket creates a LeakyBucket with the given policy.
func NewLeakyBucket(p LeakyBucketPolicy) *LeakyBucket {
	if p.Capacity <= 0 {
		p.Capacity = DefaultLeakyBucketPolicy().Capacity
	}
	if p.LeakRate <= 0 {
		p.LeakRate = DefaultLeakyBucketPolicy().LeakRate
	}
	if p.LeakEvery <= 0 {
		p.LeakEvery = DefaultLeakyBucketPolicy().LeakEvery
	}
	return &LeakyBucket{policy: p, lastAt: time.Now()}
}

// Allow returns true if there is capacity in the bucket after leaking.
func (b *LeakyBucket) Allow() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.leak()
	if b.level < b.policy.Capacity {
		b.level++
		return true
	}
	return false
}

// Level returns the current fill level.
func (b *LeakyBucket) Level() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.leak()
	return b.level
}

// Reset drains the bucket.
func (b *LeakyBucket) Reset() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.level = 0
	b.lastAt = time.Now()
}

func (b *LeakyBucket) leak() {
	now := time.Now()
	elapsed := now.Sub(b.lastAt)
	periods := int(elapsed / b.policy.LeakEvery)
	if periods > 0 {
		b.level -= periods * b.policy.LeakRate
		if b.level < 0 {
			b.level = 0
		}
		b.lastAt = b.lastAt.Add(time.Duration(periods) * b.policy.LeakEvery)
	}
}
