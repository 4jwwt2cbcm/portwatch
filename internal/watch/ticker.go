package watch

import "time"

// TickerFactory creates time.Ticker instances, allowing tests to substitute
// a fake clock without modifying production code.
type TickerFactory func(d time.Duration) *time.Ticker

// DefaultTickerFactory returns the standard time.NewTicker.
var DefaultTickerFactory TickerFactory = time.NewTicker

// IntervalFromSeconds converts an integer number of seconds into a
// time.Duration, clamping to a minimum of 1 second to prevent runaway loops.
func IntervalFromSeconds(secs int) time.Duration {
	if secs < 1 {
		return time.Second
	}
	return time.Duration(secs) * time.Second
}
