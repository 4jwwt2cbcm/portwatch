package watch

import (
	"sync"
	"time"
)

// SplitBrainPolicy configures the SplitBrain detector.
type SplitBrainPolicy struct {
	// QuorumSize is the minimum number of agreeing nodes required.
	QuorumSize int
	// Window is the time range in which votes are considered valid.
	Window time.Duration
}

// DefaultSplitBrainPolicy returns sensible defaults.
func DefaultSplitBrainPolicy() SplitBrainPolicy {
	return SplitBrainPolicy{
		QuorumSize: 2,
		Window:     10 * time.Second,
	}
}

// SplitBrain detects whether a quorum of nodes agree on a value within a
// sliding time window. It is useful for detecting split-brain conditions
// where multiple sources report conflicting state.
type SplitBrain struct {
	mu     sync.Mutex
	policy SplitBrainPolicy
	votes  map[string][]time.Time
}

// NewSplitBrain creates a new SplitBrain detector with the given policy.
// A zero-value policy falls back to DefaultSplitBrainPolicy.
func NewSplitBrain(p SplitBrainPolicy) *SplitBrain {
	if p.QuorumSize <= 0 {
		p.QuorumSize = DefaultSplitBrainPolicy().QuorumSize
	}
	if p.Window <= 0 {
		p.Window = DefaultSplitBrainPolicy().Window
	}
	return &SplitBrain{
		policy: p,
		votes:  make(map[string][]time.Time),
	}
}

// Vote records a vote for the given value at the current time.
func (s *SplitBrain) Vote(value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	s.votes[value] = append(s.evict(s.votes[value], now), now)
}

// HasQuorum returns true if any single value has received at least
// QuorumSize votes within the current window.
func (s *SplitBrain) HasQuorum() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	for _, ts := range s.votes {
		if len(s.evict(ts, now)) >= s.policy.QuorumSize {
			return true
		}
	}
	return false
}

// Conflicted returns true if two or more distinct values have at least one
// vote within the current window, indicating a potential split-brain.
func (s *SplitBrain) Conflicted() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	active := 0
	for _, ts := range s.votes {
		if len(s.evict(ts, now)) > 0 {
			active++
		}
		if active >= 2 {
			return true
		}
	}
	return false
}

// Reset clears all recorded votes.
func (s *SplitBrain) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.votes = make(map[string][]time.Time)
}

// evict removes timestamps older than the window from the slice.
func (s *SplitBrain) evict(ts []time.Time, now time.Time) []time.Time {
	cutoff := now.Add(-s.policy.Window)
	for len(ts) > 0 && ts[0].Before(cutoff) {
		ts = ts[1:]
	}
	return ts
}
