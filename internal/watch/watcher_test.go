package watch

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/user/portwatch/internal/metrics"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/state"
)

func makeWatcher(t *testing.T) *Watcher {
	t.Helper()

	dir := t.TempDir()
	store, err := state.NewStore(dir)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}

	sc := scanner.NewScanner(nil)
	collector := metrics.NewCollector()
	dispatcher := nil // alerts not under test here
	_ = dispatcher

	manager := state.NewManager(store, sc, nil)
	logger := log.New(os.Stderr, "test: ", 0)
	return NewWatcher(manager, collector, 50*time.Millisecond, logger)
}

func TestWatcherRunStopsOnContextCancel(t *testing.T) {
	w := makeWatcher(t)

	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan error, 1)
	go func() {
		done <- w.Run(ctx)
	}()

	time.Sleep(120 * time.Millisecond)
	cancel()

	select {
	case err := <-done:
		if err != context.Canceled {
			t.Errorf("expected context.Canceled, got %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("watcher did not stop after context cancel")
	}
}

func TestNewWatcherDefaultLogger(t *testing.T) {
	dir := t.TempDir()
	store, _ := state.NewStore(dir)
	sc := scanner.NewScanner(nil)
	manager := state.NewManager(store, sc, nil)
	collector := metrics.NewCollector()

	w := NewWatcher(manager, collector, time.Second, nil)
	if w.logger == nil {
		t.Error("expected non-nil logger when nil is passed")
	}
}
