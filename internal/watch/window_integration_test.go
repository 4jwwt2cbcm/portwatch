package watch_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/watch"
)

// TestRollingWindowIntegration verifies that the RollingWindow correctly
// evicts stale entries and caps storage under concurrent-style sequential load.
func TestRollingWindowIntegration(t *testing.T) {
	policy := watch.WindowPolicy{
		Size:     50 * time.Millisecond,
		MaxItems: 20,
	}
	w := watch.NewRollingWindow[int](policy)

	base := time.Now()
	// Inject a controlled clock via the exported field would require a helper;
	// instead we rely on real time with a short window.

	// Add a burst of items.
	for i := 0; i < 10; i++ {
		w.Add(i)
	}
	if w.Len() != 10 {
		t.Fatalf("expected 10 items after burst, got %d", w.Len())
	}

	// Wait for window to expire.
	time.Sleep(60 * time.Millisecond)

	// Add fresh items.
	for i := 100; i < 103; i++ {
		w.Add(i)
	}

	snap := w.Snapshot()
	if len(snap) != 3 {
		t.Fatalf("expected 3 fresh items after expiry, got %d (elapsed: %v)", len(snap), time.Since(base))
	}
	for _, v := range snap {
		if v < 100 {
			t.Errorf("stale value %d found in window after expiry", v)
		}
	}
}
