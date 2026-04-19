package watch

import "sync"

// RouteFunc is a function that selects a route key from an envelope.
type RouteFunc[T any] func(env Envelope[T]) string

// EnvelopeRouter routes envelopes to named channels based on a route function.
type EnvelopeRouter[T any] struct {
	mu      sync.RWMutex
	routes  map[string]chan Envelope[T]
	routeFn RouteFunc[T]
	bufSize int
}

// NewEnvelopeRouter creates a router using the provided route function.
// bufSize controls the channel buffer per route (defaults to 8).
func NewEnvelopeRouter[T any](fn RouteFunc[T], bufSize int) *EnvelopeRouter[T] {
	if bufSize <= 0 {
		bufSize = 8
	}
	if fn == nil {
		fn = func(env Envelope[T]) string { return "default" }
	}
	return &EnvelopeRouter[T]{
		routes:  make(map[string]chan Envelope[T]),
		routeFn: fn,
		bufSize: bufSize,
	}
}

// Send routes the envelope to the appropriate channel, creating it if needed.
// Returns false if the destination channel is full.
func (r *EnvelopeRouter[T]) Send(env Envelope[T]) bool {
	key := r.routeFn(env)
	r.mu.Lock()
	ch, ok := r.routes[key]
	if !ok {
		ch = make(chan Envelope[T], r.bufSize)
		r.routes[key] = ch
	}
	r.mu.Unlock()

	select {
	case ch <- env:
		return true
	default:
		return false
	}
}

// Channel returns the receive channel for a given route key.
func (r *EnvelopeRouter[T]) Channel(key string) <-chan Envelope[T] {
	r.mu.Lock()
	defer r.mu.Unlock()
	ch, ok := r.routes[key]
	if !ok {
		ch = make(chan Envelope[T], r.bufSize)
		r.routes[key] = ch
	}
	return ch
}

// Routes returns the current set of known route keys.
func (r *EnvelopeRouter[T]) Routes() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	keys := make([]string, 0, len(r.routes))
	for k := range r.routes {
		keys = append(keys, k)
	}
	return keys
}
