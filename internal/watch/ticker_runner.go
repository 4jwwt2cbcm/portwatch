package watch

import (
	"context"
	"time"
)

// TickerRunner calls a function on every tick of a time.Ticker until the
// context is cancelled or the callback returns an error.
type TickerRunner struct {
	interval time.Duration
	fn       func(ctx context.Context) error
	logger   *Logger
}

// NewTickerRunner creates a TickerRunner with the given interval and callback.
// If logger is nil a default stderr logger is used.
func NewTickerRunner(interval time.Duration, fn func(ctx context.Context) error, logger *Logger) *TickerRunner {
	if logger == nil {
		logger = NewLogger(nil)
	}
	return &TickerRunner{
		interval: ClampInterval(interval),
		fn:       fn,
		logger:   logger,
	}
}

// Run starts the ticker loop. It fires the callback immediately on the first
// tick and then on every subsequent interval. Returns when ctx is done or fn
// returns a non-nil error.
func (r *TickerRunner) Run(ctx context.Context) error {
	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := r.fn(ctx); err != nil {
				r.logger.Error("ticker callback error", err)
				return err
			}
		}
	}
}
