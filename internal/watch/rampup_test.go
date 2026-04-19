package watch

import (
	"testing"
	"time"
)

func makeRampUp() *RampUp {
	return NewRampUp(RampUpPolicy{
		Steps:   4,
		Initial: 100 * time.Millisecond,
		Target:  500 * time.Millisecond,
	})
}

func TestDefaultRampUpPolicyValues(t *testing.T) {
	p := DefaultRampUpPolicy()
	if p.Steps != 5 {
		t.Errorf("expected 5 steps, got %d", p.Steps)
	}
	if p.Initial != 200*time.Millisecond {
		t.Errorf("unexpected initial: %v", p.Initial)
	}
	if p.Target != 2*time.Second {
		t.Errorf("unexpected target: %v", p.Target)
	}
}

func TestRampUpFirstNextReturnsInitial(t *testing.T) {
	r := makeRampUp()
	got := r.Next()
	if got != 100*time.Millisecond {
		t.Errorf("expected 100ms, got %v", got)
	}
}

func TestRampUpGrowsTowardsTarget(t *testing.T) {
	r := makeRampUp()
	var last time.Duration
	for i := 0; i < 4; i++ {
		v := r.Next()
		if v < last {
			t.Errorf("interval decreased at step %d: %v < %v", i, v, last)
		}
		last = v
	}
}

func TestRampUpDoneAfterAllSteps(t *testing.T) {
	r := makeRampUp()
	for i := 0; i <= 4; i++ {
		r.Next()
	}
	if !r.Done() {
		t.Error("expected Done() to be true after all steps")
	}
}

func TestRampUpNotDoneInitially(t *testing.T) {
	r := makeRampUp()
	if r.Done() {
		t.Error("expected Done() to be false on init")
	}
}

func TestRampUpResetRestoresInitial(t *testing.T) {
	r := makeRampUp()
	for i := 0; i <= 4; i++ {
		r.Next()
	}
	r.Reset()
	if r.Done() {
		t.Error("expected Done() false after reset")
	}
	got := r.Next()
	if got != 100*time.Millisecond {
		t.Errorf("expected 100ms after reset, got %v", got)
	}
}

func TestRampUpClampsToTarget(t *testing.T) {
	r := makeRampUp()
	var last time.Duration
	for i := 0; i < 10; i++ {
		v := r.Next()
		if v > 500*time.Millisecond {
			t.Errorf("exceeded target at step %d: %v", i, v)
		}
		last = v
	}
	if last != 500*time.Millisecond {
		t.Errorf("expected final value to be target 500ms, got %v", last)
	}
}
