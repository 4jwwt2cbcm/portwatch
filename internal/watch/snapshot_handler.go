package watch

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"text/tabwriter"
	"time"
)

// SnapshotHandler formats and writes a Snapshot to an output.
type SnapshotHandler struct {
	out    io.Writer
	format string
}

// NewSnapshotHandler returns a SnapshotHandler writing to out in the given format ("text" or "json").
func NewSnapshotHandler(out io.Writer, format string) *SnapshotHandler {
	if out == nil {
		out = os.Stdout
	}
	if format == "" {
		format = "text"
	}
	return &SnapshotHandler{out: out, format: format}
}

// Write renders the snapshot to the configured output.
func (h *SnapshotHandler) Write(snap *Snapshot) error {
	if snap == nil {
		return fmt.Errorf("snapshot is nil")
	}
	switch h.format {
	case "json":
		return h.writeJSON(snap)
	default:
		return h.writeText(snap)
	}
}

func (h *SnapshotHandler) writeText(snap *Snapshot) error {
	w := tabwriter.NewWriter(h.out, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "Snapshot at %s\n", snap.CapturedAt.Format(time.RFC3339))
	fmt.Fprintln(w, "PORT\tPROTOCOL\tSTATE")
	for _, p := range snap.Ports {
		state := "closed"
		if p.Open {
			state = "open"
		}
		fmt.Fprintf(w, "%d\t%s\t%s\n", p.Port, p.Protocol, state)
	}
	return w.Flush()
}

type jsonSnapshot struct {
	CapturedAt string            `json:"captured_at"`
	Ports      []jsonPort        `json:"ports"`
}

type jsonPort struct {
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
	Open     bool   `json:"open"`
}

func (h *SnapshotHandler) writeJSON(snap *Snapshot) error {
	out := jsonSnapshot{
		CapturedAt: snap.CapturedAt.Format(time.RFC3339),
	}
	for _, p := range snap.Ports {
		out.Ports = append(out.Ports, jsonPort{Port: p.Port, Protocol: p.Protocol, Open: p.Open})
	}
	enc := json.NewEncoder(h.out)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
