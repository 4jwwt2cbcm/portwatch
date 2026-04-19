package watch

import (
	"sync"
	"time"
)

// RotatorPolicy configures the log rotator.
type RotatorPolicy struct {
	MaxEntries int
	MaxAge     time.Duration
}

// DefaultRotatorPolicy returns sensible defaults.
func DefaultRotatorPolicy() RotatorPolicy {
	return RotatorPolicy{
		MaxEntries: 1000,
		MaxAge:     24 * time.Hour,
	}
}

// RotatorEntry holds a timestamped value.
type RotatorEntry[T any] struct {
	Value     T
	Timestamp time.Time
}

// Rotator is a fixed-capacity ring buffer that evicts entries by age or count.
type Rotator[T any] struct {
	mu      sync.Mutex
	policy  RotatorPolicy
	entries []RotatorEntry[T]
	now     func() time.Time
}

// NewRotator creates a Rotator with the given policy.
func NewRotator[T any](policy RotatorPolicy) *Rotator[T] {
	if policy.MaxEntries <= 0 {
		policy.MaxEntries = DefaultRotatorPolicy().MaxEntries
	}
	if policy.MaxAge <= 0 {
		policy.MaxAge = DefaultRotatorPolicy().MaxAge
	}
	return &Rotator[T]{
		policy: policy,
		now:    time.Now,
	}
}

// Add appends a value, evicting stale or excess entries as needed.
func (r *Rotator[T]) Add(v T) {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := r.now()
	cutoff := now.Add(-r.policy.MaxAge)
	// evict by age
	start := 0
	for start < len(r.entries) && r.entries[start].Timestamp.Before(cutoff) {
		start++
	}
	r.entries = r.entries[start:]
	// evict by count
	if len(r.entries) >= r.policy.MaxEntries {
		r.entries = r.entries[len(r.entries)-r.policy.MaxEntries+1:]
	}
	r.entries = append(r.entries, RotatorEntry[T]{Value: v, Timestamp: now})
}

// Snapshot returns a copy of current entries.
func (r *Rotator[T]) Snapshot() []RotatorEntry[T] {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]RotatorEntry[T], len(r.entries))
	copy(out, r.entries)
	return out
}

// Len returns the number of entries.
func (r *Rotator[T]) Len() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.entries)
}

// Clear removes all entries.
func (r *Rotator[T]) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.entries = nil
}
