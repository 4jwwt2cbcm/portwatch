package watch

import (
	"context"
	"log"
	"time"

	"github.com/user/portwatch/internal/metrics"
	"github.com/user/portwatch/internal/state"
)

// Watcher orchestrates periodic port scanning cycles.
type Watcher struct {
	manager  *state.Manager
	collector *metrics.Collector
	interval time.Duration
	logger   *log.Logger
}

// NewWatcher creates a Watcher with the given manager, collector, interval, and logger.
func NewWatcher(manager *state.Manager, collector *metrics.Collector, interval time.Duration, logger *log.Logger) *Watcher {
	if logger == nil {
		logger = log.Default()
	}
	return &Watcher{
		manager:   manager,
		collector: collector,
		interval:  interval,
		logger:    logger,
	}
}

// Run starts the watch loop, executing a scan cycle at each tick until ctx is cancelled.
func (w *Watcher) Run(ctx context.Context) error {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	w.logger.Printf("watcher started with interval %s", w.interval)

	for {
		select {
		case <-ctx.Done():
			w.logger.Println("watcher stopped")
			return ctx.Err()
		case <-ticker.C:
			if err := w.tick(ctx); err != nil {
				w.logger.Printf("scan cycle error: %v", err)
			}
		}
	}
}

func (w *Watcher) tick(ctx context.Context) error {
	diff, ports, err := w.manager.Cycle(ctx)
	if err != nil {
		return err
	}
	w.collector.RecordScan(ports, diff)
	return nil
}
