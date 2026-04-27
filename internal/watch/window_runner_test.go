package watch

import (
	"errors"
	"testing"
	"time"
)

func TestWindowRunnerRecordsSuccessfulResult(t *testing.T) {
	win := makeWindow[int](time.Second, 10)
	calls := 0
	r := NewWindowRunner(win, func() (int, error) {
		calls++
		return calls, nil
	})
	for i := 0; i < 3; i++ {
		r.Run() //nolint:errcheck
	}
	snap := win.Snapshot()
	if len(snap) != 3 {
		t.Fatalf("expected 3 items, got %d", len(snap))
	}
	if snap[2] != 3 {
		t.Errorf("expected last value 3, got %d", snap[2])
	}
}

func TestWindowRunnerDoesNotRecordOnError(t *testing.T) {
	win := makeWindow[int](time.Second, 10)
	errBoom := errors.New("boom")
	r := NewWindowRunner(win, func() (int, error) {
		return 0, errBoom
	})
	_, err := r.Run()
	if !errors.Is(err, errBoom) {
		t.Fatalf("expected errBoom, got %v", err)
	}
	if win.Len() != 0 {
		t.Errorf("expected empty window on error, got %d", win.Len())
	}
}

func TestWindowRunnerPropagatesError(t *testing.T) {
	win := makeWindow[string](time.Second, 10)
	errFail := errors.New("fail")
	r := NewWindowRunner(win, func() (string, error) {
		return "", errFail
	})
	_, err := r.Run()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestWindowRunnerNilArgDefaults(t *testing.T) {
	r := NewWindowRunner[int](nil, nil)
	if r.win == nil {
		t.Fatal("expected non-nil window")
	}
	v, err := r.Run()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != 0 {
		t.Errorf("expected zero value, got %d", v)
	}
}

func TestWindowRunnerWindowAccessor(t *testing.T) {
	win := makeWindow[int](time.Second, 5)
	r := NewWindowRunner(win, func() (int, error) { return 42, nil })
	if r.Window() != win {
		t.Error("Window() should return the same window instance")
	}
}
