package watch

import (
	"sync"
	"time"
)

// Throttle limits how frequently a repeated action can fire.
// After a trigger, subsequent calls within the cooldown window are suppressed.
type Throttle struct {
	mu       sync.Mutex
	cooldown time.Duration
	lastFire time.Time
	now      func() time.Time
}

// NewThrottle creates a Throttle with the given cooldown duration.
func NewThrottle(cooldown time.Duration) *Throttle {
	return &Throttle{
		cooldown: cooldown,
		now:      time.Now,
	}
}

// Allow returns true if enough time has passed since the last allowed call.
// If allowed, it records the current time as the last fire time.
func (t *Throttle) Allow() bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := t.now()
	if t.lastFire.IsZero() || now.Sub(t.lastFire) >= t.cooldown {
		t.lastFire = now
		return true
	}
	return false
}

// Reset clears the last fire time, allowing the next call to Always succeed.
func (t *Throttle) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.lastFire = time.Time{}
}

// Remaining returns how long until the throttle will allow another call.
// Returns zero if the throttle is ready.
func (t *Throttle) Remaining() time.Duration {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.lastFire.IsZero() {
		return 0
	}
	elapsed := t.now().Sub(t.lastFire)
	if elapsed >= t.cooldown {
		return 0
	}
	return t.cooldown - elapsed
}
