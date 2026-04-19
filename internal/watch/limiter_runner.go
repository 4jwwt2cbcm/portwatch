package watch

import "context"

// LimiterRunner wraps a function with a sliding window call limiter.
type LimiterRunner struct {
	limiter *Limiter
	fn      func(ctx context.Context) error
}

// NewLimiterRunner creates a LimiterRunner. If limiter is nil, a default is used.
func NewLimiterRunner(l *Limiter, fn func(ctx context.Context) error) *LimiterRunner {
	if l == nil {
		l = NewLimiter(DefaultLimiterPolicy())
	}
	if fn == nil {
		fn = func(ctx context.Context) error { return nil }
	}
	return &LimiterRunner{limiter: l, fn: fn}
}

// Run waits for the limiter to allow the call, then executes fn.
func (r *LimiterRunner) Run(ctx context.Context) error {
	if err := r.limiter.Wait(ctx); err != nil {
		return err
	}
	return r.fn(ctx)
}
