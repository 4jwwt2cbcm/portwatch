package watch

import (
	"sync"
	"time"
)

// SlotPolicy configures the time slot behaviour.
type SlotPolicy struct {
	Duration time.Duration // length of each slot
	MaxSlots int           // maximum number of distinct slots tracked
}

// DefaultSlotPolicy returns sensible defaults.
func DefaultSlotPolicy() SlotPolicy {
	return SlotPolicy{
		Duration: time.Minute,
		MaxSlots: 60,
	}
}

// Slot divides time into fixed-width buckets and counts events per bucket.
type Slot struct {
	mu     sync.Mutex
	policy SlotPolicy
	buckets map[int64]int
}

// NewSlot creates a Slot with the given policy. Zero-value fields are
// replaced with defaults.
func NewSlot(p SlotPolicy) *Slot {
	if p.Duration <= 0 {
		p.Duration = DefaultSlotPolicy().Duration
	}
	if p.MaxSlots <= 0 {
		p.MaxSlots = DefaultSlotPolicy().MaxSlots
	}
	return &Slot{
		policy:  p,
		buckets: make(map[int64]int),
	}
}

// key returns the bucket key for the given time.
func (s *Slot) key(t time.Time) int64 {
	return t.UnixNano() / int64(s.policy.Duration)
}

// Record increments the counter for the current time slot.
func (s *Slot) Record(t time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()
	k := s.key(t)
	s.buckets[k]++
	s.evict(k)
}

// Count returns the number of events recorded in the slot containing t.
func (s *Slot) Count(t time.Time) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.buckets[s.key(t)]
}

// Reset clears all slot data.
func (s *Slot) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.buckets = make(map[int64]int)
}

// evict removes old buckets when the map exceeds MaxSlots, keeping the newest.
func (s *Slot) evict(current int64) {
	if len(s.buckets) <= s.policy.MaxSlots {
		return
	}
	// find and remove the oldest key
	var oldest int64 = current
	for k := range s.buckets {
		if k < oldest {
			oldest = k
		}
	}
	delete(s.buckets, oldest)
}
