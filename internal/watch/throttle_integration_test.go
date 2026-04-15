package watch_test

import (
	"testing"
	"time"

	"github.com/yourorg/portwatch/internal/watch"
)

// TestThrottleIntegration verifies real-time throttling behaviour end-to-end
// using actual wall-clock time with a short cooldown.
func TestThrottleIntegration(t *testing.T) {
	cooldown := 50 * time.Millisecond
	th := watch.NewThrottle(cooldown)

	// First call should always be allowed.
	if !th.Allow() {
		t.Fatal("expected first Allow() to return true")
	}

	// Immediate second call should be suppressed.
	if th.Allow() {
		t.Fatal("expected immediate second Allow() to return false")
	}

	// Remaining should be positive and within the cooldown window.
	remaining := th.Remaining()
	if remaining <= 0 || remaining > cooldown {
		t.Fatalf("unexpected Remaining() value: %v", remaining)
	}

	// Wait for the cooldown to elapse.
	time.Sleep(cooldown + 10*time.Millisecond)

	// Should be allowed again after cooldown.
	if !th.Allow() {
		t.Fatal("expected Allow() to return true after cooldown elapsed")
	}

	// Remaining should now be zero (or close to it) right after firing.
	if r := th.Remaining(); r > cooldown {
		t.Fatalf("Remaining() too large immediately after Allow(): %v", r)
	}
}
