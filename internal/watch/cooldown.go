package watch

import (
	"sync"
	"time"
)

// CooldownPolicy defines the parameters for a cooldown tracker.
type CooldownPolicy struct {
	Duration time.Duration
}

// DefaultCooldownPolicy returns a policy with a 30-second cooldown.
func DefaultCooldownPolicy() CooldownPolicy {
	return CooldownPolicy{
		Duration: 30 * time.Second,
	}
}

// Cooldown tracks whether a named action is currently in a cooldown period.
type Cooldown struct {
	mu       sync.Mutex
	policy   CooldownPolicy
	expiries map[string]time.Time
	now      func() time.Time
}

// NewCooldown creates a new Cooldown with the given policy.
func NewCooldown(policy CooldownPolicy) *Cooldown {
	if policy.Duration <= 0 {
		policy.Duration = DefaultCooldownPolicy().Duration
	}
	return &Cooldown{
		policy:   policy,
		expiries: make(map[string]time.Time),
		now:      time.Now,
	}
}

// Allow returns true and records the cooldown if the key is not currently cooling down.
// Returns false if the key is still within its cooldown window.
func (c *Cooldown) Allow(key string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := c.now()
	if exp, ok := c.expiries[key]; ok && now.Before(exp) {
		return false
	}
	c.expiries[key] = now.Add(c.policy.Duration)
	return true
}

// Reset clears the cooldown for a specific key, allowing it immediately.
func (c *Cooldown) Reset(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.expiries, key)
}

// Active returns true if the key is currently in a cooldown period.
func (c *Cooldown) Active(key string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	exp, ok := c.expiries[key]
	return ok && c.now().Before(exp)
}
