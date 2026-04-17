package watch

import "sync"

// HookEvent represents the type of lifecycle event.
type HookEvent string

const (
	HookBeforeScan HookEvent = "before_scan"
	HookAfterScan  HookEvent = "after_scan"
	HookOnError    HookEvent = "on_error"
)

// HookFunc is a callback invoked on a lifecycle event.
type HookFunc func(event HookEvent, meta map[string]any)

// HookRegistry manages lifecycle hooks for the watcher.
type HookRegistry struct {
	mu    sync.RWMutex
	hooks map[HookEvent][]HookFunc
}

// NewHookRegistry returns an empty HookRegistry.
func NewHookRegistry() *HookRegistry {
	return &HookRegistry{
		hooks: make(map[HookEvent][]HookFunc),
	}
}

// Register adds a hook for the given event.
func (r *HookRegistry) Register(event HookEvent, fn HookFunc) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.hooks[event] = append(r.hooks[event], fn)
}

// Fire invokes all hooks registered for the given event.
func (r *HookRegistry) Fire(event HookEvent, meta map[string]any) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, fn := range r.hooks[event] {
		fn(event, meta)
	}
}

// Count returns the number of hooks registered for an event.
func (r *HookRegistry) Count(event HookEvent) int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.hooks[event])
}

// Clear removes all hooks for the given event.
func (r *HookRegistry) Clear(event HookEvent) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.hooks, event)
}
