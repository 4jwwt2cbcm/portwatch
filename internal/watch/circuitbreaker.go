package watch

import (
	"errors"
	"sync"
	"time"
)

// ErrCircuitOpen is returned when the circuit breaker is open.
var ErrCircuitOpen = errors.New("circuit breaker is open")

type cbState int

const (
	cbClosed cbState = iota
	cbOpen
	cbHalfOpen
)

// CircuitBreakerPolicy configures the circuit breaker.
type CircuitBreakerPolicy struct {
	MaxFailures int
	OpenDuration time.Duration
}

// DefaultCircuitBreakerPolicy returns sensible defaults.
func DefaultCircuitBreakerPolicy() CircuitBreakerPolicy {
	return CircuitBreakerPolicy{
		MaxFailures:  5,
		OpenDuration: 30 * time.Second,
	}
}

// CircuitBreaker prevents repeated calls when a dependency is failing.
type CircuitBreaker struct {
	mu       sync.Mutex
	policy   CircuitBreakerPolicy
	state    cbState
	failures int
	openedAt time.Time
}

// NewCircuitBreaker creates a new CircuitBreaker with the given policy.
func NewCircuitBreaker(policy CircuitBreakerPolicy) *CircuitBreaker {
	return &CircuitBreaker{policy: policy}
}

// Allow returns nil if the call should proceed, or ErrCircuitOpen if not.
func (cb *CircuitBreaker) Allow() error {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	switch cb.state {
	case cbOpen:
		if time.Since(cb.openedAt) >= cb.policy.OpenDuration {
			cb.state = cbHalfOpen
			return nil
		}
		return ErrCircuitOpen
	}
	return nil
}

// RecordSuccess resets the circuit breaker to closed.
func (cb *CircuitBreaker) RecordSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.failures = 0
	cb.state = cbClosed
}

// RecordFailure records a failure and may open the circuit.
func (cb *CircuitBreaker) RecordFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.failures++
	if cb.failures >= cb.policy.MaxFailures {
		cb.state = cbOpen
		cb.openedAt = time.Now()
	}
}

// State returns the current state label for observability.
func (cb *CircuitBreaker) State() string {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	switch cb.state {
	case cbOpen:
		return "open"
	case cbHalfOpen:
		return "half-open"
	default:
		return "closed"
	}
}
