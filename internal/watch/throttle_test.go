package watch

import (
	"testing"
	"time"
)

func makeThrottle(cooldown time.Duration) (*Throttle, *time.Time) {
	now := time.Now()
	t := NewThrottle(cooldown)
	t.now = func() time.Time { return now }
	return t, &now
}

func TestThrottleFirstCallAllowed(t *testing.T) {
	th, _ := makeThrottle(5 * time.Second)
	if !th.Allow() {
		t.Fatal("expected first call to be allowed")
	}
}

func TestThrottleSecondCallSuppressed(t *testing.T) {
	th, _ := makeThrottle(5 * time.Second)
	th.Allow()
	if th.Allow() {
		t.Fatal("expected second call within cooldown to be suppressed")
	}
}

func TestThrottleAllowedAfterCooldown(t *testing.T) {
	now := time.Now()
	th := NewThrottle(5 * time.Second)
	th.now = func() time.Time { return now }

	th.Allow()

	// advance time past the cooldown
	now = now.Add(6 * time.Second)
	th.now = func() time.Time { return now }

	if !th.Allow() {
		t.Fatal("expected call after cooldown to be allowed")
	}
}

func TestThrottleResetAllowsImmediately(t *testing.T) {
	th, _ := makeThrottle(5 * time.Second)
	th.Allow()
	th.Reset()
	if !th.Allow() {
		t.Fatal("expected call after reset to be allowed")
	}
}

func TestThrottleRemainingZeroBeforeFirstCall(t *testing.T) {
	th, _ := makeThrottle(5 * time.Second)
	if r := th.Remaining(); r != 0 {
		t.Fatalf("expected 0 remaining before first call, got %v", r)
	}
}

func TestThrottleRemainingAfterCall(t *testing.T) {
	now := time.Now()
	th := NewThrottle(10 * time.Second)
	th.now = func() time.Time { return now }

	th.Allow()

	now = now.Add(3 * time.Second)
	th.now = func() time.Time { return now }

	got := th.Remaining()
	want := 7 * time.Second
	if got != want {
		t.Fatalf("expected remaining=%v, got %v", want, got)
	}
}

func TestThrottleRemainingZeroAfterCooldown(t *testing.T) {
	now := time.Now()
	th := NewThrottle(5 * time.Second)
	th.now = func() time.Time { return now }

	th.Allow()

	now = now.Add(10 * time.Second)
	th.now = func() time.Time { return now }

	if r := th.Remaining(); r != 0 {
		t.Fatalf("expected 0 remaining after cooldown, got %v", r)
	}
}
