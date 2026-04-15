package config

import (
	"os"
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Interval != 30*time.Second {
		t.Errorf("expected default interval 30s, got %v", cfg.Interval)
	}
	if len(cfg.Ports) != 0 {
		t.Errorf("expected empty default ports, got %v", cfg.Ports)
	}
	if cfg.AlertOutput != "" {
		t.Errorf("expected empty alert output, got %q", cfg.AlertOutput)
	}
	if len(cfg.Protocols) != 1 || cfg.Protocols[0] != "tcp" {
		t.Errorf("expected default protocols [tcp], got %v", cfg.Protocols)
	}
}

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "portwatch-config-*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestLoadFullConfig(t *testing.T) {
	path := writeTemp(t, `{
		"interval_seconds": 60,
		"ports": [80, 443, 8080],
		"alert_output": "/var/log/portwatch.log",
		"protocols": ["tcp", "udp"]
	}`)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}
	if cfg.Interval != 60*time.Second {
		t.Errorf("expected 60s interval, got %v", cfg.Interval)
	}
	if len(cfg.Ports) != 3 || cfg.Ports[1] != 443 {
		t.Errorf("unexpected ports: %v", cfg.Ports)
	}
	if cfg.AlertOutput != "/var/log/portwatch.log" {
		t.Errorf("unexpected alert output: %q", cfg.AlertOutput)
	}
	if len(cfg.Protocols) != 2 {
		t.Errorf("expected 2 protocols, got %v", cfg.Protocols)
	}
}

func TestLoadPartialConfig(t *testing.T) {
	path := writeTemp(t, `{"interval_seconds": 10}`)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}
	if cfg.Interval != 10*time.Second {
		t.Errorf("expected 10s interval, got %v", cfg.Interval)
	}
	// Unset fields should retain defaults.
	if len(cfg.Protocols) != 1 || cfg.Protocols[0] != "tcp" {
		t.Errorf("expected default protocols, got %v", cfg.Protocols)
	}
}

func TestLoadMissingFile(t *testing.T) {
	_, err := Load("/nonexistent/path/portwatch.json")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}
