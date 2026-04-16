package watch

import (
	"fmt"
	"io"
	"os"
	"time"
)

// EventLogWriter writes EventLog snapshots to an io.Writer on a ticker.
type EventLogWriter struct {
	log    *EventLog
	out    io.Writer
	ticker *time.Ticker
}

// NewEventLogWriter returns an EventLogWriter that flushes every interval.
func NewEventLogWriter(log *EventLog, interval time.Duration, out io.Writer) *EventLogWriter {
	if out == nil {
		out = os.Stderr
	}
	return &EventLogWriter{
		log:    log,
		out:    out,
		ticker: time.NewTicker(interval),
	}
}

// Run writes the event log on each tick until ctx is done.
func (w *EventLogWriter) Run(ctx interface{ Done() <-chan struct{} }) {
	for {
		select {
		case <-ctx.Done():
			w.ticker.Stop()
			return
		case <-w.ticker.C:
			w.flush()
		}
	}
}

func (w *EventLogWriter) flush() {
	events := w.log.Snapshot()
	if len(events) == 0 {
		return
	}
	for _, e := range events {
		fmt.Fprintf(w.out, "%s [%s] %s\n",
			e.At.Format(time.RFC3339), e.Kind, e.Message)
	}
}
