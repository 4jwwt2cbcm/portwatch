package watch

import (
	"context"
	"errors"
)

// ErrFuseBlown is returned when the fuse has tripped.
var ErrFuseBlown = errors.New("fuse blown")

// FuseRunner wraps a function and records errors into a Fuse.
// If the fuse is blown, the wrapped function is not called.
type FuseRunner struct {
	fuse *Fuse
	fn   func(context.Context) error
}

// NewFuseRunner returns a FuseRunner. If fuse is nil a default is created.
// If fn is nil a noop is used.
func NewFuseRunner(fuse *Fuse, fn func(context.Context) error) *FuseRunner {
	if fuse == nil {
		fuse = NewFuse(DefaultFusePolicy())
	}
	if fn == nil {
		fn = func(context.Context) error { return nil }
	}
	return &FuseRunner{fuse: fuse, fn: fn}
}

// Run executes the wrapped function unless the fuse is blown.
// Errors returned by fn are recorded against the fuse.
func (r *FuseRunner) Run(ctx context.Context) error {
	if r.fuse.Blown() {
		return ErrFuseBlown
	}
	if err := r.fn(ctx); err != nil {
		r.fuse.Record()
		return err
	}
	return nil
}

// Fuse exposes the underlying fuse for inspection or reset.
func (r *FuseRunner) Fuse() *Fuse {
	return r.fuse
}
