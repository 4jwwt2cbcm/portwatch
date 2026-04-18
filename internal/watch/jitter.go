package watch

import (
	"math/rand"
	"time"
)

// JitterPolicy configures jitter behaviour.
type JitterPolicy struct {
	// MaxFraction is the maximum fraction of the base duration to add as jitter.
	// E.g. 0.2 adds up to 20% of the base duration.
	MaxFraction float64
}

// DefaultJitterPolicy returns a JitterPolicy with sensible defaults.
func DefaultJitterPolicy() JitterPolicy {
	return JitterPolicy{MaxFraction: 0.25}
}

// Jitter adds a random duration up to MaxFraction of base to base.
type Jitter struct {
	policy JitterPolicy
	rng    *rand.Rand
}

// NewJitter creates a Jitter with the given policy.
func NewJitter(policy JitterPolicy) *Jitter {
	if policy.MaxFraction <= 0 {
		policy.MaxFraction = DefaultJitterPolicy().MaxFraction
	}
	return &Jitter{
		policy: policy,
		rng:    rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Apply returns base plus a random jitter up to MaxFraction of base.
func (j *Jitter) Apply(base time.Duration) time.Duration {
	if base <= 0 {
		return base
	}
	max := float64(base) * j.policy.MaxFraction
	offset := time.Duration(j.rng.Float64() * max)
	return base + offset
}

// Reset re-seeds the internal RNG.
func (j *Jitter) Reset() {
	j.rng = rand.New(rand.NewSource(time.Now().UnixNano()))
}
