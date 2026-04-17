package watch

import (
	"testing"
)

func TestWorkQueueDefaultsCapOnZero(t *testing.T) {
	q := NewWorkQueue(0)
	if q.cap != 16 {
		t.Fatalf("expected cap 16, got %d", q.cap)
	}
}

func TestWorkQueuePushAndPop(t *testing.T) {
	q := NewWorkQueue(4)
	if !q.Push("a") {
		t.Fatal("expected push to succeed")
	}
	item, ok := q.Pop()
	if !ok || item != "a" {
		t.Fatalf("expected 'a', got %q %v", item, ok)
	}
}

func TestWorkQueuePopEmptyReturnsFalse(t *testing.T) {
	q := NewWorkQueue(4)
	_, ok := q.Pop()
	if ok {
		t.Fatal("expected false on empty pop")
	}
}

func TestWorkQueuePushReturnsFalseWhenFull(t *testing.T) {
	q := NewWorkQueue(2)
	q.Push("a")
	q.Push("b")
	if q.Push("c") {
		t.Fatal("expected push to fail when full")
	}
}

func TestWorkQueueLen(t *testing.T) {
	q := NewWorkQueue(4)
	q.Push("x")
	q.Push("y")
	if q.Len() != 2 {
		t.Fatalf("expected len 2, got %d", q.Len())
	}
}

func TestWorkQueueDrain(t *testing.T) {
	q := NewWorkQueue(4)
	q.Push("a")
	q.Push("b")
	out := q.Drain()
	if len(out) != 2 {
		t.Fatalf("expected 2 items, got %d", len(out))
	}
	if q.Len() != 0 {
		t.Fatal("expected queue empty after drain")
	}
}

func TestWorkQueueNotifySignalsOnPush(t *testing.T) {
	q := NewWorkQueue(4)
	q.Push("item")
	select {
	case <-q.Notify():
	default:
		t.Fatal("expected notify signal after push")
	}
}
