package watch

// LatchRunner wraps a function and ensures it is only executed once,
// using a Latch to guard execution. Subsequent calls are no-ops.
type LatchRunner struct {
	latch *Latch
	fn    func() error
}

// NewLatchRunner returns a LatchRunner that will call fn at most once.
func NewLatchRunner(fn func() error) *LatchRunner {
	return &LatchRunner{
		latch: NewLatch(),
		fn:    fn,
	}
}

// Run executes the wrapped function if it has not been called before.
// Returns the function's error on first call, nil on subsequent calls.
func (r *LatchRunner) Run() error {
	var runErr error
	r.latch.SetOnce(func() {
		runErr = r.fn()
	})
	return runErr
}

// HasRun reports whether the function has been invoked.
func (r *LatchRunner) HasRun() bool {
	return r.latch.IsSet()
}

// Reset allows the function to be called again.
func (r *LatchRunner) Reset() {
	r.latch.Reset()
}
