package watch

import (
	"sync"
	"time"
)

// DefaultStaggerPolicy returns a StaggerPolicy with sensible defaults.
func DefaultStaggerPolicy() StaggerPolicy {
	return StaggerPolicy{
		Delay:    100 * time.Millisecond,
		MaxItems: 10,
	}
}

// StaggerPolicy controls how a Stagger distributes work over time.
type StaggerPolicy struct {
	Delay    time.Duration
	MaxItems int
}

// Stagger spaces out sequential function calls by a fixed delay,
// preventing thundering-herd bursts when many tasks start at once.
type Stagger struct {
	mu     sync.Mutex
	policy StaggerPolicy
	next   time.Time
	count  int
}

// NewStagger creates a Stagger with the given policy.
// Zero-value fields are replaced with defaults.
func NewStagger(p StaggerPolicy) *Stagger {
	def := DefaultStaggerPolicy()
	if p.Delay <= 0 {
		p.Delay = def.Delay
	}
	if p.MaxItems <= 0 {
		p.MaxItems = def.MaxItems
	}
	return &Stagger{policy: p}
}

// Next returns the time at which the next item should be dispatched.
// Each call advances the internal cursor by one delay slot.
func (s *Stagger) Next() time.Time {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	if s.next.IsZero() || s.next.Before(now) {
		s.next = now
	}
	t := s.next
	s.next = s.next.Add(s.policy.Delay)
	s.count++
	return t
}

// WaitNext blocks until the next staggered slot is reached.
func (s *Stagger) WaitNext() {
	t := s.Next()
	wait := time.Until(t)
	if wait > 0 {
		time.Sleep(wait)
	}
}

// Count returns the total number of slots that have been issued.
func (s *Stagger) Count() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.count
}

// Reset clears the stagger cursor and count.
func (s *Stagger) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.next = time.Time{}
	s.count = 0
}
