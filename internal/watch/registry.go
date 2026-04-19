package watch

import (
	"fmt"
	"sync"
)

// Registry holds named runners and allows starting/stopping them by name.
type Registry struct {
	mu      sync.RWMutex
	entries map[string]func() error
}

// NewRegistry returns an empty Registry.
func NewRegistry() *Registry {
	return &Registry{
		entries: make(map[string]func() error),
	}
}

// Register adds a named runner function to the registry.
// Returns an error if the name is already registered.
func (r *Registry) Register(name string, fn func() error) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.entries[name]; exists {
		return fmt.Errorf("registry: %q already registered", name)
	}
	r.entries[name] = fn
	return nil
}

// Unregister removes a named runner from the registry.
func (r *Registry) Unregister(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.entries, name)
}

// Run executes the runner registered under name.
// Returns an error if the name is not found or the runner fails.
func (r *Registry) Run(name string) error {
	r.mu.RLock()
	fn, ok := r.entries[name]
	r.mu.RUnlock()
	if !ok {
		return fmt.Errorf("registry: %q not found", name)
	}
	return fn()
}

// Names returns a sorted snapshot of registered names.
func (r *Registry) Names() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	names := make([]string, 0, len(r.entries))
	for k := range r.entries {
		names = append(names, k)
	}
	return names
}

// Count returns the number of registered runners.
func (r *Registry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.entries)
}
