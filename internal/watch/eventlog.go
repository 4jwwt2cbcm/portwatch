package watch

import (
	"sync"
	"time"
)

// EventKind classifies a log entry.
type EventKind string

const (
	EventScanOK    EventKind = "scan_ok"
	EventScanError EventKind = "scan_error"
	EventPortAdded EventKind = "port_added"
	EventPortGone  EventKind = "port_gone"
)

// Event is a single structured log entry.
type Event struct {
	Kind    EventKind
	Message string
	At      time.Time
}

// EventLog holds a bounded ring of recent events.
type EventLog struct {
	mu     sync.Mutex
	events []Event
	cap    int
}

// NewEventLog returns an EventLog that retains at most cap entries.
func NewEventLog(cap int) *EventLog {
	if cap <= 0 {
		cap = 100
	}
	return &EventLog{cap: cap}
}

// Append adds an event, evicting the oldest when the log is full.
func (l *EventLog) Append(kind EventKind, msg string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if len(l.events) >= l.cap {
		l.events = l.events[1:]
	}
	l.events = append(l.events, Event{Kind: kind, Message: msg, At: time.Now()})
}

// Snapshot returns a copy of all current events.
func (l *EventLog) Snapshot() []Event {
	l.mu.Lock()
	defer l.mu.Unlock()
	out := make([]Event, len(l.events))
	copy(out, l.events)
	return out
}

// Len returns the current number of stored events.
func (l *EventLog) Len() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return len(l.events)
}

// Clear removes all events.
func (l *EventLog) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.events = l.events[:0]
}
