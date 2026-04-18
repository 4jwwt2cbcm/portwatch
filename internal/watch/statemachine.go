package watch

import "sync"

// State represents a named state in a state machine.
type State string

// Transition defines a valid move from one state to another.
type Transition struct {
	From State
	To   State
}

// StateMachine is a simple thread-safe finite state machine.
type StateMachine struct {
	mu       sync.RWMutex
	current  State
	allowed  map[Transition]struct{}
	onEnter  map[State]func(from State)
}

// NewStateMachine creates a StateMachine with the given initial state and
// allowed transitions.
func NewStateMachine(initial State, transitions []Transition) *StateMachine {
	allowed := make(map[Transition]struct{}, len(transitions))
	for _, t := range transitions {
		allowed[t] = struct{}{}
	}
	return &StateMachine{
		current: initial,
		allowed: allowed,
		onEnter: make(map[State]func(from State)),
	}
}

// Current returns the current state.
func (sm *StateMachine) Current() State {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.current
}

// OnEnter registers a callback invoked whenever the machine enters state s.
func (sm *StateMachine) OnEnter(s State, fn func(from State)) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.onEnter[s] = fn
}

// Transition attempts to move to next. Returns false if the transition is not
// allowed.
func (sm *StateMachine) Transition(next State) bool {
	sm.mu.Lock()
	t := Transition{From: sm.current, To: next}
	if _, ok := sm.allowed[t]; !ok {
		sm.mu.Unlock()
		return false
	}
	prev := sm.current
	sm.current = next
	cb := sm.onEnter[next]
	sm.mu.Unlock()
	if cb != nil {
		cb(prev)
	}
	return true
}

// Is returns true when the current state equals s.
func (sm *StateMachine) Is(s State) bool {
	return sm.Current() == s
}
