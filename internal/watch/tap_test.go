package watch

import (
	"testing"
)

func makeTap() *Tap[int] {
	return NewTap[int](8, nil)
}

func TestTapDefaultCapOnZero(t *testing.T) {
	tp := NewTap[int](0, nil)
	for i := 0; i < 64; i++ {
		tp.Record(i)
	}
	if tp.Len() != 64 {
		t.Fatalf("expected 64, got %d", tp.Len())
	}
}

func TestTapRecordStoresValue(t *testing.T) {
	tp := makeTap()
	tp.Record(42)
	if tp.Len() != 1 {
		t.Fatalf("expected 1, got %d", tp.Len())
	}
}

func TestTapRecordReturnsValue(t *testing.T) {
	tp := makeTap()
	out := tp.Record(7)
	if out != 7 {
		t.Fatalf("expected 7, got %d", out)
	}
}

func TestTapCapLimitsStorage(t *testing.T) {
	tp := NewTap[int](3, nil)
	for i := 0; i < 10; i++ {
		tp.Record(i)
	}
	if tp.Len() != 3 {
		t.Fatalf("expected 3, got %d", tp.Len())
	}
}

func TestTapSnapshotReturnsCopy(t *testing.T) {
	tp := makeTap()
	tp.Record(1)
	tp.Record(2)
	snap := tp.Snapshot()
	snap[0] = 99
	if tp.Snapshot()[0] != 1 {
		t.Fatal("snapshot mutation affected tap")
	}
}

func TestTapClearResetsLen(t *testing.T) {
	tp := makeTap()
	tp.Record(1)
	tp.Clear()
	if tp.Len() != 0 {
		t.Fatalf("expected 0, got %d", tp.Len())
	}
}

func TestTapOnTapCallbackFired(t *testing.T) {
	var fired []int
	tp := NewTap[int](8, func(v int) { fired = append(fired, v) })
	tp.Record(5)
	tp.Record(6)
	if len(fired) != 2 || fired[0] != 5 || fired[1] != 6 {
		t.Fatalf("unexpected fired values: %v", fired)
	}
}
