package watch

import (
	"context"
	"os"
)

// SignalRunner wraps a run function and cancels it gracefully on OS signal.
type SignalRunner struct {
	handler *SignalHandler
	logger  *Logger
}

// NewSignalRunner creates a SignalRunner with the given signals.
func NewSignalRunner(logger *Logger, sigs ...os.Signal) *SignalRunner {
	if logger == nil {
		logger = NewLogger(nil)
	}
	return &SignalRunner{
		handler: NewSignalHandler(sigs...),
		logger:  logger,
	}
}

// Run executes fn within a context that is cancelled on signal receipt.
// The provided parent context is also respected.
func (r *SignalRunner) Run(parent context.Context, fn func(ctx context.Context) error) error {
	ctx, cancel := r.handler.WithCancel(parent)
	defer cancel()

	r.logger.Info("signal runner started")
	err := fn(ctx)
	if err != nil && ctx.Err() == nil {
		r.logger.Error("run function returned error", err)
		return err
	}
	r.logger.Info("signal runner stopped")
	return nil
}
