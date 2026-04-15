package watch

import "time"

// BackoffPolicy defines how retry delays are calculated after scan errors.
type BackoffPolicy struct {
	Initial    time.Duration
	Multiplier float64
	Max        time.Duration
}

// DefaultBackoffPolicy returns a sensible exponential backoff configuration.
func DefaultBackoffPolicy() BackoffPolicy {
	return BackoffPolicy{
		Initial:    2 * time.Second,
		Multiplier: 2.0,
		Max:        60 * time.Second,
	}
}

// Backoff tracks consecutive failures and computes the next wait duration.
type Backoff struct {
	policy  BackoffPolicy
	current time.Duration
	failures int
}

// NewBackoff creates a Backoff using the given policy.
func NewBackoff(policy BackoffPolicy) *Backoff {
	return &Backoff{
		policy:  policy,
		current: policy.Initial,
	}
}

// Failure records a failed attempt and returns the duration to wait before retrying.
func (b *Backoff) Failure() time.Duration {
	b.failures++
	delay := b.current
	next := time.Duration(float64(b.current) * b.policy.Multiplier)
	if next > b.policy.Max {
		next = b.policy.Max
	}
	b.current = next
	return delay
}

// Reset clears the failure count and resets the delay to the initial value.
func (b *Backoff) Reset() {
	b.failures = 0
	b.current = b.policy.Initial
}

// Failures returns the number of consecutive failures recorded.
func (b *Backoff) Failures() int {
	return b.failures
}
