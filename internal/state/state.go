package state

import (
	"encoding/json"
	"os"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Snapshot represents a persisted port scan result.
type Snapshot struct {
	Timestamp time.Time          `json:"timestamp"`
	Ports     []scanner.PortState `json:"ports"`
}

// Store handles reading and writing port state snapshots to disk.
type Store struct {
	path string
}

// NewStore creates a new Store backed by the given file path.
func NewStore(path string) *Store {
	return &Store{path: path}
}

// Save writes the current port states to disk as a JSON snapshot.
func (s *Store) Save(ports []scanner.PortState) error {
	snap := Snapshot{
		Timestamp: time.Now().UTC(),
		Ports:     ports,
	}
	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o644)
}

// Load reads the last saved snapshot from disk.
// Returns an empty snapshot and no error if the file does not exist yet.
func (s *Store) Load() (Snapshot, error) {
	data, err := os.ReadFile(s.path)
	if os.IsNotExist(err) {
		return Snapshot{}, nil
	}
	if err != nil {
		return Snapshot{}, err
	}
	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return Snapshot{}, err
	}
	return snap, nil
}
