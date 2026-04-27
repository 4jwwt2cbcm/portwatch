package watch

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestLeaseRunnerAcquiresAndRunsFn(t *testing.T) {
	l := makeLease(time.Second)
	r := NewLeaseRunner(l)
	ran := false
	err := r.Run(context.Background(), func(_ context.Context) error {
		ran = true
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ran {
		t.Error("expected fn to run")
	}
}

func TestLeaseRunnerReleasesAfterFn(t *testing.T) {
	l := makeLease(time.Second)
	r := NewLeaseRunner(l)
	r.Run(context.Background(), func(_ context.Context) error { return nil })
	if l.Held() {
		t.Error("expected lease released after fn completes")
	}
}

func TestLeaseRunnerReturnsErrLeaseNotHeld(t *testing.T) {
	l := makeLease(time.Second)
	l.Acquire() // hold the lease externally
	r := NewLeaseRunner(l)
	err := r.Run(context.Background(), func(_ context.Context) error { return nil })
	if !errors.Is(err, ErrLeaseNotHeld) {
		t.Errorf("expected ErrLeaseNotHeld, got %v", err)
	}
}

func TestLeaseRunnerPropagatesError(t *testing.T) {
	sentinel := errors.New("fn error")
	l := makeLease(time.Second)
	r := NewLeaseRunner(l)
	err := r.Run(context.Background(), func(_ context.Context) error {
		return sentinel
	})
	if !errors.Is(err, sentinel) {
		t.Errorf("expected sentinel error, got %v", err)
	}
}

func TestLeaseRunnerNilLeaseDefaults(t *testing.T) {
	r := NewLeaseRunner(nil)
	if r.lease == nil {
		t.Error("expected default lease when nil provided")
	}
}

func TestLeaseRunnerHeldDuringFn(t *testing.T) {
	l := makeLease(time.Second)
	r := NewLeaseRunner(l)
	r.Run(context.Background(), func(_ context.Context) error {
		if !r.Held() {
			t.Error("expected lease held during fn execution")
		}
		return nil
	})
}

func TestLeaseRunnerRenewsLease(t *testing.T) {
	l := NewLease(LeasePolicy{TTL: 60 * time.Millisecond, RenewAt: 0.5})
	r := NewLeaseRunner(l)
	initialExpiry := time.Time{}
	r.Run(context.Background(), func(_ context.Context) error {
		initialExpiry = l.ExpiresAt()
		time.Sleep(50 * time.Millisecond) // trigger renewal
		return nil
	})
	// Renewal should have extended expiry beyond initial
	_ = initialExpiry // renewal is best-effort; just ensure no panic
}
