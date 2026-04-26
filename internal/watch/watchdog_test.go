package watch

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

func makeWatchdog(timeout, interval time.Duration, onFire func()) *Watchdog {
	return NewWatchdog(WatchdogPolicy{Timeout: timeout, Interval: interval}, onFire)
}

func TestDefaultWatchdogPolicyValues(t *testing.T) {
	p := DefaultWatchdogPolicy()
	if p.Timeout != 30*time.Second {
		t.Errorf("expected 30s timeout, got %v", p.Timeout)
	}
	if p.Interval != 5*time.Second {
		t.Errorf("expected 5s interval, got %v", p.Interval)
	}
}

func TestWatchdogDefaultsOnZero(t *testing.T) {
	w := NewWatchdog(WatchdogPolicy{}, nil)
	if w.policy.Timeout != 30*time.Second {
		t.Errorf("expected default timeout, got %v", w.policy.Timeout)
	}
	if w.policy.Interval != 5*time.Second {
		t.Errorf("expected default interval, got %v", w.policy.Interval)
	}
}

func TestWatchdogKickUpdatesLastKick(t *testing.T) {
	w := makeWatchdog(time.Second, 10*time.Millisecond, nil)
	before := w.LastKick()
	time.Sleep(5 * time.Millisecond)
	w.Kick()
	after := w.LastKick()
	if !after.After(before) {
		t.Error("expected LastKick to advance after Kick")
	}
}

func TestWatchdogFiresWhenExpired(t *testing.T) {
	var fired atomic.Int32
	w := makeWatchdog(20*time.Millisecond, 5*time.Millisecond, func() {
		fired.Add(1)
	})
	ctx, cancel := context.WithTimeout(context.Background(), 80*time.Millisecond)
	defer cancel()
	// do not kick — let it expire
	_ = w.Run(ctx)
	if fired.Load() == 0 {
		t.Error("expected watchdog to fire at least once")
	}
}

func TestWatchdogDoesNotFireWhenKicked(t *testing.T) {
	var fired atomic.Int32
	w := makeWatchdog(50*time.Millisecond, 5*time.Millisecond, func() {
		fired.Add(1)
	})
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		w.Run(ctx) //nolint:errcheck
		close(done)
	}()
	// kick regularly so the timeout never elapses
	for i := 0; i < 10; i++ {
		time.Sleep(5 * time.Millisecond)
		w.Kick()
	}
	cancel()
	<-done
	if fired.Load() != 0 {
		t.Errorf("expected no fires, got %d", fired.Load())
	}
}

func TestWatchdogStopsOnContextCancel(t *testing.T) {
	w := makeWatchdog(time.Second, 5*time.Millisecond, func() {})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := w.Run(ctx)
	if err == nil {
		t.Error("expected non-nil error on cancelled context")
	}
}
