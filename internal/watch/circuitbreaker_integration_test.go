package watch_test

import (
	"errors"
	"testing"
	"time"

	"github.com/user/portwatch/internal/watch"
)

func TestCircuitBreakerIntegration(t *testing.T) {
	policy := watch.CircuitBreakerPolicy{
		MaxFailures:  3,
		OpenDuration: 30 * time.Millisecond,
	}
	cb := watch.NewCircuitBreaker(policy)

	callCount := 0
	operation := func() error {
		callCount++
		return errors.New("service unavailable")
	}

	// Drive failures until open.
	for i := 0; i < 3; i++ {
		if err := cb.Allow(); err != nil {
			t.Fatalf("unexpected block before open: %v", err)
		}
		cb.RecordFailure()
		_ = operation()
	}

	if cb.State() != "open" {
		t.Fatalf("expected open state, got %s", cb.State())
	}

	// Calls should be blocked.
	if err := cb.Allow(); !errors.Is(err, watch.ErrCircuitOpen) {
		t.Fatalf("expected ErrCircuitOpen, got %v", err)
	}

	// Wait for recovery window.
	time.Sleep(50 * time.Millisecond)

	if err := cb.Allow(); err != nil {
		t.Fatalf("expected half-open to allow, got %v", err)
	}

	if cb.State() != "half-open" {
		t.Fatalf("expected half-open state after recovery window, got %s", cb.State())
	}

	// Simulate success.
	cb.RecordSuccess()
	if cb.State() != "closed" {
		t.Fatalf("expected closed after recovery, got %s", cb.State())
	}

	if callCount != 3 {
		t.Errorf("expected 3 calls, got %d", callCount)
	}
}
