package alert

import (
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

var fixedTime = time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)

func makeTestPort() scanner.PortState {
	return scanner.PortState{Port: 8080, Proto: "tcp", PID: 1234}
}

func TestNewFormatterDefaultsToText(t *testing.T) {
	f := NewFormatter("unknown")
	if f.format != FormatText {
		t.Errorf("expected text format, got %q", f.format)
	}
}

func TestNewFormatterJSON(t *testing.T) {
	f := NewFormatter("json")
	if f.format != FormatJSON {
		t.Errorf("expected json format, got %q", f.format)
	}
}

func TestRenderTextFormat(t *testing.T) {
	f := NewFormatter("text")
	out := f.Render("ADDED", makeTestPort(), fixedTime)

	if !strings.Contains(out, "ADDED") {
		t.Errorf("expected ADDED in output: %s", out)
	}
	if !strings.Contains(out, "port=8080") {
		t.Errorf("expected port=8080 in output: %s", out)
	}
	if !strings.Contains(out, "proto=tcp") {
		t.Errorf("expected proto=tcp in output: %s", out)
	}
	if !strings.Contains(out, "pid=1234") {
		t.Errorf("expected pid=1234 in output: %s", out)
	}
}

func TestRenderJSONFormat(t *testing.T) {
	f := NewFormatter("json")
	out := f.Render("REMOVED", makeTestPort(), fixedTime)

	if !strings.HasPrefix(out, "{") || !strings.HasSuffix(out, "}") {
		t.Errorf("expected JSON object, got: %s", out)
	}
	if !strings.Contains(out, `"event":"REMOVED"`) {
		t.Errorf("expected event field in JSON: %s", out)
	}
	if !strings.Contains(out, `"port":8080`) {
		t.Errorf("expected port field in JSON: %s", out)
	}
	if !strings.Contains(out, `"proto":"tcp"`) {
		t.Errorf("expected proto field in JSON: %s", out)
	}
}

func TestRenderTimestampIncluded(t *testing.T) {
	f := NewFormatter("text")
	out := f.Render("ADDED", makeTestPort(), fixedTime)
	if !strings.Contains(out, "2024-06-01T12:00:00Z") {
		t.Errorf("expected timestamp in output: %s", out)
	}
}
