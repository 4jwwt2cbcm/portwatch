package watch_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/portwatch/internal/watch"
)

// TestWatchdogIntegration verifies that a Watchdog fires when the monitored
// goroutine stops kicking, and stays silent while kicks arrive on time.
func TestWatchdogIntegration(t *testing.T) {
	var fireCount atomic.Int32

	policy := watch.WatchdogPolicy{
		Timeout:  30 * time.Millisecond,
		Interval: 5 * time.Millisecond,
	}
	wd := watch.NewWatchdog(policy, func() {
		fireCount.Add(1)
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() { wd.Run(ctx) }() //nolint:errcheck

	// Phase 1: kick regularly — watchdog must stay quiet.
	for i := 0; i < 8; i++ {
		time.Sleep(10 * time.Millisecond)
		wd.Kick()
	}
	if fireCount.Load() != 0 {
		t.Fatalf("phase 1: unexpected fires: %d", fireCount.Load())
	}

	// Phase 2: stop kicking — watchdog must fire.
	time.Sleep(80 * time.Millisecond)
	if fireCount.Load() == 0 {
		t.Fatal("phase 2: watchdog did not fire after timeout")
	}
}
