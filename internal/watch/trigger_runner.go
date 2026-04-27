package watch

import (
	"context"
	"errors"
)

// ErrTriggerNotFired is returned by TriggerRunner when the trigger has not
// been fired and the context is done.
var ErrTriggerNotFired = errors.New("trigger: not fired")

// TriggerRunner waits for a Trigger signal before invoking a function.
// It blocks until the trigger fires or the context is cancelled.
type TriggerRunner struct {
	trigger *Trigger
	fn      func(ctx context.Context) error
}

// NewTriggerRunner creates a TriggerRunner. If trigger is nil a no-op trigger
// is used. If fn is nil it defaults to a no-op.
func NewTriggerRunner(trigger *Trigger, fn func(ctx context.Context) error) *TriggerRunner {
	if trigger == nil {
		trigger = NewTrigger(DefaultTriggerPolicy())
	}
	if fn == nil {
		fn = func(_ context.Context) error { return nil }
	}
	return &TriggerRunner{trigger: trigger, fn: fn}
}

// RunOnce blocks until the trigger fires or ctx is cancelled, then calls fn
// once. Returns ErrTriggerNotFired if the context is cancelled before the
// trigger fires.
func (r *TriggerRunner) RunOnce(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ErrTriggerNotFired
	case <-r.trigger.C():
		return r.fn(ctx)
	}
}

// RunLoop continuously waits for trigger signals and calls fn until the
// context is cancelled or fn returns an error.
func (r *TriggerRunner) RunLoop(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-r.trigger.C():
			if err := r.fn(ctx); err != nil {
				return err
			}
		}
	}
}
