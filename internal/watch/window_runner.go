package watch

// WindowRunner wraps a function and records its return value (or error)
// into a RollingWindow of results on each invocation.
type WindowRunner[T any] struct {
	win *RollingWindow[T]
	fn  func() (T, error)
}

// NewWindowRunner creates a WindowRunner that stores results in win.
// If win is nil a default window is created. If fn is nil a no-op is used.
func NewWindowRunner[T any](win *RollingWindow[T], fn func() (T, error)) *WindowRunner[T] {
	if win == nil {
		win = NewRollingWindow[T](DefaultWindowPolicy())
	}
	if fn == nil {
		var zero T
		fn = func() (T, error) { return zero, nil }
	}
	return &WindowRunner[T]{win: win, fn: fn}
}

// Run invokes the wrapped function. On success the result is stored in the
// window. The error (if any) is returned to the caller unchanged.
func (r *WindowRunner[T]) Run() (T, error) {
	v, err := r.fn()
	if err == nil {
		r.win.Add(v)
	}
	return v, err
}

// Window returns the underlying RollingWindow.
func (r *WindowRunner[T]) Window() *RollingWindow[T] {
	return r.win
}
