package watch

import (
	"testing"
	"time"
)

func TestDefaultBackoffPolicyValues(t *testing.T) {
	p := DefaultBackoffPolicy()
	if p.Initial != 500*time.Millisecond {
		t.Errorf("expected Initial=500ms, got %v", p.Initial)
	}
	if p.Max != 30*time.Second {
		t.Errorf("expected Max=30s, got %v", p.Max)
	}
	if p.Factor != 2.0 {
		t.Errorf("expected Factor=2.0, got %v", p.Factor)
	}
}

func TestBackoffFirstFailureReturnsInitial(t *testing.T) {
	b := NewBackoff(BackoffPolicy{Initial: 100 * time.Millisecond, Max: 1 * time.Second, Factor: 2.0})
	d := b.Next()
	if d != 100*time.Millisecond {
		t.Errorf("expected 100ms, got %v", d)
	}
}

func TestBackoffExponentialGrowth(t *testing.T) {
	b := NewBackoff(BackoffPolicy{Initial: 100 * time.Millisecond, Max: 10 * time.Second, Factor: 2.0})
	b.Next() // 100ms
	d := b.Next() // 200ms
	if d != 200*time.Millisecond {
		t.Errorf("expected 200ms, got %v", d)
	}
}

func TestBackoffCapsAtMax(t *testing.T) {
	b := NewBackoff(BackoffPolicy{Initial: 1 * time.Second, Max: 2 * time.Second, Factor: 10.0})
	b.Next()
	d := b.Next()
	if d > 2*time.Second {
		t.Errorf("expected cap at 2s, got %v", d)
	}
}

func TestBackoffResetRestoresInitial(t *testing.T) {
	b := NewBackoff(BackoffPolicy{Initial: 100 * time.Millisecond, Max: 1 * time.Second, Factor: 2.0})
	b.Next()
	b.Next()
	b.Reset()
	if b.Current() != 100*time.Millisecond {
		t.Errorf("expected reset to 100ms, got %v", b.Current())
	}
}

func TestBackoffCurrentDoesNotAdvance(t *testing.T) {
	b := NewBackoff(BackoffPolicy{Initial: 200 * time.Millisecond, Max: 1 * time.Second, Factor: 2.0})
	_ = b.Current()
	_ = b.Current()
	if b.Current() != 200*time.Millisecond {
		t.Errorf("expected Current to not advance, got %v", b.Current())
	}
}
