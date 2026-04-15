package watch_test

import (
	"context"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/metrics"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/state"
	"github.com/user/portwatch/internal/watch"
)

// TestCycleRunnerIntegration verifies that multiple consecutive cycles
// accumulate scan counts correctly and do not error.
func TestCycleRunnerIntegration(t *testing.T) {
	dir := t.TempDir()
	store, err := state.NewStore(dir + "/state.json")
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}

	sc := scanner.NewScanner(scanner.NewFilter(nil))
	mgr := state.NewManager(store, sc)

	collector := metrics.NewCollector()
	notifier := alert.NewNotifier(nil)
	fmt := alert.NewFormatter("text")
	dispatcher := alert.NewDispatcher(notifier, fmt)
	logger := watch.NewLogger(nil, nil)

	cr := watch.NewCycleRunner(mgr, dispatcher, collector, logger)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	const cycles = 3
	for i := 0; i < cycles; i++ {
		if err := cr.Run(ctx); err != nil {
			t.Fatalf("cycle %d failed: %v", i+1, err)
		}
	}

	snap := collector.Snapshot()
	if snap.TotalScans != cycles {
		t.Errorf("expected TotalScans=%d, got %d", cycles, snap.TotalScans)
	}
}
