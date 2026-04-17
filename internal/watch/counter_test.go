package watch

import (
	"sync"
	"testing"
)

func TestCounterZeroOnInit(t *testing.T) {
	c := NewCounter()
	if c.Value() != 0 {
		t.Fatalf("expected 0, got %d", c.Value())
	}
}

func TestCounterInc(t *testing.T) {
	c := NewCounter()
	v := c.Inc()
	if v != 1 {
		t.Fatalf("expected 1, got %d", v)
	}
	if c.Value() != 1 {
		t.Fatalf("expected Value 1, got %d", c.Value())
	}
}

func TestCounterAdd(t *testing.T) {
	c := NewCounter()
	v := c.Add(5)
	if v != 5 {
		t.Fatalf("expected 5, got %d", v)
	}
	c.Add(3)
	if c.Value() != 8 {
		t.Fatalf("expected 8, got %d", c.Value())
	}
}

func TestCounterReset(t *testing.T) {
	c := NewCounter()
	c.Add(10)
	prev := c.Reset()
	if prev != 10 {
		t.Fatalf("expected prev 10, got %d", prev)
	}
	if c.Value() != 0 {
		t.Fatalf("expected 0 after reset, got %d", c.Value())
	}
}

func TestCounterConcurrentInc(t *testing.T) {
	c := NewCounter()
	var wg sync.WaitGroup
	const n = 100
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.Inc()
		}()
	}
	wg.Wait()
	if c.Value() != n {
		t.Fatalf("expected %d, got %d", n, c.Value())
	}
}
