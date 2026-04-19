package watch

import (
	"container/heap"
	"sync"
)

// PriorityItem wraps a value with a priority level.
type PriorityItem[T any] struct {
	Value    T
	Priority int
	index    int
}

type priorityHeap[T any] []*PriorityItem[T]

func (h priorityHeap[T]) Len() int            { return len(h) }
func (h priorityHeap[T]) Less(i, j int) bool  { return h[i].Priority > h[j].Priority }
func (h priorityHeap[T]) Swap(i, j int)       { h[i], h[j] = h[j], h[i]; h[i].index = i; h[j].index = j }
func (h *priorityHeap[T]) Push(x any)         { item := x.(*PriorityItem[T]); item.index = len(*h); *h = append(*h, item) }
func (h *priorityHeap[T]) Pop() any           { old := *h; n := len(old); item := old[n-1]; old[n-1] = nil; *h = old[:n-1]; return item }

// PriorityQueue is a thread-safe generic max-priority queue.
type PriorityQueue[T any] struct {
	mu   sync.Mutex
	h    priorityHeap[T]
	cap  int
}

// NewPriorityQueue returns a PriorityQueue with optional max capacity (0 = unlimited).
func NewPriorityQueue[T any](cap int) *PriorityQueue[T] {
	h := priorityHeap[T]{}
	heap.Init(&h)
	return &PriorityQueue[T]{h: h, cap: cap}
}

// Push adds an item. Returns false if the queue is at capacity.
func (q *PriorityQueue[T]) Push(value T, priority int) bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.cap > 0 && len(q.h) >= q.cap {
		return false
	}
	heap.Push(&q.h, &PriorityItem[T]{Value: value, Priority: priority})
	return true
}

// Pop removes and returns the highest-priority item. Returns false if empty.
func (q *PriorityQueue[T]) Pop() (T, int, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.h) == 0 {
		var zero T
		return zero, 0, false
	}
	item := heap.Pop(&q.h).(*PriorityItem[T])
	return item.Value, item.Priority, true
}

// Len returns the current number of items.
func (q *PriorityQueue[T]) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.h)
}
