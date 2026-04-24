package watch

import (
	"sync"
	"time"
)

// DefaultShedderPolicy returns a ShedderPolicy with sensible defaults.
func DefaultShedderPolicy() ShedderPolicy {
	return ShedderPolicy{
		MaxLoad:  0.80,
		Window:   5 * time.Second,
		Cooldown: 2 * time.Second,
	}
}

// ShedderPolicy controls load-shedding thresholds.
type ShedderPolicy struct {
	// MaxLoad is the fraction of capacity [0,1] above which work is shed.
	MaxLoad float64
	// Window is the measurement window for load averaging.
	Window time.Duration
	// Cooldown is the minimum time between successive sheds.
	Cooldown time.Duration
}

// Shedder drops work when the system is above a configured load threshold.
type Shedder struct {
	mu       sync.Mutex
	policy   ShedderPolicy
	samples  []loadSample
	lastShed time.Time
}

type loadSample struct {
	at    time.Time
	value float64
}

// NewShedder creates a Shedder with the given policy.
func NewShedder(p ShedderPolicy) *Shedder {
	if p.MaxLoad <= 0 {
		p.MaxLoad = DefaultShedderPolicy().MaxLoad
	}
	if p.Window <= 0 {
		p.Window = DefaultShedderPolicy().Window
	}
	if p.Cooldown <= 0 {
		p.Cooldown = DefaultShedderPolicy().Cooldown
	}
	return &Shedder{policy: p}
}

// Record adds a load observation in [0,1].
func (s *Shedder) Record(load float64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	s.evict(now)
	s.samples = append(s.samples, loadSample{at: now, value: load})
}

// ShouldShed reports whether the current load average exceeds the threshold
// and the cooldown period has elapsed since the last shed.
func (s *Shedder) ShouldShed() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	s.evict(now)
	if len(s.samples) == 0 {
		return false
	}
	var sum float64
	for _, sample := range s.samples {
		sum += sample.value
	}
	avg := sum / float64(len(s.samples))
	if avg < s.policy.MaxLoad {
		return false
	}
	if now.Sub(s.lastShed) < s.policy.Cooldown {
		return false
	}
	s.lastShed = now
	return true
}

// Load returns the current average load across the window.
func (s *Shedder) Load() float64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.evict(time.Now())
	if len(s.samples) == 0 {
		return 0
	}
	var sum float64
	for _, sample := range s.samples {
		sum += sample.value
	}
	return sum / float64(len(s.samples))
}

func (s *Shedder) evict(now time.Time) {
	cutoff := now.Add(-s.policy.Window)
	i := 0
	for i < len(s.samples) && s.samples[i].at.Before(cutoff) {
		i++
	}
	s.samples = s.samples[i:]
}
