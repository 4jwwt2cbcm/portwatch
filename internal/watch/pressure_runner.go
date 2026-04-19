package watch

import "context"

// PressureRunner wraps a function and skips execution when the
// PressureTracker reports high pressure.
type PressureRunner struct {
	tracker *PressureTracker
	fn      func(ctx context.Context) error
}

// NewPressureRunner creates a PressureRunner. A nil tracker defaults to
// a tracker with default policy (never high on init).
func NewPressureRunner(pt *PressureTracker, fn func(ctx context.Context) error) *PressureRunner {
	if pt == nil {
		pt = NewPressureTracker(DefaultPressurePolicy())
	}
	if fn == nil {
		fn = func(_ context.Context) error { return nil }
	}
	return &PressureRunner{tracker: pt, fn: fn}
}

// Run executes the wrapped function only when pressure is not high.
// Returns ErrSkipped when skipped due to high pressure.
func (pr *PressureRunner) Run(ctx context.Context) error {
	if pr.tracker.High() {
		return ErrSkipped
	}
	return pr.fn(ctx)
}

// ErrSkipped is returned by PressureRunner when the call is suppressed.
type skippedError struct{}

func (skippedError) Error() string { return "skipped: high pressure" }

// ErrSkipped sentinel.
var ErrSkipped error = skippedError{}
