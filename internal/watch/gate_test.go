package watch

import (
	"sync"
	"testing"
)

func TestGateClosedByDefault(t *testing.T) {
	g := NewGate(false)
	if g.IsOpen() {
		t.Fatal("expected gate to be closed")
	}
}

func TestGateOpenByDefault(t *testing.T) {
	g := NewGate(true)
	if !g.IsOpen() {
		t.Fatal("expected gate to be open")
	}
}

func TestGateOpenAllows(t *testing.T) {
	g := NewGate(false)
	g.Open()
	if !g.Allow() {
		t.Fatal("expected Allow to return true after Open")
	}
}

func TestGateCloseBlocks(t *testing.T) {
	g := NewGate(true)
	g.Close()
	if g.Allow() {
		t.Fatal("expected Allow to return false after Close")
	}
}

func TestGateOnOpenCallback(t *testing.T) {
	g := NewGate(false)
	called := false
	g.OnOpen(func() { called = true })
	g.Open()
	if !called {
		t.Fatal("expected onOpen callback to be called")
	}
}

func TestGateOnOpenCallbackNotFiredIfAlreadyOpen(t *testing.T) {
	g := NewGate(true)
	called := false
	g.OnOpen(func() { called = true })
	g.Open()
	if called {
		t.Fatal("onOpen should not fire if gate was already open")
	}
}

func TestGateOnCloseCallback(t *testing.T) {
	g := NewGate(true)
	called := false
	g.OnClose(func() { called = true })
	g.Close()
	if !called {
		t.Fatal("expected onClose callback to be called")
	}
}

func TestGateConcurrentAccess(t *testing.T) {
	g := NewGate(false)
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(2)
		go func() { defer wg.Done(); g.Open() }()
		go func() { defer wg.Done(); _ = g.IsOpen() }()
	}
	wg.Wait()
}
