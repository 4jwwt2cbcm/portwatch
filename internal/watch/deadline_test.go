package watch

import (
	"context"
	"testing"
	"time"
)

func makeDeadline(at time.Time) *Deadline {
	return NewDeadline(DeadlinePolicy{At: at})
}

func TestDefaultDeadlinePolicyValues(t *testing.T) {
	p := DefaultDeadlinePolicy()
	if !p.At.IsZero() {
		t.Fatalf("expected zero time, got %v", p.At)
	}
}

func TestDeadlineNotExpiredOnZero(t *testing.T) {
	d := makeDeadline(time.Time{})
	if d.Expired() {
		t.Fatal("expected not expired for zero deadline")
	}
}

func TestDeadlineExpiredInPast(t *testing.T) {
	d := makeDeadline(time.Now().Add(-time.Second))
	if !d.Expired() {
		t.Fatal("expected expired for past deadline")
	}
}

func TestDeadlineNotExpiredInFuture(t *testing.T) {
	d := makeDeadline(time.Now().Add(time.Hour))
	if d.Expired() {
		t.Fatal("expected not expired for future deadline")
	}
}

func TestDeadlineSetUpdates(t *testing.T) {
	d := makeDeadline(time.Time{})
	now := time.Now().Add(-time.Second)
	d.Set(now)
	if !d.Expired() {
		t.Fatal("expected expired after Set to past")
	}
}

func TestDeadlineAtReturnsValue(t *testing.T) {
	at := time.Now().Add(time.Minute)
	d := makeDeadline(at)
	if !d.At().Equal(at) {
		t.Fatalf("expected %v, got %v", at, d.At())
	}
}

func TestDeadlineWaitCancelledContext(t *testing.T) {
	d := makeDeadline(time.Now().Add(time.Hour))
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := d.Wait(ctx)
	if err != context.Canceled {
		t.Fatalf("expected Canceled, got %v", err)
	}
}

func TestDeadlineWaitFiresDeadline(t *testing.T) {
	d := makeDeadline(time.Now().Add(20 * time.Millisecond))
	ctx := context.Background()
	err := d.Wait(ctx)
	if err != context.DeadlineExceeded {
		t.Fatalf("expected DeadlineExceeded, got %v", err)
	}
}
