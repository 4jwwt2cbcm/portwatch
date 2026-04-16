package watch

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Snapshot holds a point-in-time view of scanned ports.
type Snapshot struct {
	Ports     []scanner.PortState
	CapturedAt time.Time
}

// SnapshotStore keeps the most recent snapshot thread-safely.
type SnapshotStore struct {
	mu       sync.RWMutex
	current  *Snapshot
}

// NewSnapshotStore returns an empty SnapshotStore.
func NewSnapshotStore() *SnapshotStore {
	return &SnapshotStore{}
}

// Set replaces the current snapshot.
func (s *SnapshotStore) Set(ports []scanner.PortState) {
	snap := &Snapshot{
		Ports:      ports,
		CapturedAt: time.Now(),
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.current = snap
}

// Get returns the current snapshot and whether one exists.
func (s *SnapshotStore) Get() (*Snapshot, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.current == nil {
		return nil, false
	}
	copy := &Snapshot{
		Ports:      append([]scanner.PortState(nil), s.current.Ports...),
		CapturedAt: s.current.CapturedAt,
	}
	return copy, true
}

// Clear removes the current snapshot.
func (s *SnapshotStore) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.current = nil
}
