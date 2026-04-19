package watch

import "sync"

// Tee splits the result of a function call to multiple sinks.
type Tee[T any] struct {
	mu    sync.Mutex
	sinks []func(T, error)
}

// NewTee creates a new Tee with the given sinks.
func NewTee[T any](sinks ...func(T, error)) *Tee[T] {
	return &Tee[T]{sinks: sinks}
}

// Add appends a sink to the Tee.
func (t *Tee[T]) Add(fn func(T, error)) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.sinks = append(t.sinks, fn)
}

// Emit calls all registered sinks with the given value and error.
func (t *Tee[T]) Emit(v T, err error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	for _, s := range t.sinks {
		s(v, err)
	}
}

// Wrap returns a function that calls fn and emits its result to all sinks.
func (t *Tee[T]) Wrap(fn func() (T, error)) func() (T, error) {
	return func() (T, error) {
		v, err := fn()
		t.Emit(v, err)
		return v, err
	}
}

// Count returns the number of registered sinks.
func (t *Tee[T]) Count() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return len(t.sinks)
}

// Clear removes all sinks.
func (t *Tee[T]) Clear() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.sinks = nil
}
