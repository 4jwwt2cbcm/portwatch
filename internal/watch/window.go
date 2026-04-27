package watch

import (
	"sync"
	"time"
)

// DefaultWindowPolicy returns a WindowPolicy with sensible defaults.
func DefaultWindowPolicy() WindowPolicy {
	return WindowPolicy{
		Size:     10 * time.Second,
		MaxItems: 100,
	}
}

// WindowPolicy configures a RollingWindow.
type WindowPolicy struct {
	Size     time.Duration
	MaxItems int
}

// windowEntry is a timestamped value stored in the rolling window.
type windowEntry[T any] struct {
	at    time.Time
	value T
}

// RollingWindow is a generic time-bounded sliding window that retains
// at most MaxItems entries within the configured Size duration.
type RollingWindow[T any] struct {
	mu     sync.Mutex
	policy WindowPolicy
	items  []windowEntry[T]
	now    func() time.Time
}

// NewRollingWindow creates a RollingWindow with the given policy.
// A zero policy falls back to DefaultWindowPolicy.
func NewRollingWindow[T any](p WindowPolicy) *RollingWindow[T] {
	if p.Size <= 0 {
		p.Size = DefaultWindowPolicy().Size
	}
	if p.MaxItems <= 0 {
		p.MaxItems = DefaultWindowPolicy().MaxItems
	}
	return &RollingWindow[T]{
		policy: p,
		now:    time.Now,
	}
}

// Add inserts a value into the window, evicting entries outside the window
// and capping storage at MaxItems.
func (w *RollingWindow[T]) Add(v T) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.evict()
	if len(w.items) >= w.policy.MaxItems {
		w.items = w.items[1:]
	}
	w.items = append(w.items, windowEntry[T]{at: w.now(), value: v})
}

// Snapshot returns a copy of all values currently within the window.
func (w *RollingWindow[T]) Snapshot() []T {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.evict()
	out := make([]T, len(w.items))
	for i, e := range w.items {
		out[i] = e.value
	}
	return out
}

// Len returns the number of entries currently in the window.
func (w *RollingWindow[T]) Len() int {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.evict()
	return len(w.items)
}

// Clear removes all entries from the window.
func (w *RollingWindow[T]) Clear() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.items = nil
}

// evict removes entries older than the window size. Must be called with mu held.
func (w *RollingWindow[T]) evict() {
	cutoff := w.now().Add(-w.policy.Size)
	i := 0
	for i < len(w.items) && w.items[i].at.Before(cutoff) {
		i++
	}
	w.items = w.items[i:]
}
