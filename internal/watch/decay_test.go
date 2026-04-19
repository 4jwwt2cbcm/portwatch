package watch

import (
	"testing"
	"time"
)

func makeDecayCounter(halfLife time.Duration) *DecayCounter {
	dc := NewDecayCounter(DecayPolicy{HalfLife: halfLife})
	return dc
}

func TestDefaultDecayPolicyValues(t *testing.T) {
	p := DefaultDecayPolicy()
	if p.HalfLife != 30*time.Second {
		t.Errorf("expected 30s half-life, got %v", p.HalfLife)
	}
}

func TestDecayCounterZeroOnInit(t *testing.T) {
	dc := makeDecayCounter(time.Second)
	if v := dc.Value(); v != 0 {
		t.Errorf("expected 0, got %f", v)
	}
}

func TestDecayCounterAddIncreasesValue(t *testing.T) {
	dc := makeDecayCounter(time.Minute)
	dc.Add(10)
	if v := dc.Value(); v <= 0 {
		t.Errorf("expected positive value after Add, got %f", v)
	}
}

func TestDecayCounterValueDecreasesOverTime(t *testing.T) {
	dc := makeDecayCounter(time.Second)
	now := time.Now()
	dc.nowFn = func() time.Time { return now }
	dc.Add(100)

	// Advance time by one half-life; value should be ~50.
	dc.nowFn = func() time.Time { return now.Add(time.Second) }
	v := dc.Value()
	if v < 45 || v > 55 {
		t.Errorf("expected ~50 after one half-life, got %f", v)
	}
}

func TestDecayCounterResetZerosValue(t *testing.T) {
	dc := makeDecayCounter(time.Minute)
	dc.Add(42)
	dc.Reset()
	if v := dc.Value(); v != 0 {
		t.Errorf("expected 0 after Reset, got %f", v)
	}
}

func TestDecayCounterDefaultsHalfLifeOnZero(t *testing.T) {
	dc := NewDecayCounter(DecayPolicy{HalfLife: 0})
	if dc.policy.HalfLife != DefaultDecayPolicy().HalfLife {
		t.Errorf("expected default half-life, got %v", dc.policy.HalfLife)
	}
}
