package watch

import (
	"errors"
	"testing"
)

func TestTeeNoSinksIsNoop(t *testing.T) {
	tee := NewTee[int]()
	tee.Emit(42, nil) // should not panic
}

func TestTeeEmitCallsAllSinks(t *testing.T) {
	var got []int
	a := func(v int, _ error) { got = append(got, v*1) }
	b := func(v int, _ error) { got = append(got, v*2) }
	tee := NewTee(a, b)
	tee.Emit(3, nil)
	if len(got) != 2 || got[0] != 3 || got[1] != 6 {
		t.Fatalf("expected [3 6], got %v", got)
	}
}

func TestTeeEmitPassesError(t *testing.T) {
	sentinel := errors.New("boom")
	var gotErr error
	tee := NewTee(func(_ int, err error) { gotErr = err })
	tee.Emit(0, sentinel)
	if gotErr != sentinel {
		t.Fatalf("expected sentinel error, got %v", gotErr)
	}
}

func TestTeeAdd(t *testing.T) {
	tee := NewTee[string]()
	if tee.Count() != 0 {
		t.Fatal("expected 0 sinks")
	}
	tee.Add(func(_ string, _ error) {})
	if tee.Count() != 1 {
		t.Fatal("expected 1 sink")
	}
}

func TestTeeClear(t *testing.T) {
	tee := NewTee(func(_ int, _ error) {})
	tee.Clear()
	if tee.Count() != 0 {
		t.Fatal("expected 0 sinks after clear")
	}
}

func TestTeeWrap(t *testing.T) {
	var captured int
	tee := NewTee(func(v int, _ error) { captured = v })
	wrapped := tee.Wrap(func() (int, error) { return 99, nil })
	v, err := wrapped()
	if err != nil || v != 99 {
		t.Fatalf("unexpected v=%d err=%v", v, err)
	}
	if captured != 99 {
		t.Fatalf("expected captured=99, got %d", captured)
	}
}

func TestTeeWrapPropagatesError(t *testing.T) {
	sentinel := errors.New("fail")
	var gotErr error
	tee := NewTee(func(_ int, err error) { gotErr = err })
	wrapped := tee.Wrap(func() (int, error) { return 0, sentinel })
	_, err := wrapped()
	if err != sentinel || gotErr != sentinel {
		t.Fatalf("expected sentinel, got err=%v gotErr=%v", err, gotErr)
	}
}
