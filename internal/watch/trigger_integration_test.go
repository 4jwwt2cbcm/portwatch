package watch_test

import (
	"context"
	"testing"
	"time"

	"github.com/user/portwatch/internal/watch"
)

func TestTriggerRunnerIntegration(t *testing.T) {
	policy := watch.TriggerPolicy{MinInterval: 10 * time.Millisecond}
	tr := watch.NewTrigger(policy)

	results := make([]int, 0, 3)
	counter := 0

	r := watch.NewTriggerRunner(tr, func(_ context.Context) error {
		counter++
		results = append(results, counter)
		return nil
	})

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	done := make(chan error, 1)
	go func() { done <- r.RunLoop(ctx) }()

	// fire three times with enough spacing to respect MinInterval
	for i := 0; i < 3; i++ {
		time.Sleep(20 * time.Millisecond)
		tr.Fire()
	}

	<-done

	if len(results) != 3 {
		t.Fatalf("expected 3 executions, got %d", len(results))
	}
	if tr.Count() != 3 {
		t.Fatalf("expected trigger count 3, got %d", tr.Count())
	}
}
