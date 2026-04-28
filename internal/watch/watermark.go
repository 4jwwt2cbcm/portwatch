package watch

import (
	"sync"
	"time"
)

// DefaultWatermarkPolicy returns sensible defaults for a Watermark.
func DefaultWatermarkPolicy() WatermarkPolicy {
	return WatermarkPolicy{
		HighWater: 0.80,
		LowWater:  0.50,
	}
}

// WatermarkPolicy configures the high and low watermark thresholds.
// Values are expressed as fractions in the range [0, 1].
type WatermarkPolicy struct {
	HighWater float64
	LowWater  float64
}

// Watermark tracks a utilisation level and reports whether the high or low
// watermark has been crossed. It is safe for concurrent use.
type Watermark struct {
	mu       sync.Mutex
	policy   WatermarkPolicy
	level    float64
	above    bool
	updated  time.Time
}

// NewWatermark creates a Watermark with the given policy.
// If policy fields are zero the defaults are applied.
func NewWatermark(p WatermarkPolicy) *Watermark {
	def := DefaultWatermarkPolicy()
	if p.HighWater == 0 {
		p.HighWater = def.HighWater
	}
	if p.LowWater == 0 {
		p.LowWater = def.LowWater
	}
	return &Watermark{policy: p}
}

// Set updates the current utilisation level (0–1) and recalculates whether
// the tracker is above the high watermark or has dropped below the low one.
func (w *Watermark) Set(level float64) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.level = level
	w.updated = time.Now()
	if !w.above && level >= w.policy.HighWater {
		w.above = true
	} else if w.above && level < w.policy.LowWater {
		w.above = false
	}
}

// Above reports whether the level is currently above the high watermark and
// has not yet fallen back below the low watermark.
func (w *Watermark) Above() bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.above
}

// Level returns the most recently recorded utilisation level.
func (w *Watermark) Level() float64 {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.level
}

// LastUpdated returns the time of the most recent Set call.
func (w *Watermark) LastUpdated() time.Time {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.updated
}

// Reset clears the level and the above flag.
func (w *Watermark) Reset() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.level = 0
	w.above = false
	w.updated = time.Time{}
}
