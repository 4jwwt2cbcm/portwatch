package watch_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/user/portwatch/internal/watch"
)

func TestQuotaIntegration(t *testing.T) {
	const goroutines = 20
	const max = 10

	q := watch.NewQuota(watch.QuotaPolicy{Max: max, Window: time.Second})
	r := watch.NewQuotaRunner(q, func(_ context.Context) error { return nil })

	var (
		wg      sync.WaitGroup
		mu      sync.Mutex
		allowed int
		denied  int
	)

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := r.Run(context.Background())
			mu.Lock()
			defer mu.Unlock()
			if errors.Is(err, watch.ErrQuotaExceeded) {
				denied++
			} else if err == nil {
				allowed++
			}
		}()
	}
	wg.Wait()

	if allowed != max {
		t.Errorf("expected %d allowed, got %d", max, allowed)
	}
	if denied != goroutines-max {
		t.Errorf("expected %d denied, got %d", goroutines-max, denied)
	}
}
