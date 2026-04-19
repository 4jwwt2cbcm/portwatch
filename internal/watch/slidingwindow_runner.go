package watch

import "errors"

// ErrWindowFull is returned when the sliding window threshold is exceeded.
var ErrWindowFull = errors.New("sliding window: threshold exceeded")

// SlidingWindowRunner wraps a function and gates execution through a
// SlidingWindow, returning ErrWindowFull when the window is saturated.
type SlidingWindowRunner struct {
	win *SlidingWindow
	fn  func() error
}

// NewSlidingWindowRunner creates a SlidingWindowRunner.
// If win is nil a default window is used. If fn is nil a noop is used.
func NewSlidingWindowRunner(win *SlidingWindow, fn func() error) *SlidingWindowRunner {
	if win == nil {
		win = NewSlidingWindow(DefaultSlidingWindowPolicy())
	}
	if fn == nil {
		fn = func() error { return nil }
	}
	return &SlidingWindowRunner{win: win, fn: fn}
}

// Run records the event in the window. If the window is full it returns
// ErrWindowFull without invoking fn. Otherwise fn is called and its
// error returned.
func (r *SlidingWindowRunner) Run() error {
	if r.win.Record() {
		return ErrWindowFull
	}
	return r.fn()
}
