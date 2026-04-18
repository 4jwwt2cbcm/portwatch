package watch

import (
	"sync"
	"testing"
)

func TestObserverNoSubscribersIsNoop(t *testing.T) {
	o := NewObserver()
	// should not panic
	o.Publish("scan", nil)
}

func TestObserverSubscribeAndPublish(t *testing.T) {
	o := NewObserver()
	var got []any
	o.Subscribe("scan", func(_ string, p any) { got = append(got, p) })
	o.Publish("scan", 42)
	if len(got) != 1 || got[0] != 42 {
		t.Fatalf("expected [42], got %v", got)
	}
}

func TestObserverMultipleSubscribers(t *testing.T) {
	o := NewObserver()
	count := 0
	o.Subscribe("e", func(_ string, _ any) { count++ })
	o.Subscribe("e", func(_ string, _ any) { count++ })
	o.Publish("e", nil)
	if count != 2 {
		t.Fatalf("expected 2, got %d", count)
	}
}

func TestObserverUnsubscribe(t *testing.T) {
	o := NewObserver()
	count := 0
	unsub := o.Subscribe("e", func(_ string, _ any) { count++ })
	unsub()
	o.Publish("e", nil)
	if count != 0 {
		t.Fatalf("expected 0 after unsubscribe, got %d", count)
	}
}

func TestObserverCount(t *testing.T) {
	o := NewObserver()
	if o.Count("e") != 0 {
		t.Fatal("expected 0")
	}
	o.Subscribe("e", func(_ string, _ any) {})
	o.Subscribe("e", func(_ string, _ any) {})
	if o.Count("e") != 2 {
		t.Fatalf("expected 2, got %d", o.Count("e"))
	}
}

func TestObserverClear(t *testing.T) {
	o := NewObserver()
	o.Subscribe("e", func(_ string, _ any) {})
	o.Clear("e")
	if o.Count("e") != 0 {
		t.Fatal("expected 0 after clear")
	}
}

func TestObserverConcurrentPublish(t *testing.T) {
	o := NewObserver()
	var mu sync.Mutex
	results := []any{}
	o.Subscribe("e", func(_ string, p any) {
		mu.Lock()
		results = append(results, p)
		mu.Unlock()
	})
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(v int) { defer wg.Done(); o.Publish("e", v) }(i)
	}
	wg.Wait()
	if len(results) != 20 {
		t.Fatalf("expected 20, got %d", len(results))
	}
}
