package watch

import (
	"testing"
)

func TestPriorityQueueEmptyOnInit(t *testing.T) {
	q := NewPriorityQueue[string](0)
	if q.Len() != 0 {
		t.Fatalf("expected 0, got %d", q.Len())
	}
}

func TestPriorityQueuePopEmptyReturnsFalse(t *testing.T) {
	q := NewPriorityQueue[int](0)
	_, _, ok := q.Pop()
	if ok {
		t.Fatal("expected false on empty pop")
	}
}

func TestPriorityQueueOrdersByPriority(t *testing.T) {
	q := NewPriorityQueue[string](0)
	q.Push("low", 1)
	q.Push("high", 10)
	q.Push("mid", 5)

	v, p, ok := q.Pop()
	if !ok || v != "high" || p != 10 {
		t.Fatalf("expected high/10, got %s/%d", v, p)
	}
	v, p, ok = q.Pop()
	if !ok || v != "mid" || p != 5 {
		t.Fatalf("expected mid/5, got %s/%d", v, p)
	}
}

func TestPriorityQueueCapacityLimit(t *testing.T) {
	q := NewPriorityQueue[int](2)
	if !q.Push(1, 1) {
		t.Fatal("first push should succeed")
	}
	if !q.Push(2, 2) {
		t.Fatal("second push should succeed")
	}
	if q.Push(3, 3) {
		t.Fatal("third push should fail at capacity")
	}
	if q.Len() != 2 {
		t.Fatalf("expected len 2, got %d", q.Len())
	}
}

func TestPriorityQueueLenDecrementsOnPop(t *testing.T) {
	q := NewPriorityQueue[string](0)
	q.Push("a", 1)
	q.Push("b", 2)
	q.Pop()
	if q.Len() != 1 {
		t.Fatalf("expected len 1 after pop, got %d", q.Len())
	}
}

func TestPriorityQueueUnlimitedCapacity(t *testing.T) {
	q := NewPriorityQueue[int](0)
	for i := 0; i < 100; i++ {
		if !q.Push(i, i) {
			t.Fatalf("push %d should succeed with unlimited cap", i)
		}
	}
	if q.Len() != 100 {
		t.Fatalf("expected 100, got %d", q.Len())
	}
}
