package watch

import (
	"testing"
	"time"
)

func makeCircuitBreaker(maxFailures int, openDuration time.Duration) *CircuitBreaker {
	return NewCircuitBreaker(CircuitBreakerPolicy{
		MaxFailures:  maxFailures,
		OpenDuration: openDuration,
	})
}

func TestCircuitBreakerDefaultPolicy(t *testing.T) {
	p := DefaultCircuitBreakerPolicy()
	if p.MaxFailures != 5 {
		t.Errorf("expected MaxFailures=5, got %d", p.MaxFailures)
	}
	if p.OpenDuration != 30*time.Second {
		t.Errorf("expected OpenDuration=30s, got %v", p.OpenDuration)
	}
}

func TestCircuitBreakerClosedByDefault(t *testing.T) {
	cb := makeCircuitBreaker(3, time.Second)
	if err := cb.Allow(); err != nil {
		t.Errorf("expected nil, got %v", err)
	}
	if cb.State() != "closed" {
		t.Errorf("expected closed, got %s", cb.State())
	}
}

func TestCircuitBreakerOpensAfterMaxFailures(t *testing.T) {
	cb := makeCircuitBreaker(3, time.Minute)
	for i := 0; i < 3; i++ {
		cb.RecordFailure()
	}
	if cb.State() != "open" {
		t.Errorf("expected open, got %s", cb.State())
	}
	if err := cb.Allow(); err != ErrCircuitOpen {
		t.Errorf("expected ErrCircuitOpen, got %v", err)
	}
}

func TestCircuitBreakerSuccessResetsClosed(t *testing.T) {
	cb := makeCircuitBreaker(2, time.Minute)
	cb.RecordFailure()
	cb.RecordFailure()
	cb.RecordSuccess()
	if cb.State() != "closed" {
		t.Errorf("expected closed after success, got %s", cb.State())
	}
	if err := cb.Allow(); err != nil {
		t.Errorf("expected nil after reset, got %v", err)
	}
}

func TestCircuitBreakerHalfOpenAfterDuration(t *testing.T) {
	cb := makeCircuitBreaker(1, 10*time.Millisecond)
	cb.RecordFailure()
	time.Sleep(20 * time.Millisecond)
	if err := cb.Allow(); err != nil {
		t.Errorf("expected nil in half-open, got %v", err)
	}
	if cb.State() != "half-open" {
		t.Errorf("expected half-open, got %s", cb.State())
	}
}

func TestCircuitBreakerDoesNotOpenBeforeMax(t *testing.T) {
	cb := makeCircuitBreaker(5, time.Minute)
	for i := 0; i < 4; i++ {
		cb.RecordFailure()
	}
	if cb.State() != "closed" {
		t.Errorf("expected closed, got %s", cb.State())
	}
}

func TestCircuitBreakerHalfOpenFailureReopens(t *testing.T) {
	cb := makeCircuitBreaker(1, 10*time.Millisecond)
	cb.RecordFailure()
	time.Sleep(20 * time.Millisecond)
	// Transition to half-open
	if err := cb.Allow(); err != nil {
		t.Fatalf("expected nil in half-open, got %v", err)
	}
	// A failure in half-open should reopen the circuit
	cb.RecordFailure()
	if cb.State() != "open" {
		t.Errorf("expected open after half-open failure, got %s", cb.State())
	}
	if err := cb.Allow(); err != ErrCircuitOpen {
		t.Errorf("expected ErrCircuitOpen after reopen, got %v", err)
	}
}
