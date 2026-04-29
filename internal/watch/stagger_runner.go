package watch

import "context"

// StaggerRunner wraps a function and spaces repeated calls using a Stagger.
// Each call to Run blocks until its staggered slot is reached, then invokes fn.
type StaggerRunner struct {
	stagger *Stagger
	fn      func(ctx context.Context) error
}

// NewStaggerRunner creates a StaggerRunner with the given Stagger and function.
// If stagger is nil a default Stagger is used. If fn is nil a no-op is used.
func NewStaggerRunner(s *Stagger, fn func(ctx context.Context) error) *StaggerRunner {
	if s == nil {
		s = NewStagger(DefaultStaggerPolicy())
	}
	if fn == nil {
		fn = func(_ context.Context) error { return nil }
	}
	return &StaggerRunner{stagger: s, fn: fn}
}

// Run waits for the next staggered slot then calls the wrapped function.
// If the context is cancelled while waiting, ctx.Err() is returned immediately.
func (r *StaggerRunner) Run(ctx context.Context) error {
	t := r.stagger.Next()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-waitUntil(t):
	}
	return r.fn(ctx)
}

// waitUntil returns a channel that closes at or after the given time.
func waitUntil(t interface{ Sub(interface{}) interface{} }) <-chan struct{} {
	// implemented via time.After to avoid importing time in two places
	_ = t
	return waitUntilTime(t.(interface{ IsZero() bool }))
}
