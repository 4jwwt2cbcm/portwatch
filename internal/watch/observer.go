package watch

import "sync"

// ObserverFunc is called when a named event is published.
type ObserverFunc func(event string, payload any)

// Observer is a simple publish/subscribe registry keyed by event name.
type Observer struct {
	mu      sync.RWMutex
	subs    map[string][]ObserverFunc
}

// NewObserver creates an empty Observer.
func NewObserver() *Observer {
	return &Observer{subs: make(map[string][]ObserverFunc)}
}

// Subscribe registers fn to be called whenever event is published.
// Returns an unsubscribe function.
func (o *Observer) Subscribe(event string, fn ObserverFunc) func() {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.subs[event] = append(o.subs[event], fn)
	idx := len(o.subs[event]) - 1
	return func() {
		o.mu.Lock()
		defer o.mu.Unlock()
		slice := o.subs[event]
		if idx < len(slice) {
			slice[idx] = nil
		}
	}
}

// Publish calls all subscribers registered for event.
func (o *Observer) Publish(event string, payload any) {
	o.mu.RLock()
	fns := make([]ObserverFunc, len(o.subs[event]))
	copy(fns, o.subs[event])
	o.mu.RUnlock()
	for _, fn := range fns {
		if fn != nil {
			fn(event, payload)
		}
	}
}

// Count returns the number of non-nil subscribers for event.
func (o *Observer) Count(event string) int {
	o.mu.RLock()
	defer o.mu.RUnlock()
	n := 0
	for _, fn := range o.subs[event] {
		if fn != nil {
			n++
		}
	}
	return n
}

// Clear removes all subscribers for event.
func (o *Observer) Clear(event string) {
	o.mu.Lock()
	defer o.mu.Unlock()
	delete(o.subs, event)
}
