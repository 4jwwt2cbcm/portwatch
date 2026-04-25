package watch

import (
	"sync"
	"time"
)

// ShedCounterPolicy configures the ShedCounter.
type ShedCounterPolicy struct {
	Window   time.Duration
	Capacity int
}

// DefaultShedCounterPolicy returns sensible defaults.
func DefaultShedCounterPolicy() ShedCounterPolicy {
	return ShedCounterPolicy{
		Window:   30 * time.Second,
		Capacity: 1000,
	}
}

// ShedCounter tracks how many tasks have been shed within a rolling window.
type ShedCounter struct {
	policy  ShedCounterPolicy
	mu      sync.Mutex
	entries []time.Time
}

// NewShedCounter creates a ShedCounter using the given policy.
// Zero-value fields fall back to defaults.
func NewShedCounter(p ShedCounterPolicy) *ShedCounter {
	def := DefaultShedCounterPolicy()
	if p.Window <= 0 {
		p.Window = def.Window
	}
	if p.Capacity <= 0 {
		p.Capacity = def.Capacity
	}
	return &ShedCounter{policy: p}
}

// Record marks one shed event at the current time.
func (s *ShedCounter) Record() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.evict()
	if len(s.entries) < s.policy.Capacity {
		s.entries = append(s.entries, time.Now())
	}
}

// Count returns the number of shed events within the current window.
func (s *ShedCounter) Count() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.evict()
	return len(s.entries)
}

// Reset clears all recorded shed events.
func (s *ShedCounter) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries = s.entries[:0]
}

// evict removes entries older than the window. Must be called with mu held.
func (s *ShedCounter) evict() {
	cutoff := time.Now().Add(-s.policy.Window)
	i := 0
	for i < len(s.entries) && s.entries[i].Before(cutoff) {
		i++
	}
	s.entries = s.entries[i:]
}
