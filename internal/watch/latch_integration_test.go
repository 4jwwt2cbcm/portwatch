package watch_test

import (
	"sync"
	"testing"

	"github.com/user/portwatch/internal/watch"
)

func TestLatchIntegration(t *testing.T) {
	// Simulate a one-shot startup signal across goroutines.
	l := watch.NewLatch()
	var wg sync.WaitGroup
	notified := make([]int, 0, 10)
	var mu sync.Mutex

	for i := 0; i < 10; i++ {
		wg.Add(1)
		id := i
		go func() {
			defer wg.Done()
			l.SetOnce(func() {
				mu.Lock()
				notified = append(notified, id)
				mu.Unlock()
			})
		}()
	}
	wg.Wait()

	if len(notified) != 1 {
		t.Fatalf("expected exactly one notification, got %d", len(notified))
	}
	if !l.IsSet() {
		t.Fatal("expected latch to remain set after goroutines finish")
	}
}
