package watch

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// SignalHandler listens for OS signals and cancels a context.
type SignalHandler struct {
	signals []os.Signal
	notify  func(chan<- os.Signal, ...os.Signal)
	stop    func(chan<- os.Signal)
}

// NewSignalHandler creates a SignalHandler for SIGINT and SIGTERM by default.
func NewSignalHandler(sigs ...os.Signal) *SignalHandler {
	if len(sigs) == 0 {
		sigs = []os.Signal{syscall.SIGINT, syscall.SIGTERM}
	}
	return &SignalHandler{
		signals: sigs,
		notify:  signal.Notify,
		stop:    signal.Stop,
	}
}

// Run blocks until a signal is received or ctx is cancelled.
// Returns the received signal, or nil if context was cancelled first.
func (h *SignalHandler) Run(ctx context.Context) os.Signal {
	ch := make(chan os.Signal, 1)
	h.notify(ch, h.signals...)
	defer h.stop(ch)

	select {
	case sig := <-ch:
		return sig
	case <-ctx.Done():
		return nil
	}
}

// WithCancel returns a derived context that is cancelled when a signal arrives.
func (h *SignalHandler) WithCancel(parent context.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(parent)
	go func() {
		if sig := h.Run(parent); sig != nil {
			_ = sig
			cancel()
		}
	}()
	return ctx, cancel
}
