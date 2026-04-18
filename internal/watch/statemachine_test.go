package watch

import (
	"sync"
	"testing"
)

const (
	stateIdle    State = "idle"
	stateRunning State = "running"
	stateStopped State = "stopped"
)

func makeStateMachine() *StateMachine {
	return NewStateMachine(stateIdle, []Transition{
		{From: stateIdle, To: stateRunning},
		{From: stateRunning, To: stateStopped},
		{From: stateStopped, To: stateIdle},
	})
}

func TestStateMachineInitialState(t *testing.T) {
	sm := makeStateMachine()
	if sm.Current() != stateIdle {
		t.Fatalf("expected idle, got %s", sm.Current())
	}
}

func TestStateMachineValidTransition(t *testing.T) {
	sm := makeStateMachine()
	if !sm.Transition(stateRunning) {
		t.Fatal("expected transition to succeed")
	}
	if !sm.Is(stateRunning) {
		t.Fatalf("expected running, got %s", sm.Current())
	}
}

func TestStateMachineInvalidTransition(t *testing.T) {
	sm := makeStateMachine()
	if sm.Transition(stateStopped) {
		t.Fatal("expected transition to fail")
	}
	if !sm.Is(stateIdle) {
		t.Fatal("state should remain idle")
	}
}

func TestStateMachineOnEnterCallback(t *testing.T) {
	sm := makeStateMachine()
	var got State
	sm.OnEnter(stateRunning, func(from State) { got = from })
	sm.Transition(stateRunning)
	if got != stateIdle {
		t.Fatalf("expected callback from idle, got %s", got)
	}
}

func TestStateMachineConcurrentTransitions(t *testing.T) {
	sm := NewStateMachine(stateIdle, []Transition{
		{From: stateIdle, To: stateRunning},
		{From: stateRunning, To: stateIdle},
	})
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sm.Transition(stateRunning)
			sm.Transition(stateIdle)
		}()
	}
	wg.Wait()
}
