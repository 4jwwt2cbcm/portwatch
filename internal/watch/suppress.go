package watch

import (
	"sync"
	"time"
)

// DefaultSuppressPolicy returns a policy that suppresses repeated events
// for the same key within a 30-second window.
func DefaultSuppressPolicy() SuppressPolicy {
	return SuppressPolicy{
		Window: 30 * time.Second,
		MaxKeys: 1024,
	}
}

// SuppressPolicy configures the Suppressor.
type SuppressPolicy struct {
	Window  time.Duration
	MaxKeys int
}

type suppressEntry struct {
	last time.Time
	count int
}

// Suppressor prevents duplicate events for a key within a time window.
type Suppressor struct {
	policy  SuppressPolicy
	mu      sync.Mutex
	entries map[string]*suppressEntry
}

// NewSuppressor creates a Suppressor with the given policy.
func NewSuppressor(p SuppressPolicy) *Suppressor {
	if p.Window <= 0 {
		p.Window = DefaultSuppressPolicy().Window
	}
	if p.MaxKeys <= 0 {
		p.MaxKeys = DefaultSuppressPolicy().MaxKeys
	}
	return &Suppressor{
		policy:  p,
		entries: make(map[string]*suppressEntry),
	}
}

// Allow returns true if the event for key should be allowed through.
// Subsequent calls within the window are suppressed.
func (s *Suppressor) Allow(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	if e, ok := s.entries[key]; ok {
		if now.Sub(e.last) < s.policy.Window {
			e.count++
			return false
		}
	}
	if len(s.entries) >= s.policy.MaxKeys {
		s.entries = make(map[string]*suppressEntry)
	}
	s.entries[key] = &suppressEntry{last: now, count: 1}
	return true
}

// Count returns how many times key has been suppressed in the current window.
func (s *Suppressor) Count(key string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	if e, ok := s.entries[key]; ok {
		return e.count
	}
	return 0
}

// Reset clears suppression state for all keys.
func (s *Suppressor) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries = make(map[string]*suppressEntry)
}
