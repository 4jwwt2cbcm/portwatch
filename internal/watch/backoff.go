package watch

import "time"

// BackoffPolicy defines the parameters for exponential backoff.
type BackoffPolicy struct {
	Initial time.Duration
	Max     time.Duration
	Factor  float64
}

// DefaultBackoffPolicy returns a sensible default backoff policy.
func DefaultBackoffPolicy() BackoffPolicy {
	return BackoffPolicy{
		Initial: 500 * time.Millisecond,
		Max:     30 * time.Second,
		Factor:  2.0,
	}
}

// Backoff tracks exponential backoff state.
type Backoff struct {
	policy  BackoffPolicy
	current time.Duration
}

// NewBackoff creates a new Backoff with the given policy.
func NewBackoff(policy BackoffPolicy) *Backoff {
	return &Backoff{
		policy:  policy,
		current: policy.Initial,
	}
}

// Next returns the current wait duration and advances the backoff.
func (b *Backoff) Next() time.Duration {
	d := b.current
	next := time.Duration(float64(b.current) * b.policy.Factor)
	if next > b.policy.Max {
		next = b.policy.Max
	}
	b.current = next
	return d
}

// Reset restores the backoff to its initial duration.
func (b *Backoff) Reset() {
	b.current = b.policy.Initial
}

// Current returns the current wait duration without advancing.
func (b *Backoff) Current() time.Duration {
	return b.current
}
