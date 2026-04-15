package watch

import (
	"testing"
	"time"
)

func TestDefaultBackoffPolicyValues(t *testing.T) {
	p := DefaultBackoffPolicy()
	if p.Initial != 2*time.Second {
		t.Errorf("expected Initial=2s, got %v", p.Initial)
	}
	if p.Multiplier != 2.0 {
		t.Errorf("expected Multiplier=2.0, got %v", p.Multiplier)
	}
	if p.Max != 60*time.Second {
		t.Errorf("expected Max=60s, got %v", p.Max)
	}
}

func TestBackoffFirstFailureReturnsInitial(t *testing.T) {
	b := NewBackoff(DefaultBackoffPolicy())
	delay := b.Failure()
	if delay != 2*time.Second {
		t.Errorf("expected 2s, got %v", delay)
	}
}

func TestBackoffExponentialGrowth(t *testing.T) {
	b := NewBackoff(DefaultBackoffPolicy())
	delays := make([]time.Duration, 4)
	for i := range delays {
		delays[i] = b.Failure()
	}
	expected := []time.Duration{2 * time.Second, 4 * time.Second, 8 * time.Second, 16 * time.Second}
	for i, d := range delays {
		if d != expected[i] {
			t.Errorf("delay[%d]: expected %v, got %v", i, expected[i], d)
		}
	}
}

func TestBackoffCapsAtMax(t *testing.T) {
	p := BackoffPolicy{Initial: 30 * time.Second, Multiplier: 2.0, Max: 60 * time.Second}
	b := NewBackoff(p)
	b.Failure() // 30s
	delay := b.Failure() // would be 60s
	if delay != 60*time.Second {
		t.Errorf("expected cap at 60s, got %v", delay)
	}
	delay = b.Failure() // should still be 60s
	if delay != 60*time.Second {
		t.Errorf("expected cap maintained at 60s, got %v", delay)
	}
}

func TestBackoffResetRestoresInitial(t *testing.T) {
	b := NewBackoff(DefaultBackoffPolicy())
	b.Failure()
	b.Failure()
	b.Reset()
	if b.Failures() != 0 {
		t.Errorf("expected 0 failures after reset, got %d", b.Failures())
	}
	delay := b.Failure()
	if delay != 2*time.Second {
		t.Errorf("expected initial delay after reset, got %v", delay)
	}
}

func TestBackoffFailuresCount(t *testing.T) {
	b := NewBackoff(DefaultBackoffPolicy())
	for i := 1; i <= 3; i++ {
		b.Failure()
		if b.Failures() != i {
			t.Errorf("expected %d failures, got %d", i, b.Failures())
		}
	}
}
