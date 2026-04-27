package watch

import (
	"testing"
	"time"
)

func TestEpochZeroOnInit(t *testing.T) {
	e := NewEpoch()
	if e.Generation() != 0 {
		t.Fatalf("expected generation 0, got %d", e.Generation())
	}
}

func TestEpochAdvanceIncrementsGeneration(t *testing.T) {
	e := NewEpoch()
	gen := e.Advance()
	if gen != 1 {
		t.Fatalf("expected Advance to return 1, got %d", gen)
	}
	if e.Generation() != 1 {
		t.Fatalf("expected generation 1, got %d", e.Generation())
	}
}

func TestEpochAdvanceMultipleTimes(t *testing.T) {
	e := NewEpoch()
	for i := uint64(1); i <= 5; i++ {
		gen := e.Advance()
		if gen != i {
			t.Fatalf("expected generation %d, got %d", i, gen)
		}
	}
}

func TestEpochResetClearsGeneration(t *testing.T) {
	e := NewEpoch()
	e.Advance()
	e.Advance()
	e.Reset()
	if e.Generation() != 0 {
		t.Fatalf("expected generation 0 after reset, got %d", e.Generation())
	}
}

func TestEpochSnapshotIsConsistentCopy(t *testing.T) {
	e := NewEpoch()
	e.Advance()
	e.Advance()
	snap := e.Snapshot()
	if snap.Generation != 2 {
		t.Fatalf("expected snapshot generation 2, got %d", snap.Generation)
	}
	if snap.StartedAt.IsZero() {
		t.Fatal("expected StartedAt to be set")
	}
	if snap.ResetAt.IsZero() {
		t.Fatal("expected ResetAt to be set")
	}
}

func TestEpochSnapshotResetAtUpdatesOnAdvance(t *testing.T) {
	now := time.Now()
	e := NewEpoch()
	e.now = func() time.Time { return now.Add(5 * time.Second) }
	e.Advance()
	snap := e.Snapshot()
	if !snap.ResetAt.Equal(now.Add(5 * time.Second)) {
		t.Fatalf("expected ResetAt to reflect mocked time")
	}
}

func TestEpochSinceReturnsTrueForOlderGeneration(t *testing.T) {
	e := NewEpoch()
	e.Advance()
	e.Advance()
	if !e.Since(1) {
		t.Fatal("expected Since(1) to return true when current generation is 2")
	}
}

func TestEpochSinceReturnsFalseForCurrentGeneration(t *testing.T) {
	e := NewEpoch()
	e.Advance()
	if e.Since(1) {
		t.Fatal("expected Since(1) to return false when current generation is 1")
	}
}

func TestEpochSinceReturnsFalseForFutureGeneration(t *testing.T) {
	e := NewEpoch()
	if e.Since(5) {
		t.Fatal("expected Since(5) to return false when current generation is 0")
	}
}
