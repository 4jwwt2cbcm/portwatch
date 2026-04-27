package watch

import (
	"sync"
	"time"
)

// DefaultSurgePolicy returns a SurgePolicy with sensible defaults.
func DefaultSurgePolicy() SurgePolicy {
	return SurgePolicy{
		Window:    5 * time.Second,
		Threshold: 10,
	}
}

// SurgePolicy configures surge detection behaviour.
type SurgePolicy struct {
	Window    time.Duration
	Threshold int
}

// Surge detects sudden spikes in event rate within a sliding time window.
type Surge struct {
	mu       sync.Mutex
	policy   SurgePolicy
	events   []time.Time
	now      func() time.Time
}

// NewSurge creates a new Surge detector with the given policy.
// If now is nil, time.Now is used.
func NewSurge(policy SurgePolicy, now func() time.Time) *Surge {
	if policy.Window <= 0 {
		policy.Window = DefaultSurgePolicy().Window
	}
	if policy.Threshold <= 0 {
		policy.Threshold = DefaultSurgePolicy().Threshold
	}
	if now == nil {
		now = time.Now
	}
	return &Surge{policy: policy, now: now}
}

// Record registers a new event and reports whether a surge is detected.
func (s *Surge) Record() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := s.now()
	cutoff := now.Add(-s.policy.Window)
	filtered := s.events[:0]
	for _, t := range s.events {
		if t.After(cutoff) {
			filtered = append(filtered, t)
		}
	}
	filtered = append(filtered, now)
	s.events = filtered
	return len(s.events) >= s.policy.Threshold
}

// Count returns the number of events within the current window.
func (s *Surge) Count() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := s.now()
	cutoff := now.Add(-s.policy.Window)
	count := 0
	for _, t := range s.events {
		if t.After(cutoff) {
			count++
		}
	}
	return count
}

// Reset clears all recorded events.
func (s *Surge) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.events = s.events[:0]
}
