package watch

import "errors"

// ErrHighWatermark is returned by WatermarkRunner when the watermark is above
// the high threshold and the function call is suppressed.
var ErrHighWatermark = errors.New("watermark: above high watermark")

// WatermarkRunner wraps a Watermark and guards a function call: if the
// watermark is currently above the high threshold the function is skipped and
// ErrHighWatermark is returned instead.
type WatermarkRunner struct {
	wm *Watermark
	fn func() error
}

// NewWatermarkRunner creates a WatermarkRunner.
// If wm is nil a default Watermark is used.
// If fn is nil a no-op function is used.
func NewWatermarkRunner(wm *Watermark, fn func() error) *WatermarkRunner {
	if wm == nil {
		wm = NewWatermark(DefaultWatermarkPolicy())
	}
	if fn == nil {
		fn = func() error { return nil }
	}
	return &WatermarkRunner{wm: wm, fn: fn}
}

// Run executes the wrapped function only when the watermark is not above the
// high threshold. If the watermark is above, ErrHighWatermark is returned and
// the function is not called.
func (r *WatermarkRunner) Run() error {
	if r.wm.Above() {
		return ErrHighWatermark
	}
	return r.fn()
}

// Watermark returns the underlying Watermark so callers can update the level.
func (r *WatermarkRunner) Watermark() *Watermark {
	return r.wm
}
