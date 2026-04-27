package watch

import (
	"sync"
	"time"
)

// DefaultLeasePolicy returns a LeasePolicy with sensible defaults.
func DefaultLeasePolicy() LeasePolicy {
	return LeasePolicy{
		TTL:     30 * time.Second,
		RenewAt: 0.75,
	}
}

// LeasePolicy controls lease duration and renewal threshold.
type LeasePolicy struct {
	// TTL is how long a lease is valid after acquisition or renewal.
	TTL time.Duration
	// RenewAt is the fraction of TTL remaining at which renewal is suggested (0–1).
	RenewAt float64
}

// Lease represents an expiring, renewable ownership token.
type Lease struct {
	mu      sync.Mutex
	policy  LeasePolicy
	held    bool
	acquired time.Time
	expires  time.Time
	now      func() time.Time
}

// NewLease creates a Lease with the given policy.
// If policy.TTL is zero, DefaultLeasePolicy is used.
func NewLease(policy LeasePolicy) *Lease {
	if policy.TTL == 0 {
		policy = DefaultLeasePolicy()
	}
	if policy.RenewAt <= 0 || policy.RenewAt >= 1 {
		policy.RenewAt = 0.75
	}
	return &Lease{
		policy: policy,
		now:    time.Now,
	}
}

// Acquire attempts to take the lease. Returns false if already held and not expired.
func (l *Lease) Acquire() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := l.now()
	if l.held && now.Before(l.expires) {
		return false
	}
	l.held = true
	l.acquired = now
	l.expires = now.Add(l.policy.TTL)
	return true
}

// Renew extends the lease TTL from now. Returns false if lease is not held.
func (l *Lease) Renew() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := l.now()
	if !l.held || now.After(l.expires) {
		return false
	}
	l.expires = now.Add(l.policy.TTL)
	return true
}

// Release relinquishes the lease.
func (l *Lease) Release() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.held = false
}

// Held reports whether the lease is currently held and not expired.
func (l *Lease) Held() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.held && l.now().Before(l.expires)
}

// ShouldRenew reports whether the lease is held and past the renewal threshold.
func (l *Lease) ShouldRenew() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	if !l.held {
		return false
	}
	now := l.now()
	if now.After(l.expires) {
		return false
	}
	remaining := l.expires.Sub(now)
	threshold := time.Duration(float64(l.policy.TTL) * (1 - l.policy.RenewAt))
	return remaining <= threshold
}

// ExpiresAt returns the current expiry time of the lease.
func (l *Lease) ExpiresAt() time.Time {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.expires
}
