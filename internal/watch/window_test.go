package watch

import (
	"testing"
	"time"
)

func makeWindow[T any](size time.Duration, max int) *RollingWindow[T] {
	return NewRollingWindow[T](WindowPolicy{Size: size, MaxItems: max})
}

func TestDefaultWindowPolicyValues(t *testing.T) {
	p := DefaultWindowPolicy()
	if p.Size != 10*time.Second {
		t.Errorf("expected 10s, got %v", p.Size)
	}
	if p.MaxItems != 100 {
		t.Errorf("expected 100, got %d", p.MaxItems)
	}
}

func TestRollingWindowEmptyOnInit(t *testing.T) {
	w := makeWindow[int](time.Second, 10)
	if w.Len() != 0 {
		t.Fatalf("expected 0, got %d", w.Len())
	}
}

func TestRollingWindowAddAndSnapshot(t *testing.T) {
	w := makeWindow[string](time.Second, 10)
	w.Add("a")
	w.Add("b")
	snap := w.Snapshot()
	if len(snap) != 2 {
		t.Fatalf("expected 2 items, got %d", len(snap))
	}
	if snap[0] != "a" || snap[1] != "b" {
		t.Errorf("unexpected snapshot: %v", snap)
	}
}

func TestRollingWindowEvictsExpiredEntries(t *testing.T) {
	w := makeWindow[int](100*time.Millisecond, 10)
	base := time.Now()
	w.now = func() time.Time { return base }
	w.Add(1)
	w.Add(2)
	// advance time past window
	w.now = func() time.Time { return base.Add(200 * time.Millisecond) }
	w.Add(3)
	snap := w.Snapshot()
	if len(snap) != 1 || snap[0] != 3 {
		t.Errorf("expected [3], got %v", snap)
	}
}

func TestRollingWindowCapsAtMaxItems(t *testing.T) {
	w := makeWindow[int](time.Minute, 3)
	for i := 0; i < 5; i++ {
		w.Add(i)
	}
	if w.Len() != 3 {
		t.Fatalf("expected 3, got %d", w.Len())
	}
	snap := w.Snapshot()
	// oldest items should have been evicted
	if snap[0] != 2 || snap[2] != 4 {
		t.Errorf("unexpected snapshot after cap: %v", snap)
	}
}

func TestRollingWindowClear(t *testing.T) {
	w := makeWindow[int](time.Second, 10)
	w.Add(1)
	w.Add(2)
	w.Clear()
	if w.Len() != 0 {
		t.Errorf("expected 0 after clear, got %d", w.Len())
	}
}

func TestRollingWindowDefaultsOnZeroPolicy(t *testing.T) {
	w := NewRollingWindow[int](WindowPolicy{})
	if w.policy.Size != DefaultWindowPolicy().Size {
		t.Errorf("expected default size, got %v", w.policy.Size)
	}
	if w.policy.MaxItems != DefaultWindowPolicy().MaxItems {
		t.Errorf("expected default max items, got %d", w.policy.MaxItems)
	}
}
