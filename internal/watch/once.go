package watch

import "sync"

// OnceRunner wraps a function so it only executes once, regardless of how
// many times Run is called. Subsequent calls return the result of the first.
type OnceRunner struct {
	mu   sync.Mutex
	once sync.Once
	fn   func() error
	err  error
	ran  bool
}

// NewOnceRunner returns a OnceRunner that will call fn at most once.
func NewOnceRunner(fn func() error) *OnceRunner {
	if fn == nil {
		fn = func() error { return nil }
	}
	return &OnceRunner{fn: fn}
}

// Run executes the wrapped function the first time it is called.
// All subsequent calls return the same error without invoking fn again.
func (o *OnceRunner) Run() error {
	o.once.Do(func() {
		o.mu.Lock()
		defer o.mu.Unlock()
		o.err = o.fn()
		o.ran = true
	})
	o.mu.Lock()
	defer o.mu.Unlock()
	return o.err
}

// HasRun reports whether the function has been executed at least once.
func (o *OnceRunner) HasRun() bool {
	o.mu.Lock()
	defer o.mu.Unlock()
	return o.ran
}

// Reset allows the function to be called again on the next Run invocation.
func (o *OnceRunner) Reset() {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.once = sync.Once{}
	o.ran = false
	o.err = nil
}
