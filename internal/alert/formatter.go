package alert

import (
	"fmt"
	"strings"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Format controls the output format for alerts.
type Format string

const (
	FormatText Format = "text"
	FormatJSON  Format = "json"
)

// Formatter renders a port alert as a string in the given format.
type Formatter struct {
	format Format
}

// NewFormatter returns a Formatter for the given format string.
// Falls back to FormatText if the value is unrecognised.
func NewFormatter(format string) *Formatter {
	f := Format(strings.ToLower(format))
	if f != FormatJSON {
		f = FormatText
	}
	return &Formatter{format: f}
}

// Render returns a formatted string for the given event and port.
func (f *Formatter) Render(event string, p scanner.PortState, ts time.Time) string {
	switch f.format {
	case FormatJSON:
		return fmt.Sprintf(
			`{"time":%q,"event":%q,"port":%d,"proto":%q,"pid":%d}`,
			ts.UTC().Format(time.RFC3339),
			event,
			p.Port,
			p.Proto,
			p.PID,
		)
	default:
		return fmt.Sprintf(
			"[%s] %s port=%d proto=%s pid=%d",
			ts.UTC().Format(time.RFC3339),
			event,
			p.Port,
			p.Proto,
			p.PID,
		)
	}
}
