package watch

import (
	"sync/atomic"
	"testing"
	"time"
)

func makeDebounce(wait time.Duration) *Debounce {
	return NewDebounce(DebouncePolicy{Wait: wait})
}

func TestDefaultDebouncePolicy(t *testing.T) {
	p := DefaultDebouncePolicy()
	if p.Wait != 500*time.Millisecond {
		t.Fatalf("expected 500ms, got %v", p.Wait)
	}
}

func TestDebounceFiresAfterWait(t *testing.T) {
	d := makeDebounce(50 * time.Millisecond)
	var called int32
	d.Trigger(func() { atomic.StoreInt32(&called, 1) })
	time.Sleep(100 * time.Millisecond)
	if atomic.LoadInt32(&called) != 1 {
		t.Fatal("expected fn to be called after wait")
	}
}

func TestDebounceCollapsesCalls(t *testing.T) {
	d := makeDebounce(80 * time.Millisecond)
	var count int32
	for i := 0; i < 5; i++ {
		d.Trigger(func() { atomic.AddInt32(&count, 1) })
		time.Sleep(20 * time.Millisecond)
	}
	time.Sleep(150 * time.Millisecond)
	if n := atomic.LoadInt32(&count); n != 1 {
		t.Fatalf("expected 1 call, got %d", n)
	}
}

func TestDebouncePendingBeforeFire(t *testing.T) {
	d := makeDebounce(100 * time.Millisecond)
	d.Trigger(func() {})
	if !d.Pending() {
		t.Fatal("expected pending to be true before wait elapses")
	}
	time.Sleep(150 * time.Millisecond)
	if d.Pending() {
		t.Fatal("expected pending to be false after fire")
	}
}

func TestDebounceCancelPreventsCall(t *testing.T) {
	d := makeDebounce(80 * time.Millisecond)
	var called int32
	d.Trigger(func() { atomic.StoreInt32(&called, 1) })
	d.Cancel()
	time.Sleep(120 * time.Millisecond)
	if atomic.LoadInt32(&called) != 0 {
		t.Fatal("expected fn not to be called after cancel")
	}
	if d.Pending() {
		t.Fatal("expected pending false after cancel")
	}
}

func TestDebounceDefaultsWaitOnZero(t *testing.T) {
	d := NewDebounce(DebouncePolicy{Wait: 0})
	if d.policy.Wait != DefaultDebouncePolicy().Wait {
		t.Fatalf("expected default wait, got %v", d.policy.Wait)
	}
}
