package main

import (
	"testing"
)

func TestParseFlagsDefaults(t *testing.T) {
	f := ParseFlags([]string{})
	if f.ConfigPath != "portwatch.yaml" {
		t.Errorf("expected default config path, got %q", f.ConfigPath)
	}
	if f.Version {
		t.Error("expected Version to be false by default")
	}
}

func TestParseFlagsCustomConfig(t *testing.T) {
	f := ParseFlags([]string{"-config", "/etc/portwatch/config.yaml"})
	if f.ConfigPath != "/etc/portwatch/config.yaml" {
		t.Errorf("unexpected config path: %q", f.ConfigPath)
	}
}

func TestParseFlagsVersionFlag(t *testing.T) {
	f := ParseFlags([]string{"-version"})
	if !f.Version {
		t.Error("expected Version to be true")
	}
}
