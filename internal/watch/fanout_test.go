package watch

import (
	"testing"
	"time"
)

func TestFanOutNoSubscribers(t *testing.T) {
	f := NewFanOut[int](2)
	// Should not panic with no subscribers.
	f.Publish(42)
}

func TestFanOutSingleSubscriber(t *testing.T) {
	f := NewFanOut[int](4)
	ch := f.Subscribe()
	f.Publish(7)
	select {
	case v := <-ch:
		if v != 7 {
			t.Fatalf("expected 7, got %d", v)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("timed out waiting for value")
	}
}

func TestFanOutMultipleSubscribers(t *testing.T) {
	f := NewFanOut[string](4)
	a := f.Subscribe()
	b := f.Subscribe()
	f.Publish("hello")
	for _, ch := range []<-chan string{a, b} {
		select {
		case v := <-ch:
			if v != "hello" {
				t.Fatalf("expected hello, got %s", v)
			}
		case <-time.After(100 * time.Millisecond):
			t.Fatal("timed out waiting for value")
		}
	}
}

func TestFanOutFullChannelSkipped(t *testing.T) {
	f := NewFanOut[int](1)
	ch := f.Subscribe()
	f.Publish(1) // fills buffer
	f.Publish(2) // should be skipped, not block
	if len(ch) != 1 {
		t.Fatalf("expected 1 item in channel, got %d", len(ch))
	}
}

func TestFanOutCloseSignalsSubscribers(t *testing.T) {
	f := NewFanOut[int](2)
	ch := f.Subscribe()
	f.Close()
	_, ok := <-ch
	if ok {
		t.Fatal("expected channel to be closed")
	}
}

func TestNewFanOutDefaultsBufSize(t *testing.T) {
	f := NewFanOut[int](0)
	if f.bufSize != 1 {
		t.Fatalf("expected bufSize 1, got %d", f.bufSize)
	}
}
