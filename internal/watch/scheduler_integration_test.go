package watch

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

func TestSchedulerIntegration(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	var count atomic.Int64
	policy := SchedulerPolicy{
		InitialDelay: 0,
		Interval:     30 * time.Millisecond,
		Jitter:       0,
	}
	s := NewScheduler(policy, func(ctx context.Context) error {
		count.Add(1)
		return nil
	})

	_ = s.Run(ctx)

	got := count.Load()
	if got < 3 {
		t.Errorf("expected at least 3 invocations in 200ms, got %d", got)
	}
	if s.Fires() != int(got) {
		t.Errorf("Fires() %d does not match callback count %d", s.Fires(), got)
	}
}
