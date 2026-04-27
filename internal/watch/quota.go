package watch

import (
	"sync"
	"time"
)

// DefaultQuotaPolicy returns a sensible default quota policy.
func DefaultQuotaPolicy() QuotaPolicy {
	return QuotaPolicy{
		Max:    100,
		Window: time.Minute,
	}
}

// QuotaPolicy configures a Quota.
type QuotaPolicy struct {
	Max    int
	Window time.Duration
}

// Quota enforces a rolling usage limit over a fixed time window.
type Quota struct {
	mu      sync.Mutex
	policy  QuotaPolicy
	used    int
	resetsAt time.Time
	now     func() time.Time
}

// NewQuota creates a new Quota with the given policy.
func NewQuota(p QuotaPolicy) *Quota {
	if p.Max <= 0 {
		p.Max = DefaultQuotaPolicy().Max
	}
	if p.Window <= 0 {
		p.Window = DefaultQuotaPolicy().Window
	}
	return &Quota{
		policy:  p,
		now:     time.Now,
		resetsAt: time.Now().Add(p.Window),
	}
}

// Allow reports whether the quota permits one more unit of usage.
func (q *Quota) Allow() bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.maybeReset()
	if q.used >= q.policy.Max {
		return false
	}
	q.used++
	return true
}

// Remaining returns the number of allowed units left in the current window.
func (q *Quota) Remaining() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.maybeReset()
	r := q.policy.Max - q.used
	if r < 0 {
		return 0
	}
	return r
}

// Reset manually resets the quota counter and starts a new window.
func (q *Quota) Reset() {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.used = 0
	q.resetsAt = q.now().Add(q.policy.Window)
}

func (q *Quota) maybeReset() {
	if q.now().After(q.resetsAt) {
		q.used = 0
		q.resetsAt = q.now().Add(q.policy.Window)
	}
}
