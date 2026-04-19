package watch

import "context"

// TeeRunner runs a function and emits the result to a Tee.
type TeeRunner[T any] struct {
	tee *Tee[T]
	fn  func(context.Context) (T, error)
}

// NewTeeRunner creates a TeeRunner that calls fn and forwards results to tee.
// If tee is nil a no-op Tee is used. If fn is nil it returns the zero value.
func NewTeeRunner[T any](tee *Tee[T], fn func(context.Context) (T, error)) *TeeRunner[T] {
	if tee == nil {
		tee = NewTee[T]()
	}
	if fn == nil {
		fn = func(_ context.Context) (T, error) { var z T; return z, nil }
	}
	return &TeeRunner[T]{tee: tee, fn: fn}
}

// Run executes the function and emits its result to all sinks.
func (r *TeeRunner[T]) Run(ctx context.Context) (T, error) {
	v, err := r.fn(ctx)
	r.tee.Emit(v, err)
	return v, err
}
