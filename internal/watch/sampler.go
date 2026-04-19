package watch

import (
	"sync"
	"time"
)

// SamplerPolicy configures sampling behaviour.
type SamplerPolicy struct {
	Rate     int           // max samples to keep per window
	Window   time.Duration // rolling window duration
}

// DefaultSamplerPolicy returns sensible defaults.
func DefaultSamplerPolicy() SamplerPolicy {
	return SamplerPolicy{
		Rate:   10,
		Window: time.Minute,
	}
}

// sample holds a recorded value and its timestamp.
type sample[T any] struct {
	value T
	at    time.Time
}

// Sampler records up to Rate values per Window, dropping older entries.
type Sampler[T any] struct {
	mu      sync.Mutex
	policy  SamplerPolicy
	samples []sample[T]
}

// NewSampler creates a Sampler with the given policy.
// Zero-value policy fields are replaced with defaults.
func NewSampler[T any](p SamplerPolicy) *Sampler[T] {
	def := DefaultSamplerPolicy()
	if p.Rate <= 0 {
		p.Rate = def.Rate
	}
	if p.Window <= 0 {
		p.Window = def.Window
	}
	return &Sampler[T]{policy: p}
}

// Record attempts to add v to the sampler.
// Returns true if the value was accepted, false if the rate limit was reached.
func (s *Sampler[T]) Record(v T) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	s.evict(now)

	if len(s.samples) >= s.policy.Rate {
		return false
	}
	s.samples = append(s.samples, sample[T]{value: v, at: now})
	return true
}

// Snapshot returns a copy of all samples currently in the window.
func (s *Sampler[T]) Snapshot() []T {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.evict(time.Now())
	out := make([]T, len(s.samples))
	for i, sm := range s.samples {
		out[i] = sm.value
	}
	return out
}

// Len returns the number of samples currently in the window.
func (s *Sampler[T]) Len() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.evict(time.Now())
	return len(s.samples)
}

// evict removes samples outside the rolling window. Must be called with mu held.
func (s *Sampler[T]) evict(now time.Time) {
	cutoff := now.Add(-s.policy.Window)
	i := 0
	for i < len(s.samples) && s.samples[i].at.Before(cutoff) {
		i++
	}
	s.samples = s.samples[i:]
}
