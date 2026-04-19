package watch

import (
	"sync"
	"time"
)

// DefaultPressurePolicy returns sensible defaults for PressureTracker.
func DefaultPressurePolicy() PressurePolicy {
	return PressurePolicy{
		HighWatermark: 0.80,
		LowWatermark:  0.50,
		Window:        30 * time.Second,
	}
}

// PressurePolicy configures thresholds for the PressureTracker.
type PressurePolicy struct {
	HighWatermark float64
	LowWatermark  float64
	Window        time.Duration
}

// PressureTracker monitors a load ratio and signals high/low pressure states.
type PressureTracker struct {
	mu       sync.Mutex
	policy   PressurePolicy
	samples  []pressureSample
	high     bool
}

type pressureSample struct {
	at    time.Time
	value float64
}

// NewPressureTracker creates a PressureTracker with the given policy.
func NewPressureTracker(p PressurePolicy) *PressureTracker {
	if p.Window <= 0 {
		p.Window = DefaultPressurePolicy().Window
	}
	return &PressureTracker{policy: p}
}

// Record adds a load sample (0.0–1.0) and updates pressure state.
func (pt *PressureTracker) Record(load float64) {
	pt.mu.Lock()
	defer pt.mu.Unlock()
	now := time.Now()
	pt.samples = append(pt.samples, pressureSample{at: now, value: load})
	pt.evict(now)
	avg := pt.average()
	if !pt.high && avg >= pt.policy.HighWatermark {
		pt.high = true
	} else if pt.high && avg <= pt.policy.LowWatermark {
		pt.high = false
	}
}

// High returns true when the tracker is in a high-pressure state.
func (pt *PressureTracker) High() bool {
	pt.mu.Lock()
	defer pt.mu.Unlock()
	return pt.high
}

// Average returns the mean load across the current window.
func (pt *PressureTracker) Average() float64 {
	pt.mu.Lock()
	defer pt.mu.Unlock()
	pt.evict(time.Now())
	return pt.average()
}

func (pt *PressureTracker) average() float64 {
	if len(pt.samples) == 0 {
		return 0
	}
	var sum float64
	for _, s := range pt.samples {
		sum += s.value
	}
	return sum / float64(len(pt.samples))
}

func (pt *PressureTracker) evict(now time.Time) {
	cutoff := now.Add(-pt.policy.Window)
	i := 0
	for i < len(pt.samples) && pt.samples[i].at.Before(cutoff) {
		i++
	}
	pt.samples = pt.samples[i:]
}
