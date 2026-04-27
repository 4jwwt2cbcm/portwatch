package watch

import (
	"sync"
	"time"
)

// TriggerPolicy configures Trigger behaviour.
type TriggerPolicy struct {
	// MinInterval is the minimum time that must elapse between consecutive
	// triggers. Zero means no minimum.
	MinInterval time.Duration
}

// DefaultTriggerPolicy returns sensible defaults.
func DefaultTriggerPolicy() TriggerPolicy {
	return TriggerPolicy{
		MinInterval: 0,
	}
}

// Trigger is a one-shot or repeatable manual signal that can be fired and
// waited on. It respects an optional minimum interval between firings.
type Trigger struct {
	policy TriggerPolicy
	mu      sync.Mutex
	ch      chan struct{}
	lastFired time.Time
	count   int
}

// NewTrigger creates a new Trigger with the given policy.
func NewTrigger(p TriggerPolicy) *Trigger {
	return &Trigger{
		policy: p,
		ch:     make(chan struct{}, 1),
	}
}

// Fire signals the trigger. Returns true if the signal was delivered, false if
// suppressed by the minimum interval or a pending signal already exists.
func (t *Trigger) Fire() bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.policy.MinInterval > 0 && !t.lastFired.IsZero() {
		if time.Since(t.lastFired) < t.policy.MinInterval {
			return false
		}
	}

	select {
	case t.ch <- struct{}{}:
		t.lastFired = time.Now()
		t.count++
		return true
	default:
		return false
	}
}

// C returns the channel that receives trigger signals.
func (t *Trigger) C() <-chan struct{} {
	return t.ch
}

// Count returns the total number of times the trigger has fired.
func (t *Trigger) Count() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.count
}

// Reset clears the count and last-fired time.
func (t *Trigger) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.count = 0
	t.lastFired = time.Time{}
	// drain any pending signal
	select {
	case <-t.ch:
	default:
	}
}
