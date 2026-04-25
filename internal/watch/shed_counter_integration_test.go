package watch_test

import (
	"testing"
	"time"

	"portwatch/internal/watch"
)

func TestShedCounterIntegration(t *testing.T) {
	policy := watch.ShedCounterPolicy{
		Window:   80 * time.Millisecond,
		Capacity: 50,
	}
	s := watch.NewShedCounter(policy)

	// Record a burst of shed events.
	for i := 0; i < 10; i++ {
		s.Record()
	}
	if got := s.Count(); got != 10 {
		t.Fatalf("expected 10 events recorded, got %d", got)
	}

	// Wait for the window to expire; count should drop to zero.
	time.Sleep(100 * time.Millisecond)
	if got := s.Count(); got != 0 {
		t.Fatalf("expected 0 after window expiry, got %d", got)
	}

	// Record again and verify reset works.
	s.Record()
	s.Record()
	s.Reset()
	if got := s.Count(); got != 0 {
		t.Fatalf("expected 0 after reset, got %d", got)
	}
}
