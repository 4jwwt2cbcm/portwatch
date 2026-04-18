package watch

import (
	"sync"
	"testing"
)

func TestLatchNotSetOnInit(t *testing.T) {
	l := NewLatch()
	if l.IsSet() {
		t.Fatal("expected latch to be unset on init")
	}
}

func TestLatchSetReturnsTrue(t *testing.T) {
	l := NewLatch()
	if !l.Set() {
		t.Fatal("expected first Set to return true")
	}
	if !l.IsSet() {
		t.Fatal("expected latch to be set")
	}
}

func TestLatchSetTwiceReturnsFalse(t *testing.T) {
	l := NewLatch()
	l.Set()
	if l.Set() {
		t.Fatal("expected second Set to return false")
	}
}

func TestLatchResetAllowsReuse(t *testing.T) {
	l := NewLatch()
	l.Set()
	l.Reset()
	if l.IsSet() {
		t.Fatal("expected latch to be unset after reset")
	}
	if !l.Set() {
		t.Fatal("expected Set to return true after reset")
	}
}

func TestLatchSetOnceCallsFn(t *testing.T) {
	l := NewLatch()
	called := 0
	l.SetOnce(func() { called++ })
	l.SetOnce(func() { called++ })
	if called != 1 {
		t.Fatalf("expected fn called once, got %d", called)
	}
}

func TestLatchConcurrentSet(t *testing.T) {
	l := NewLatch()
	var wg sync.WaitGroup
	winners := 0
	var mu sync.Mutex
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if l.Set() {
				mu.Lock()
				winners++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	if winners != 1 {
		t.Fatalf("expected exactly 1 winner, got %d", winners)
	}
}
