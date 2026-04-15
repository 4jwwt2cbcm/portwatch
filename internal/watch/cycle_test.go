package watch

import (
	"context"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/metrics"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/state"
)

func makeCycleRunner(t *testing.T) (*CycleRunner, *metrics.Collector) {
	t.Helper()

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

	logger := NewLogger(nil, nil)

	return NewCycleRunner(mgr, dispatcher, collector, logger), collector
}

func TestCycleRunnerRunCompletesSuccessfully(t *testing.T) {
	cr, _ := makeCycleRunner(t)
	ctx := context.Background()

	if err := cr.Run(ctx); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestCycleRunnerRecordsMetrics(t *testing.T) {
	cr, collector := makeCycleRunner(t)
	ctx := context.Background()

	if err := cr.Run(ctx); err != nil {
		t.Fatalf("Run: %v", err)
	}

	snap := collector.Snapshot()
	if snap.TotalScans != 1 {
		t.Errorf("expected TotalScans=1, got %d", snap.TotalScans)
	}
}

func TestCycleRunnerCancelledContextReturnsError(t *testing.T) {
	cr, _ := makeCycleRunner(t)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if err := cr.Run(ctx); err == nil {
		t.Fatal("expected error on cancelled context, got nil")
	}
}
