package watch

import (
	"testing"
	"time"
)

func makeDedupWindow(window time.Duration) *DedupWindow {
	d := NewDedupWindow(window)
	return d
}

func TestDedupFirstCallNotDuplicate(t *testing.T) {
	d := makeDedupWindow(5 * time.Second)
	if d.IsDuplicate("port:tcp:8080") {
		t.Fatal("expected first call to not be a duplicate")
	}
}

func TestDedupSecondCallWithinWindowIsDuplicate(t *testing.T) {
	d := makeDedupWindow(5 * time.Second)
	d.IsDuplicate("port:tcp:8080")
	if !d.IsDuplicate("port:tcp:8080") {
		t.Fatal("expected second call within window to be a duplicate")
	}
}

func TestDedupAfterWindowExpires(t *testing.T) {
	now := time.Now()
	d := makeDedupWindow(5 * time.Second)
	d.nowFunc = func() time.Time { return now }
	d.IsDuplicate("port:tcp:9090")

	// Advance time past the window
	d.nowFunc = func() time.Time { return now.Add(6 * time.Second) }
	if d.IsDuplicate("port:tcp:9090") {
		t.Fatal("expected call after window expiry to not be a duplicate")
	}
}

func TestDedupDifferentKeysIndependent(t *testing.T) {
	d := makeDedupWindow(5 * time.Second)
	d.IsDuplicate("port:tcp:80")
	if d.IsDuplicate("port:tcp:443") {
		t.Fatal("expected different key to not be a duplicate")
	}
}

func TestDedupEvictRemovesExpired(t *testing.T) {
	now := time.Now()
	d := makeDedupWindow(2 * time.Second)
	d.nowFunc = func() time.Time { return now }
	d.IsDuplicate("key1")

	d.nowFunc = func() time.Time { return now.Add(3 * time.Second) }
	d.Evict()

	if len(d.seen) != 0 {
		t.Fatalf("expected seen map to be empty after eviction, got %d entries", len(d.seen))
	}
}

func TestDedupResetClearsAll(t *testing.T) {
	d := makeDedupWindow(10 * time.Second)
	d.IsDuplicate("a")
	d.IsDuplicate("b")
	d.Reset()
	if d.IsDuplicate("a") {
		t.Fatal("expected key to not be duplicate after reset")
	}
}
