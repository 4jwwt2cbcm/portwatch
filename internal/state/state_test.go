package state_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/state"
)

func tempPath(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "state.json")
}

func TestSaveAndLoad(t *testing.T) {
	ports := []scanner.PortState{
		{Port: 80, Proto: "tcp", State: "open"},
		{Port: 443, Proto: "tcp", State: "open"},
	}

	store := state.NewStore(tempPath(t))

	if err := store.Save(ports); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	snap, err := store.Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if len(snap.Ports) != len(ports) {
		t.Fatalf("expected %d ports, got %d", len(ports), len(snap.Ports))
	}
	if snap.Ports[0].Port != 80 {
		t.Errorf("expected port 80, got %d", snap.Ports[0].Port)
	}
	if snap.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestLoadMissingFileReturnsEmpty(t *testing.T) {
	store := state.NewStore("/tmp/portwatch_nonexistent_abc123.json")

	snap, err := store.Load()
	if err != nil {
		t.Fatalf("Load() unexpected error: %v", err)
	}
	if len(snap.Ports) != 0 {
		t.Errorf("expected empty ports, got %d", len(snap.Ports))
	}
}

func TestSaveOverwritesPreviousState(t *testing.T) {
	path := tempPath(t)
	store := state.NewStore(path)

	first := []scanner.PortState{{Port: 22, Proto: "tcp", State: "open"}}
	second := []scanner.PortState{{Port: 8080, Proto: "tcp", State: "open"}}

	_ = store.Save(first)
	_ = store.Save(second)

	snap, err := store.Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if len(snap.Ports) != 1 || snap.Ports[0].Port != 8080 {
		t.Errorf("expected port 8080 after overwrite, got %+v", snap.Ports)
	}
}

func TestLoadCorruptFileReturnsError(t *testing.T) {
	path := tempPath(t)
	_ = os.WriteFile(path, []byte("not valid json{"), 0o644)

	store := state.NewStore(path)
	_, err := store.Load()
	if err == nil {
		t.Error("expected error for corrupt JSON, got nil")
	}
}
