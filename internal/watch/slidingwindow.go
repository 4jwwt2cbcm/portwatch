package watch

import (
	"sync"
	"time"
)

// SlidingWindowPolicy configures a SlidingWindow.
type SlidingWindowPolicy struct {
	Size     int
	Interval time.Duration
}

// DefaultSlidingWindowPolicy returns sensible defaults.
func DefaultSlidingWindowPolicy() SlidingWindowPolicy {
	return SlidingWindowPolicy{
		Size:     10,
		Interval: 60 * time.Second,
	}
}

// SlidingWindow tracks events over a rolling time window and reports
// whether the count exceeds the configured size threshold.
type SlidingWindow struct {
	policy SlidingWindowPolicy
	mu     sync.Mutex
	times  []time.Time
	now    func() time.Time
}

// NewSlidingWindow creates a SlidingWindow with the given policy.
// If policy.Size <= 0 or policy.Interval <= 0 defaults are applied.
func NewSlidingWindow(policy SlidingWindowPolicy) *SlidingWindow {
	def := DefaultSlidingWindowPolicy()
	if policy.Size <= 0 {
		policy.Size = def.Size
	}
	if policy.Interval <= 0 {
		policy.Interval = def.Interval
	}
	return &SlidingWindow{
		policy: policy,
		times:  make([]time.Time, 0, policy.Size),
		now:    time.Now,
	}
}

// Record adds an event timestamp and returns true if the window is full.
func (w *SlidingWindow) Record() bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	now := w.now()
	w.evict(now)
	w.times = append(w.times, now)
	return len(w.times) >= w.policy.Size
}

// Count returns the number of events within the current window.
func (w *SlidingWindow) Count() int {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.evict(w.now())
	return len(w.times)
}

// Reset clears all recorded events.
func (w *SlidingWindow) Reset() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.times = w.times[:0]
}

func (w *SlidingWindow) evict(now time.Time) {
	cutoff := now.Add(-w.policy.Interval)
	i := 0
	for i < len(w.times) && w.times[i].Before(cutoff) {
		i++
	}
	w.times = w.times[i:]
}
