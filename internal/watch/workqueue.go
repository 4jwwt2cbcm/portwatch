package watch

import "sync"

// WorkQueue is a bounded, thread-safe FIFO queue for scan work items.
type WorkQueue struct {
	mu       sync.Mutex
	items    []string
	cap      int
	notify   chan struct{}
}

// NewWorkQueue creates a WorkQueue with the given capacity (min 1).
func NewWorkQueue(cap int) *WorkQueue {
	if cap < 1 {
		cap = 16
	}
	return &WorkQueue{
		items:  make([]string, 0, cap),
		cap:    cap,
		notify: make(chan struct{}, 1),
	}
}

// Push adds an item to the queue. Returns false if the queue is full.
func (q *WorkQueue) Push(item string) bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.items) >= q.cap {
		return false
	}
	q.items = append(q.items, item)
	select {
	case q.notify <- struct{}{}:
	default:
	}
	return true
}

// Pop removes and returns the next item. Returns "", false if empty.
func (q *WorkQueue) Pop() (string, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.items) == 0 {
		return "", false
	}
	item := q.items[0]
	q.items = q.items[1:]
	return item, true
}

// Len returns the current number of items.
func (q *WorkQueue) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.items)
}

// Notify returns a channel that receives a signal when items are pushed.
func (q *WorkQueue) Notify() <-chan struct{} {
	return q.notify
}

// Drain returns all current items and clears the queue.
func (q *WorkQueue) Drain() []string {
	q.mu.Lock()
	defer q.mu.Unlock()
	out := make([]string, len(q.items))
	copy(out, q.items)
	q.items = q.items[:0]
	return out
}
