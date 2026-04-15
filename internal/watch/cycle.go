package watch

import (
	"context"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/metrics"
	"github.com/user/portwatch/internal/state"
)

// CycleRunner executes a single scan-diff-alert-persist cycle.
type CycleRunner struct {
	manager   *state.Manager
	dispatcher *alert.Dispatcher
	collector  *metrics.Collector
	logger    *Logger
}

// NewCycleRunner creates a CycleRunner wiring together state, alerting, and metrics.
func NewCycleRunner(
	manager *state.Manager,
	dispatcher *alert.Dispatcher,
	collector *metrics.Collector,
	logger *Logger,
) *CycleRunner {
	return &CycleRunner{
		manager:    manager,
		dispatcher: dispatcher,
		collector:  collector,
		logger:     logger,
	}
}

// Run performs one full scan cycle. It scans ports, computes the diff against
// the previous state, dispatches alerts for changes, and persists the new state.
// Returns early if ctx is cancelled.
func (c *CycleRunner) Run(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	diff, ports, err := c.manager.Cycle(ctx)
	if err != nil {
		c.logger.Error("scan cycle failed", err)
		return err
	}

	c.collector.RecordScan(ports, diff)
	c.dispatcher.Dispatch(diff)

	if len(diff.Added) > 0 || len(diff.Removed) > 0 {
		c.logger.Info("port changes detected: +%d/-%d", len(diff.Added), len(diff.Removed))
	}

	return nil
}
