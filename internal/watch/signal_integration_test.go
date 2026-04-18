package watch_test

import (
	"context"
	"syscall"
	"testing"
	"time"

	"github.com/example/portwatch/internal/watch"
)

func TestSignalHandlerIntegration(t *testing.T) {
	h := watch.NewSignalHandler(syscall.SIGUSR1)
	ctx, cancel := h.WithCancel(context.Background())
	defer cancel()

	go func() {
		time.Sleep(20 * time.Millisecond)
		// send signal to self
		syscall.Kill(syscall.Getpid(), syscall.SIGUSR1) //nolint:errcheck
	}()

	select {
	case <-ctx.Done():
		// success
	case <-time.After(500 * time.Millisecond):
		t.Fatal("context not cancelled after SIGUSR1")
	}
}
