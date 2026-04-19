package watch

import (
	"context"
	"testing"
	"time"
)

func TestDrainEmptyOnInit(t *testing.T) {
	d := NewDrain[int]()
	if d.Len() != 0 {
		t.Fatalf("expected 0, got %d", d.Len())
	}
}

func TestDrainCollectsItems(t *testing.T) {
	d := NewDrain[int]()
	ch := make(chan int, 3)
	ch <- 1
	ch <- 2
	ch <- 3
	close(ch)

	d.Run(context.Background(), ch)

	if d.Len() != 3 {
		t.Fatalf("expected 3 items, got %d", d.Len())
	}
	snap := d.Snapshot()
	for i, v := range snap {
		if v != i+1 {
			t.Errorf("index %d: expected %d, got %d", i, i+1, v)
		}
	}
}

func TestDrainSnapshotReturnsCopy(t *testing.T) {
	d := NewDrain[int]()
	ch := make(chan int, 1)
	ch <- 42
	close(ch)
	d.Run(context.Background(), ch)

	snap := d.Snapshot()
	snap[0] = 99
	if d.Snapshot()[0] != 42 {
		t.Fatal("snapshot mutation affected internal state")
	}
}

func TestDrainClear(t *testing.T) {
	d := NewDrain[int]()
	ch := make(chan int, 2)
	ch <- 1
	ch <- 2
	close(ch)
	d.Run(context.Background(), ch)
	d.Clear()
	if d.Len() != 0 {
		t.Fatalf("expected 0 after clear, got %d", d.Len())
	}
}

func TestDrainStopsOnContextCancel(t *testing.T) {
	d := NewDrain[int]()
	ch := make(chan int)

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		d.Run(ctx, ch)
		close(done)
	}()

	cancel()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("Run did not stop after context cancel")
	}
}
