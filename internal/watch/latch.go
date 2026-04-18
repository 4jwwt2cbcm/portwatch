package watch

import "sync"

// Latch is a one-shot boolean flag that can be set once and checked many times.
// Once set, it cannot be unset. Safe for concurrent use.
type Latch struct {
	mu  sync.RWMutex
	set bool
}

// NewLatch returns a new unset Latch.
func NewLatch() *Latch {
	return &Latch{}
}

// Set marks the latch as triggered. Subsequent calls are no-ops.
// Returns true if this call was the one that set it.
func (l *Latch) Set() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.set {
		return false
	}
	l.set = true
	return true
}

// IsSet reports whether the latch has been set.
func (l *Latch) IsSet() bool {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.set
}

// Reset clears the latch, allowing it to be set again.
func (l *Latch) Reset() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.set = false
}

// SetOnce calls fn the first time Set is called on this latch.
// fn is invoked while the latch lock is NOT held.
func (l *Latch) SetOnce(fn func()) bool {
	if !l.Set() {
		return false
	}
	fn()
	return true
}
