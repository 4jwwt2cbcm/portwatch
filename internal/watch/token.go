package watch

import (
	"context"
	"sync"
	"time"
)

// TokenPolicy defines configuration for the token pool.
type TokenPolicy struct {
	Capacity    int
	RefillEvery time.Duration
	RefillAmount int
}

// DefaultTokenPolicy returns sensible defaults.
func DefaultTokenPolicy() TokenPolicy {
	return TokenPolicy{
		Capacity:     10,
		RefillEvery:  time.Second,
		RefillAmount: 1,
	}
}

// TokenPool is a fixed-capacity pool of tokens that refills on a schedule.
type TokenPool struct {
	mu      sync.Mutex
	policy  TokenPolicy
	avail   int
	stopCh  chan struct{}
}

// NewTokenPool creates a TokenPool with the given policy.
func NewTokenPool(p TokenPolicy) *TokenPool {
	if p.Capacity <= 0 {
		p.Capacity = DefaultTokenPolicy().Capacity
	}
	if p.RefillEvery <= 0 {
		p.RefillEvery = DefaultTokenPolicy().RefillEvery
	}
	if p.RefillAmount <= 0 {
		p.RefillAmount = DefaultTokenPolicy().RefillAmount
	}
	return &TokenPool{
		policy: p,
		avail:  p.Capacity,
		stopCh: make(chan struct{}),
	}
}

// Take attempts to consume one token. Returns false if none are available.
func (tp *TokenPool) Take() bool {
	tp.mu.Lock()
	defer tp.mu.Unlock()
	if tp.avail <= 0 {
		return false
	}
	tp.avail--
	return true
}

// Available returns the current number of available tokens.
func (tp *TokenPool) Available() int {
	tp.mu.Lock()
	defer tp.mu.Unlock()
	return tp.avail
}

// Run starts the refill loop, blocking until ctx is cancelled.
func (tp *TokenPool) Run(ctx context.Context) error {
	ticker := time.NewTicker(tp.policy.RefillEvery)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			tp.mu.Lock()
			tp.avail += tp.policy.RefillAmount
			if tp.avail > tp.policy.Capacity {
				tp.avail = tp.policy.Capacity
			}
			tp.mu.Unlock()
		}
	}
}
