package watch

import "sync"

// FanOut distributes a single input value to multiple consumer channels.
type FanOut[T any] struct {
	mu       sync.RWMutex
	subs     []chan T
	bufSize  int
}

// NewFanOut creates a FanOut with the given per-subscriber channel buffer size.
func NewFanOut[T any](bufSize int) *FanOut[T] {
	if bufSize < 1 {
		bufSize = 1
	}
	return &FanOut[T]{bufSize: bufSize}
}

// Subscribe registers a new subscriber and returns its receive channel.
func (f *FanOut[T]) Subscribe() <-chan T {
	ch := make(chan T, f.bufSize)
	f.mu.Lock()
	f.subs = append(f.subs, ch)
	f.mu.Unlock()
	return ch
}

// Publish sends value to all current subscribers.
// Subscribers whose channels are full are skipped (non-blocking).
func (f *FanOut[T]) Publish(v T) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	for _, ch := range f.subs {
		select {
		case ch <- v:
		default:
		}
	}
}

// Close closes all subscriber channels.
func (f *FanOut[T]) Close() {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, ch := range f.subs {
		close(ch)
	}
	f.subs = nil
}
