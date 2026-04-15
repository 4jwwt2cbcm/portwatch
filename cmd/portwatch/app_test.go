package main

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/portwatch/internal/config"
)

func testConfig(t *testing.T) *config.Config {
	t.Helper()
	cfg := config.DefaultConfig()
	cfg.StateFile = filepath.Join(t.TempDir(), "state.json")
	cfg.Ports = []int{}
	cfg.IntervalSeconds = 1
	return cfg
}

func TestNewAppSucceeds(t *testing.T) {
	cfg := testConfig(t)
	app, err := NewApp(cfg)
	if err != nil {
		t.Fatalf("NewApp returned error: %v", err)
	}
	if app == nil {
		t.Fatal("expected non-nil app")
	}
}

func TestRunExitsOnContextCancel(t *testing.T) {
	cfg := testConfig(t)
	app, err := NewApp(cfg)
	if err != nil {
		t.Fatalf("NewApp: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan struct{})
	go func() {
		app.Run(ctx)
		close(done)
	}()

	cancel()

	select {
	case <-done:
		// success
	case <-time.After(2 * time.Second):
		t.Fatal("Run did not exit after context cancellation")
	}
}

func TestNewAppBadStateDir(t *testing.T) {
	cfg := testConfig(t)
	// Point state file at a path whose parent does not exist.
	cfg.StateFile = "/nonexistent/dir/state.json"
	_, err := NewApp(cfg)
	if err == nil {
		t.Fatal("expected error for bad state path, got nil")
	}
}
