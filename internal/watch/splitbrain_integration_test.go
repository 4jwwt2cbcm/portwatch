package watch_test

import (
	"sync"
	"testing"
	"time"

	"portwatch/internal/watch"
)

func TestSplitBrainIntegration(t *testing.T) {
	sb := watch.NewSplitBrain(watch.SplitBrainPolicy{
		QuorumSize: 3,
		Window:     200 * time.Millisecond,
	})

	var wg sync.WaitGroup

	// Simulate 3 goroutines voting for "primary".
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sb.Vote("primary")
		}()
	}

	// Simulate 2 goroutines voting for "secondary" (split-brain).
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sb.Vote("secondary")
		}()
	}

	wg.Wait()

	if !sb.HasQuorum() {
		t.Error("expected quorum after 3 votes for primary")
	}
	if !sb.Conflicted() {
		t.Error("expected conflict since secondary also has votes")
	}

	// After window expires, both quorum and conflict should clear.
	time.Sleep(250 * time.Millisecond)
	if sb.HasQuorum() {
		t.Error("expected quorum to clear after window expiry")
	}
	if sb.Conflicted() {
		t.Error("expected conflict to clear after window expiry")
	}
}
