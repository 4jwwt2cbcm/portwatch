package watch

import (
	"context"
	"os"
	"syscall"
	"testing"
	"time"
)

func makeSignalHandler(sigs ...os.Signal) (*SignalHandler, chan os.Signal) {
	ch := make(chan os.Signal, 1)
	h := NewSignalHandler(sigs...)
	h.notify = func(c chan<- os.Signal, s ...os.Signal) {
		go func() {
			for sig := range ch {
				c <- sig
			}
		}()
	}
	h.stop = func(c chan<- os.Signal) {}
	return h, ch
}

func TestSignalHandlerDefaultSignals(t *testing.T) {
	h := NewSignalHandler()
	if len(h.signals) != 2 {
		t.Fatalf("expected 2 default signals, got %d", len(h.signals))
	}
}

func TestSignalHandlerRunReceivesSignal(t *testing.T) {
	h, ch := makeSignalHandler(syscall.SIGINT)
	ctx := context.Background()

	go func() {
		time.Sleep(10 * time.Millisecond)
		ch <- syscall.SIGINT
	}()

	sig := h.Run(ctx)
	if sig != syscall.SIGINT {
		t.Fatalf("expected SIGINT, got %v", sig)
	}
}

func TestSignalHandlerRunCancelledContext(t *testing.T) {
	h, _ := makeSignalHandler(syscall.SIGINT)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	sig := h.Run(ctx)
	if sig != nil {
		t.Fatalf("expected nil signal on cancelled context, got %v", sig)
	}
}

func TestSignalHandlerWithCancel(t *testing.T) {
	h, ch := makeSignalHandler(syscall.SIGTERM)
	ctx, cancel := h.WithCancel(context.Background())
	defer cancel()

	ch <- syscall.SIGTERM

	select {
	case <-ctx.Done():
		// expected
	case <-time.After(200 * time.Millisecond):
		t.Fatal("context was not cancelled after signal")
	}
}

func TestSignalHandlerWithCancelParentCancels(t *testing.T) {
	h, _ := makeSignalHandler(syscall.SIGTERM)
	parent, parentCancel := context.WithCancel(context.Background())
	ctx, cancel := h.WithCancel(parent)
	defer cancel()

	parentCancel()

	select {
	case <-ctx.Done():
		// expected
	case <-time.After(200 * time.Millisecond):
		t.Fatal("context was not cancelled after parent cancel")
	}
}

func TestSignalHandlerWithCancelCancelFuncStopsRun(t *testing.T) {
	h, _ := makeSignalHandler(syscall.SIGTERM)
	ctx, cancel := h.WithCancel(context.Background())

	cancel()

	select {
	case <-ctx.Done():
		// expected: calling cancel() should stop the context
	case <-time.After(200 * time.Millisecond):
		t.Fatal("context was not cancelled after calling cancel")
	}
}
