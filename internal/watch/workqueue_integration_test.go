package watch

import (
	"sync"
	"testing"
)

func TestWorkQueueConcurrentPushPop(t *testing.T) {
	q := NewWorkQueue(64)
	var wg sync.WaitGroup
	const producers = 4
	const items = 10

	for i := 0; i < producers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < items; j++ {
				q.Push("work")
			}
		}()
	}
	wg.Wait()

	if q.Len() != producers*items {
		t.Fatalf("expected %d items, got %d", producers*items, q.Len())
	}

	var mu sync.Mutex
	popped := 0
	for i := 0; i < producers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				_, ok := q.Pop()
				if !ok {
					return
				}
				mu.Lock()
				popped++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()

	if popped != producers*items {
		t.Fatalf("expected %d popped, got %d", producers*items, popped)
	}
}
