package state_test

import (
	"testing"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/state"
)

func TestCycleFirstRunNoDiff(t *testing.T) {
	sc := scanner.NewScanner()
	store := state.NewStore(tempPath(t))
	mgr := state.NewManager(store, sc)

	// First cycle: no previous state, so no added/removed ports expected
	// (all currently open ports appear as "added" vs empty baseline).
	diff, err := mgr.Cycle([]int{}, "tcp")
	if err != nil {
		t.Fatalf("Cycle() error: %v", err)
	}
	// With an empty port list to scan, nothing should be added or removed.
	if len(diff.Added) != 0 {
		t.Errorf("expected no added ports, got %d", len(diff.Added))
	}
	if len(diff.Removed) != 0 {
		t.Errorf("expected no removed ports, got %d", len(diff.Removed))
	}
}

func TestCycleDetectsStateChange(t *testing.T) {
	path := tempPath(t)
	store := state.NewStore(path)

	// Seed the store with a previously open port.
	initial := []scanner.PortState{
		{Port: 9999, Proto: "tcp", State: "open"},
	}
	if err := store.Save(initial); err != nil {
		t.Fatalf("seed Save() error: %v", err)
	}

	sc := scanner.NewScanner()
	mgr := state.NewManager(store, sc)

	// Scan an empty port list — port 9999 should appear as removed.
	diff, err := mgr.Cycle([]int{}, "tcp")
	if err != nil {
		t.Fatalf("Cycle() error: %v", err)
	}

	if len(diff.Removed) != 1 || diff.Removed[0].Port != 9999 {
		t.Errorf("expected port 9999 removed, got removed=%+v", diff.Removed)
	}
}

func TestCyclePersistsNewState(t *testing.T) {
	path := tempPath(t)
	store := state.NewStore(path)
	sc := scanner.NewScanner()
	mgr := state.NewManager(store, sc)

	_, err := mgr.Cycle([]int{}, "tcp")
	if err != nil {
		t.Fatalf("Cycle() error: %v", err)
	}

	snap, err := store.Load()
	if err != nil {
		t.Fatalf("Load() after Cycle() error: %v", err)
	}
	if snap.Timestamp.IsZero() {
		t.Error("expected persisted snapshot to have a timestamp")
	}
}
