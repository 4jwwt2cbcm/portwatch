package watch

import (
	"context"
	"testing"
	"time"
)

func makeLimiter(max int, window time.Duration) *Limiter {
	return NewLimiter(LimiterPolicy{MaxCalls: max, Window: window})
}

func TestDefaultLimiterPolicyValues(t *testing.T) {
	p := DefaultLimiterPolicy()
	if p.MaxCalls != 10 {
		t.Errorf("expected MaxCalls 10, got %d", p.MaxCalls)
	}
	if p.Window != time.Minute {
		t.Errorf("expected Window 1m, got %v", p.Window)
	}
}

func TestLimiterAllowsUpToMax(t *testing.T) {
	l := makeLimiter(3, time.Second)
	for i := 0; i < 3; i++ {
		if !l.Allow() {
			t.Fatalf("expected allow on call %d", i+1)
		}
	}
	if l.Allow() {
		t.Error("expected deny after max calls")
	}
}

func TestLimiterResetClearsState(t *testing.T) {
	l := makeLimiter(1, time.Second)
	l.Allow()
	if l.Allow() {
		t.Error("expected deny before reset")
	}
	l.Reset()
	if !l.Allow() {
		t.Error("expected allow after reset")
	}
}

func TestLimiterEvictsExpiredCalls(t *testing.T) {
	l := makeLimiter(1, 50*time.Millisecond)
	base := time.Now()
	l.now = func() time.Time { return base }
	l.Allow()
	if l.Allow() {
		t.Error("expected deny within window")
	}
	l.now = func() time.Time { return base.Add(100 * time.Millisecond) }
	if !l.Allow() {
		t.Error("expected allow after window expired")
	}
}

func TestLimiterWaitReturnsOnAllow(t *testing.T) {
	l := makeLimiter(2, 50*time.Millisecond)
	l.Allow()
	l.Allow()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := l.Wait(ctx); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestLimiterWaitCancelledContext(t *testing.T) {
	l := makeLimiter(1, time.Hour)
	l.Allow()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := l.Wait(ctx); err == nil {
		t.Error("expected error on cancelled context")
	}
}
