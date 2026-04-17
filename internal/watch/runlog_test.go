package watch

import (
	"errors"
	"testing"
	"time"
)

func makeEntry(portsFound int, err error) RunEntry {
	now := time.Now()
	return RunEntry{
		StartedAt:  now,
		FinishedAt: now.Add(10 * time.Millisecond),
		Err:        err,
		PortsFound: portsFound,
	}
}

func TestRunLogEmptyOnInit(t *testing.T) {
	rl := NewRunLog(0)
	if rl.Len() != 0 {
		t.Fatalf("expected 0, got %d", rl.Len())
	}
}

func TestRunLogDefaultCapOnZero(t *testing.T) {
	rl := NewRunLog(0)
	if rl.cap != 64 {
		t.Fatalf("expected cap 64, got %d", rl.cap)
	}
}

func TestRunLogAppendAndSnapshot(t *testing.T) {
	rl := NewRunLog(10)
	rl.Append(makeEntry(3, nil))
	rl.Append(makeEntry(5, errors.New("oops")))

	snap := rl.Snapshot()
	if len(snap) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(snap))
	}
	if snap[0].PortsFound != 3 {
		t.Errorf("expected 3 ports, got %d", snap[0].PortsFound)
	}
	if snap[1].Err == nil {
		t.Error("expected error in second entry")
	}
}

func TestRunLogEvictsOldestWhenFull(t *testing.T) {
	rl := NewRunLog(3)
	for i := 0; i < 4; i++ {
		rl.Append(makeEntry(i, nil))
	}
	snap := rl.Snapshot()
	if len(snap) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(snap))
	}
	if snap[0].PortsFound != 1 {
		t.Errorf("expected oldest evicted, first entry ports=1, got %d", snap[0].PortsFound)
	}
}

func TestRunLogClear(t *testing.T) {
	rl := NewRunLog(10)
	rl.Append(makeEntry(1, nil))
	rl.Clear()
	if rl.Len() != 0 {
		t.Fatalf("expected 0 after clear, got %d", rl.Len())
	}
}

func TestRunEntryDuration(t *testing.T) {
	now := time.Now()
	e := RunEntry{
		StartedAt:  now,
		FinishedAt: now.Add(50 * time.Millisecond),
	}
	if e.Duration() != 50*time.Millisecond {
		t.Errorf("unexpected duration: %v", e.Duration())
	}
}
