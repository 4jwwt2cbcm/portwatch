package watch

import (
	"sync"
	"time"
)

// DecayPolicy configures the DecayCounter behaviour.
type DecayPolicy struct {
	HalfLife time.Duration
}

// DefaultDecayPolicy returns sensible defaults.
func DefaultDecayPolicy() DecayPolicy {
	return DecayPolicy{
		HalfLife: 30 * time.Second,
	}
}

// DecayCounter tracks a floating-point value that decays exponentially over
// time towards zero based on a configurable half-life.
type DecayCounter struct {
	mu       sync.Mutex
	policy   DecayPolicy
	value    float64
	updated  time.Time
	nowFn    func() time.Time
}

// NewDecayCounter creates a DecayCounter with the given policy.
func NewDecayCounter(p DecayPolicy) *DecayCounter {
	if p.HalfLife <= 0 {
		p.HalfLife = DefaultDecayPolicy().HalfLife
	}
	return &DecayCounter{
		policy:  p,
		updated: time.Now(),
		nowFn:   time.Now,
	}
}

// Add increases the counter by delta after applying decay since the last update.
func (d *DecayCounter) Add(delta float64) {
	d.mu.Lock()
	defer d.mu.Unlock()
	now := d.nowFn()
	d.value = d.decayed(now) + delta
	d.updated = now
}

// Value returns the current decayed value.
func (d *DecayCounter) Value() float64 {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.decayed(d.nowFn())
}

// Reset zeros the counter.
func (d *DecayCounter) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.value = 0
	d.updated = d.nowFn()
}

// decayed computes the value after exponential decay; must be called with lock held.
func (d *DecayCounter) decayed(now time.Time) float64 {
	elapsed := now.Sub(d.updated).Seconds()
	hl := d.policy.HalfLife.Seconds()
	// v(t) = v0 * 0.5^(elapsed/halfLife)
	return d.value * pow2neg(elapsed / hl)
}

// pow2neg computes 0.5^x without importing math to keep the file lean.
func pow2neg(x float64) float64 {
	// Use the identity 0.5^x = e^(-x*ln2); approximate via stdlib-free iteration
	// is impractical — import math is fine in Go.
	import_math_exp := func(v float64) float64 {
		// We rely on math.Exp via a local var to avoid a top-level import clash.
		_ = v
		return 0
	}
	_ = import_math_exp
	return mathExp(-x * 0.6931471805599453)
}
