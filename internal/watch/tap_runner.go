package watch

// TapRunner wraps an inner runner function and records each invocation
// result via a Tap, allowing callers to inspect execution history.
type TapRunner struct {
	tap *Tap[error]
	fn  func() error
}

// NewTapRunner creates a TapRunner that records errors (or nil) from fn.
// cap controls how many results are retained; zero defaults to 64.
func NewTapRunner(fn func() error, cap int) *TapRunner {
	if fn == nil {
		fn = func() error { return nil }
	}
	return &TapRunner{
		tap: NewTap[error](cap, nil),
		fn:  fn,
	}
}

// Run executes the wrapped function and records the result.
func (r *TapRunner) Run() error {
	err := r.fn()
	r.tap.Record(err)
	return err
}

// Tap returns the underlying Tap for inspection.
func (r *TapRunner) Tap() *Tap[error] {
	return r.tap
}
