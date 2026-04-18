package watch

import "sync"

// Gate is a togglable boolean guard that allows or blocks passage.
// It is safe for concurrent use.
type Gate struct {
	mu     sync.RWMutex
	open   bool
	onOpen  func()
	onClose func()
}

// NewGate returns a Gate. If startOpen is true the gate begins open.
func NewGate(startOpen bool) *Gate {
	return &Gate{open: startOpen}
}

// OnOpen registers a callback invoked when the gate transitions to open.
func (g *Gate) OnOpen(fn func()) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.onOpen = fn
}

// OnClose registers a callback invoked when the gate transitions to closed.
func (g *Gate) OnClose(fn func()) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.onClose = fn
}

// Open opens the gate and fires the onOpen callback if registered.
func (g *Gate) Open() {
	g.mu.Lock()
	already := g.open
	g.open = true
	fn := g.onOpen
	g.mu.Unlock()
	if !already && fn != nil {
		fn()
	}
}

// Close closes the gate and fires the onClose callback if registered.
func (g *Gate) Close() {
	g.mu.Lock()
	already := !g.open
	g.open = false
	fn := g.onClose
	g.mu.Unlock()
	if !already && fn != nil {
		fn()
	}
}

// IsOpen reports whether the gate is currently open.
func (g *Gate) IsOpen() bool {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.open
}

// Allow returns true if the gate is open.
func (g *Gate) Allow() bool {
	return g.IsOpen()
}
