package watch

import (
	"context"
	"testing"
	"time"
)

func TestDefaultHeartbeatPolicyValues(t *testing.T) {
	p := DefaultHeartbeatPolicy()
	if p.Interval != 30*time.Second {
		t.Errorf("expected 30s interval, got %v", p.Interval)
	}
	if p.Timeout != 5*time.Second {
		t.Errorf("expected 5s timeout, got %v", p.Timeout)
	}
}

func TestHeartbeatNotAliveOnInit(t *testing.T) {
	h := NewHeartbeat(DefaultHeartbeatPolicy())
	if h.Alive() {
		t.Error("expected not alive before first beat")
	}
}

func TestHeartbeatAliveAfterBeat(t *testing.T) {
	p := DefaultHeartbeatPolicy()
	p.Timeout = time.Second
	h := NewHeartbeat(p)
	h.Beat()
	if !h.Alive() {
		t.Error("expected alive immediately after beat")
	}
}

func TestHeartbeatNotAliveAfterTimeout(t *testing.T) {
	p := HeartbeatPolicy{Interval: time.Second, Timeout: 10 * time.Millisecond}
	h := NewHeartbeat(p)
	h.Beat()
	time.Sleep(20 * time.Millisecond)
	if h.Alive() {
		t.Error("expected not alive after timeout elapsed")
	}
}

func TestHeartbeatLastUpdatedOnBeat(t *testing.T) {
	h := NewHeartbeat(DefaultHeartbeatPolicy())
	if !h.Last().IsZero() {
		t.Error("expected zero last before beat")
	}
	before := time.Now()
	h.Beat()
	if h.Last().Before(before) {
		t.Error("last should be >= time before beat")
	}
}

func TestHeartbeatZeroPolicyFallsToDefault(t *testing.T) {
	h := NewHeartbeat(HeartbeatPolicy{})
	if h.policy.Interval != 30*time.Second {
		t.Errorf("expected default interval, got %v", h.policy.Interval)
	}
}

func TestHeartbeatRunStopsOnContextCancel(t *testing.T) {
	p := HeartbeatPolicy{Interval: 10 * time.Millisecond, Timeout: time.Second}
	h := NewHeartbeat(p)
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- h.Run(ctx) }()
	time.Sleep(25 * time.Millisecond)
	cancel()
	select {
	case err := <-done:
		if err != context.Canceled {
			t.Errorf("expected context.Canceled, got %v", err)
		}
	case <-time.After(time.Second):
		t.Error("Run did not stop after context cancel")
	}
	if h.Last().IsZero() {
		t.Error("expected at least one beat during run")
	}
}
