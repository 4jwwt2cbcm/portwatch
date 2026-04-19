package watch

import "sync"

// Tap is a pass-through observer that records values flowing through a pipeline
// without modifying them. Useful for debugging and metrics collection.
type Tap[T any] struct {
	mu       sync.Mutex
	recorded []T
	cap      int
	onTap    func(T)
}

// NewTap creates a Tap with the given capacity. Zero cap defaults to 64.
// An optional onTap callback is invoked synchronously for each value.
func NewTap[T any](cap int, onTap func(T)) *Tap[T] {
	if cap <= 0 {
		cap = 64
	}
	return &Tap[T]{cap: cap, onTap: onTap}
}

// Record passes value through the tap, storing it and invoking the callback.
func (t *Tap[T]) Record(v T) T {
	t.mu.Lock()
	if len(t.recorded) < t.cap {
		t.recorded = append(t.recorded, v)
	}
	t.mu.Unlock()
	if t.onTap != nil {
		t.onTap(v)
	}
	return v
}

// Snapshot returns a copy of all recorded values.
func (t *Tap[T]) Snapshot() []T {
	t.mu.Lock()
	defer t.mu.Unlock()
	out := make([]T, len(t.recorded))
	copy(out, t.recorded)
	return out
}

// Len returns the number of recorded values.
func (t *Tap[T]) Len() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return len(t.recorded)
}

// Clear resets the recorded values.
func (t *Tap[T]) Clear() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.recorded = t.recorded[:0]
}
