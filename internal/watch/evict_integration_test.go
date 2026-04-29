package watch_test

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/portwatch/internal/watch"
)

func TestEvictCacheIntegration(t *testing.T) {
	var calls atomic.Int32
	fn := func(_ context.Context, key string) (string, error) {
		calls.Add(1)
		return "result-" + key, nil
	}

	cache := watch.NewEvictCache[string](watch.EvictPolicy{
		TTL:      50 * time.Millisecond,
		Capacity: 4,
	})
	runner := watch.NewEvictRunner(cache, fn)

	ctx := context.Background()
	const goroutines = 8
	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			v, err := runner.Run(ctx, "shared")
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if v != "result-shared" {
				t.Errorf("unexpected value: %q", v)
			}
		}()
	}
	wg.Wait()

	// After TTL, fn should be called again.
	time.Sleep(60 * time.Millisecond)
	runner.Run(ctx, "shared")
	if calls.Load() < 2 {
		t.Fatalf("expected at least 2 calls after TTL expiry, got %d", calls.Load())
	}
}
