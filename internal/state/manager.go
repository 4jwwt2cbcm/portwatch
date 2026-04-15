package state

import (
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/scanner"
)

// Manager coordinates loading previous state, running a scan, and persisting
// the new state. It returns the diff so callers can act on changes.
type Manager struct {
	store   *Store
	scanner *scanner.Scanner
}

// NewManager creates a Manager using the provided store and scanner.
func NewManager(store *Store, sc *scanner.Scanner) *Manager {
	return &Manager{store: store, scanner: sc}
}

// Cycle performs one full scan cycle:
//  1. Load previous snapshot from disk.
//  2. Run a fresh port scan.
//  3. Diff old vs new.
//  4. Persist new snapshot.
//
// Returns the diff result and any error encountered.
func (m *Manager) Cycle(ports []int, proto string) (scanner.Diff, error) {
	prev, err := m.store.Load()
	if err != nil {
		return scanner.Diff{}, err
	}

	current, err := m.scanner.Scan(ports, proto)
	if err != nil {
		return scanner.Diff{}, err
	}

	diff := scanner.Compare(prev.Ports, current)

	if err := m.store.Save(current); err != nil {
		return diff, err
	}

	return diff, nil
}
