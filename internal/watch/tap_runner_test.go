package watch

import (
	"errors"
	"testing"
)

func TestTapRunnerRecordsNilOnSuccess(t *testing.T) {
	r := NewTapRunner(func() error { return nil }, 0)
	if err := r.Run(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Tap().Len() != 1 {
		t.Fatalf("expected 1 recorded result, got %d", r.Tap().Len())
	}
	if r.Tap().Snapshot()[0] != nil {
		t.Fatal("expected nil recorded")
	}
}

func TestTapRunnerRecordsError(t *testing.T) {
	sentinel := errors.New("boom")
	r := NewTapRunner(func() error { return sentinel }, 0)
	err := r.Run()
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel, got %v", err)
	}
	snap := r.Tap().Snapshot()
	if len(snap) != 1 || !errors.Is(snap[0], sentinel) {
		t.Fatalf("unexpected tap snapshot: %v", snap)
	}
}

func TestTapRunnerMultipleRuns(t *testing.T) {
	calls := 0
	r := NewTapRunner(func() error { calls++; return nil }, 4)
	for i := 0; i < 3; i++ {
		r.Run()
	}
	if r.Tap().Len() != 3 {
		t.Fatalf("expected 3, got %d", r.Tap().Len())
	}
}

func TestTapRunnerNilFnDefaults(t *testing.T) {
	r := NewTapRunner(nil, 0)
	if err := r.Run(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTapRunnerPropagatesError(t *testing.T) {
	sentinel := errors.New("fail")
	r := NewTapRunner(func() error { return sentinel }, 0)
	if err := r.Run(); !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel, got %v", err)
	}
}
