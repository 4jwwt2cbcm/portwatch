package alert

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Level represents the severity of an alert.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelAlert Level = "ALERT"
)

// Alert represents a single port change notification.
type Alert struct {
	Timestamp time.Time
	Level     Level
	Message   string
	Port      scanner.PortState
}

// Notifier sends alerts to a destination.
type Notifier struct {
	Writer io.Writer
}

// NewNotifier creates a Notifier that writes to the given writer.
// If w is nil, os.Stdout is used.
func NewNotifier(w io.Writer) *Notifier {
	if w == nil {
		w = os.Stdout
	}
	return &Notifier{Writer: w}
}

// Notify formats and writes an alert to the configured writer.
func (n *Notifier) Notify(a Alert) error {
	_, err := fmt.Fprintf(
		n.Writer,
		"[%s] %s %s\n",
		a.Timestamp.Format(time.RFC3339),
		a.Level,
		a.Message,
	)
	return err
}

// NewPortAlert constructs an Alert for a newly opened port.
func NewPortAlert(ps scanner.PortState) Alert {
	return Alert{
		Timestamp: time.Now(),
		Level:     LevelAlert,
		Message:   fmt.Sprintf("new port open: %s", ps),
		Port:      ps,
	}
}

// ClosedPortAlert constructs an Alert for a port that has closed.
func ClosedPortAlert(ps scanner.PortState) Alert {
	return Alert{
		Timestamp: time.Now(),
		Level:     LevelWarn,
		Message:   fmt.Sprintf("port closed: %s", ps),
		Port:      ps,
	}
}
