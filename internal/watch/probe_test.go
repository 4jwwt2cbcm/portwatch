package watch

import (
	"context"
	"net"
	"strconv"
	"testing"
	"time"
)

func startTCPServer(t *testing.T) (port int, stop func()) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	_, p, _ := net.SplitHostPort(ln.Addr().String())
	port, _ = strconv.Atoi(p)
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			conn.Close()
		}
	}()
	return port, func() { ln.Close() }
}

func TestDefaultProbePolicy(t *testing.T) {
	p := DefaultProbePolicy()
	if p.Timeout != 2*time.Second {
		t.Errorf("expected 2s timeout, got %v", p.Timeout)
	}
	if p.Retries != 2 {
		t.Errorf("expected 2 retries, got %d", p.Retries)
	}
}

func TestProbeReachableOpenPort(t *testing.T) {
	port, stop := startTCPServer(t)
	defer stop()

	pr := NewProbe(DefaultProbePolicy())
	if !pr.Reachable(context.Background(), "127.0.0.1", port) {
		t.Error("expected port to be reachable")
	}
}

func TestProbeCheckClosedPort(t *testing.T) {
	pr := NewProbe(ProbePolicy{Timeout: 200 * time.Millisecond, Retries: 0})
	err := pr.Check(context.Background(), "127.0.0.1", 1)
	if err == nil {
		t.Error("expected error for closed port")
	}
}

func TestProbeRespectsContextCancel(t *testing.T) {
	pr := NewProbe(ProbePolicy{Timeout: 2 * time.Second, Retries: 0})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := pr.Check(ctx, "127.0.0.1", 9999)
	if err == nil {
		t.Error("expected error with cancelled context")
	}
}

func TestNewProbeDefaultsOnZeroTimeout(t *testing.T) {
	pr := NewProbe(ProbePolicy{Timeout: 0, Retries: 1})
	if pr.policy.Timeout != DefaultProbePolicy().Timeout {
		t.Errorf("expected default timeout, got %v", pr.policy.Timeout)
	}
}
