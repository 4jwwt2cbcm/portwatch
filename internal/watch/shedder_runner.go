package watch

import (
	"context"
	"errors"
)

// ErrShed is returned by ShedderRunner when work is dropped.
var ErrShed = errors.New("work shed: load above threshold")

// ShedderRunner wraps a function and skips execution when the Shedder
// determines the system is overloaded.
type ShedderRunner struct {
	shedder *Shedder
	fn      func(ctx context.Context) error
}

// NewShedderRunner creates a ShedderRunner.
// If shedder is nil a no-op shedder (never sheds) is used.
// If fn is nil it defaults to a no-op.
func NewShedderRunner(shedder *Shedder, fn func(ctx context.Context) error) *ShedderRunner {
	if shedder == nil {
		shedder = NewShedder(ShedderPolicy{
			MaxLoad:  1.1, // never triggers
			Window:   DefaultShedderPolicy().Window,
			Cooldown: DefaultShedderPolicy().Cooldown,
		})
	}
	if fn == nil {
		fn = func(_ context.Context) error { return nil }
	}
	return &ShedderRunner{shedder: shedder, fn: fn}
}

// Run executes fn unless the shedder decides to drop this unit of work.
// The caller is expected to call shedder.Record with an appropriate load
// observation before or after invoking Run.
func (r *ShedderRunner) Run(ctx context.Context) error {
	if r.shedder.ShouldShed() {
		return ErrShed
	}
	return r.fn(ctx)
}
