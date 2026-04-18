package watch

import (
	"testing"
	"time"
)

func makeJitter(fraction float64) *Jitter {
	return NewJitter(JitterPolicy{MaxFraction: fraction})
}

func TestDefaultJitterPolicyValues(t *testing.T) {
	p := DefaultJitterPolicy()
	if p.MaxFraction != 0.25 {
		t.Errorf("expected 0.25, got %v", p.MaxFraction)
	}
}

func TestJitterApplyReturnsAtLeastBase(t *testing.T) {
	j := makeJitter(0.5)
	base := 100 * time.Millisecond
	for i := 0; i < 20; i++ {
		result := j.Apply(base)
		if result < base {
			t.Errorf("jitter result %v less than base %v", result, base)
		}
	}
}

func TestJitterApplyDoesNotExceedMax(t *testing.T) {
	j := makeJitter(0.2)
	base := 100 * time.Millisecond
	max := base + time.Duration(float64(base)*0.2)
	for i := 0; i < 20; i++ {
		result := j.Apply(base)
		if result > max {
			t.Errorf("jitter result %v exceeds max %v", result, max)
		}
	}
}

func TestJitterApplyZeroBaseReturnsZero(t *testing.T) {
	j := makeJitter(0.5)
	if got := j.Apply(0); got != 0 {
		t.Errorf("expected 0, got %v", got)
	}
}

func TestJitterApplyNegativeBaseReturnsUnchanged(t *testing.T) {
	j := makeJitter(0.5)
	base := -10 * time.Millisecond
	if got := j.Apply(base); got != base {
		t.Errorf("expected %v, got %v", base, got)
	}
}

func TestJitterZeroFractionDefaultsToPolicy(t *testing.T) {
	j := NewJitter(JitterPolicy{MaxFraction: 0})
	if j.policy.MaxFraction != DefaultJitterPolicy().MaxFraction {
		t.Errorf("expected default fraction, got %v", j.policy.MaxFraction)
	}
}

func TestJitterResetDoesNotPanic(t *testing.T) {
	j := makeJitter(0.1)
	j.Reset()
	base := 50 * time.Millisecond
	result := j.Apply(base)
	if result < base {
		t.Errorf("after reset, result %v less than base %v", result, base)
	}
}
