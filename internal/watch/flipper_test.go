package watch

import (
	"sync"
	"testing"
)

func TestFlipperInitialStateFalse(t *testing.T) {
	f := NewFlipper(false, nil)
	if f.State() != false {
		t.Fatalf("expected false, got true")
	}
}

func TestFlipperInitialStateTrue(t *testing.T) {
	f := NewFlipper(true, nil)
	if f.State() != true {
		t.Fatalf("expected true, got false")
	}
}

func TestFlipperFlipTogglesState(t *testing.T) {
	f := NewFlipper(false, nil)
	next := f.Flip()
	if next != true {
		t.Fatalf("expected Flip to return true")
	}
	if f.State() != true {
		t.Fatalf("expected state true after flip")
	}
}

func TestFlipperCountIncrementsOnFlip(t *testing.T) {
	f := NewFlipper(false, nil)
	f.Flip()
	f.Flip()
	f.Flip()
	if f.Count() != 3 {
		t.Fatalf("expected count 3, got %d", f.Count())
	}
}

func TestFlipperResetClearsCount(t *testing.T) {
	f := NewFlipper(false, nil)
	f.Flip()
	f.Flip()
	f.Reset(false)
	if f.Count() != 0 {
		t.Fatalf("expected count 0 after reset, got %d", f.Count())
	}
	if f.State() != false {
		t.Fatalf("expected state false after reset")
	}
}

func TestFlipperOnFlipCallbackFired(t *testing.T) {
	var got []bool
	f := NewFlipper(false, func(s bool) {
		got = append(got, s)
	})
	f.Flip()
	f.Flip()
	if len(got) != 2 {
		t.Fatalf("expected 2 callbacks, got %d", len(got))
	}
	if got[0] != true || got[1] != false {
		t.Fatalf("unexpected callback values: %v", got)
	}
}

func TestFlipperConcurrentFlips(t *testing.T) {
	f := NewFlipper(false, nil)
	var wg sync.WaitGroup
	const n = 100
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			f.Flip()
		}()
	}
	wg.Wait()
	if f.Count() != n {
		t.Fatalf("expected count %d, got %d", n, f.Count())
	}
}
