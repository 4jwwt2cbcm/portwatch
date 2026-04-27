package watch

import (
	"sync"
	"time"
)

// TickerPoolPolicy controls the behaviour of a TickerPool.
type TickerPoolPolicy struct {
	// MaxTickers is the maximum number of tickers the pool will manage.
	// Defaults to 16 when zero.
	MaxTickers int
	// DefaultInterval is used when a caller does not supply an interval.
	// Defaults to 1 second when zero.
	DefaultInterval time.Duration
}

// DefaultTickerPoolPolicy returns a sensible default policy.
func DefaultTickerPoolPolicy() TickerPoolPolicy {
	return TickerPoolPolicy{
		MaxTickers:      16,
		DefaultInterval: time.Second,
	}
}

// TickerPool manages a set of named tickers that share a common lifecycle.
type TickerPool struct {
	mu      sync.Mutex
	policy  TickerPoolPolicy
	tickers map[string]*time.Ticker
}

// NewTickerPool creates a new TickerPool with the given policy.
// Zero-value fields in the policy are replaced with defaults.
func NewTickerPool(policy TickerPoolPolicy) *TickerPool {
	if policy.MaxTickers <= 0 {
		policy.MaxTickers = DefaultTickerPoolPolicy().MaxTickers
	}
	if policy.DefaultInterval <= 0 {
		policy.DefaultInterval = DefaultTickerPoolPolicy().DefaultInterval
	}
	return &TickerPool{
		policy:  policy,
		tickers: make(map[string]*time.Ticker),
	}
}

// Add registers a named ticker with the given interval. If interval is zero
// the pool's DefaultInterval is used. Returns false if the pool is full or
// the name is already registered.
func (p *TickerPool) Add(name string, interval time.Duration) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	if _, exists := p.tickers[name]; exists {
		return false
	}
	if len(p.tickers) >= p.policy.MaxTickers {
		return false
	}
	if interval <= 0 {
		interval = p.policy.DefaultInterval
	}
	p.tickers[name] = time.NewTicker(interval)
	return true
}

// C returns the channel for a named ticker. Returns nil if not found.
func (p *TickerPool) C(name string) <-chan time.Time {
	p.mu.Lock()
	defer p.mu.Unlock()
	if t, ok := p.tickers[name]; ok {
		return t.C
	}
	return nil
}

// Remove stops and removes a named ticker.
func (p *TickerPool) Remove(name string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if t, ok := p.tickers[name]; ok {
		t.Stop()
		delete(p.tickers, name)
	}
}

// StopAll stops and removes every ticker in the pool.
func (p *TickerPool) StopAll() {
	p.mu.Lock()
	defer p.mu.Unlock()
	for name, t := range p.tickers {
		t.Stop()
		delete(p.tickers, name)
	}
}

// Len returns the number of active tickers.
func (p *TickerPool) Len() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return len(p.tickers)
}
