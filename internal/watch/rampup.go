package watch

import (
	"sync"
	"time"
)

// DefaultRampUpPolicy returns a sensible default ramp-up policy.
func DefaultRampUpPolicy() RampUpPolicy {
	return RampUpPolicy{
		Steps:    5,
		Initial:  200 * time.Millisecond,
		Target:   2 * time.Second,
	}
}

// RampUpPolicy controls how the ramp-up interval grows.
type RampUpPolicy struct {
	Steps   int
	Initial time.Duration
	Target  time.Duration
}

// RampUp gradually increases an interval from Initial to Target over Steps ticks.
type RampUp struct {
	mu      sync.Mutex
	policy  RampUpPolicy
	step    int
	current time.Duration
}

// NewRampUp creates a new RampUp with the given policy.
func NewRampUp(p RampUpPolicy) *RampUp {
	if p.Steps <= 0 {
		p.Steps = 1
	}
	if p.Initial <= 0 {
		p.Initial = 100 * time.Millisecond
	}
	if p.Target <= p.Initial {
		p.Target = p.Initial
	}
	return &RampUp{policy: p, current: p.Initial}
}

// Next returns the current interval and advances to the next step.
func (r *RampUp) Next() time.Duration {
	r.mu.Lock()
	defer r.mu.Unlock()
	v := r.current
	if r.step < r.policy.Steps {
		r.step++
		span := float64(r.policy.Target - r.policy.Initial)
		r.current = r.policy.Initial + time.Duration(span*float64(r.step)/float64(r.policy.Steps))
	}
	return v
}

// Done reports whether the ramp-up has reached the target.
func (r *RampUp) Done() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.step >= r.policy.Steps
}

// Reset restarts the ramp-up from the initial interval.
func (r *RampUp) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.step = 0
	r.current = r.policy.Initial
}
