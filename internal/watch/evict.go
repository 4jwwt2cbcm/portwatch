package watch

import (
	"sync"
	"time"
)

// DefaultEvictPolicy returns sensible defaults for the eviction policy.
func DefaultEvictPolicy() EvictPolicy {
	return EvictPolicy{
		TTL:      30 * time.Second,
		Capacity: 256,
	}
}

// EvictPolicy controls eviction behaviour.
type EvictPolicy struct {
	TTL      time.Duration
	Capacity int
}

type evictEntry[T any] struct {
	value     T
	ExpiresAt time.Time
}

// EvictCache is a capacity-bounded TTL cache that evicts expired or
// least-recently-added entries when full.
type EvictCache[T any] struct {
	mu     sync.Mutex
	policy EvictPolicy
	items  map[string]evictEntry[T]
	order  []string
}

// NewEvictCache returns a new EvictCache with the given policy.
// Zero-value fields are replaced with defaults.
func NewEvictCache[T any](p EvictPolicy) *EvictCache[T] {
	if p.TTL <= 0 {
		p.TTL = DefaultEvictPolicy().TTL
	}
	if p.Capacity <= 0 {
		p.Capacity = DefaultEvictPolicy().Capacity
	}
	return &EvictCache[T]{
		policy: p,
		items:  make(map[string]evictEntry[T]),
	}
}

// Set stores a value under key, evicting expired or oldest entries as needed.
func (c *EvictCache[T]) Set(key string, value T) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.evictExpired()
	if _, exists := c.items[key]; !exists {
		if len(c.items) >= c.policy.Capacity {
			c.evictOldest()
		}
		c.order = append(c.order, key)
	}
	c.items[key] = evictEntry[T]{value: value, ExpiresAt: time.Now().Add(c.policy.TTL)}
}

// Get retrieves a value by key. Returns the zero value and false if not found or expired.
func (c *EvictCache[T]) Get(key string) (T, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	e, ok := c.items[key]
	if !ok || time.Now().After(e.ExpiresAt) {
		var zero T
		delete(c.items, key)
		return zero, false
	}
	return e.value, true
}

// Len returns the number of non-expired entries.
func (c *EvictCache[T]) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.evictExpired()
	return len(c.items)
}

func (c *EvictCache[T]) evictExpired() {
	now := time.Now()
	for k, e := range c.items {
		if now.After(e.ExpiresAt) {
			delete(c.items, k)
		}
	}
}

func (c *EvictCache[T]) evictOldest() {
	for len(c.order) > 0 {
		oldest := c.order[0]
		c.order = c.order[1:]
		if _, ok := c.items[oldest]; ok {
			delete(c.items, oldest)
			return
		}
	}
}
