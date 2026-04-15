package scanner_test

import (
	"net"
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func TestScanDetectsOpenPort(t *testing.T) {
	// Start a real TCP listener on an ephemeral port.
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start listener: %v", err)
	}
	defer ln.Close()

	port := ln.Addr().(*net.TCPAddr).Port

	s := scanner.NewScanner()
	s.PortRange = [2]int{port, port}
	s.Protocols = []string{"tcp"}

	ports, err := s.Scan()
	if err != nil {
		t.Fatalf("Scan() returned error: %v", err)
	}

	if len(ports) != 1 {
		t.Fatalf("expected 1 open port, got %d", len(ports))
	}

	if ports[0].Port != port {
		t.Errorf("expected port %d, got %d", port, ports[0].Port)
	}
	if ports[0].Protocol != "tcp" {
		t.Errorf("expected protocol tcp, got %s", ports[0].Protocol)
	}
}

func TestScanClosedPort(t *testing.T) {
	// Port 1 is almost certainly closed in test environments.
	s := scanner.NewScanner()
	s.PortRange = [2]int{1, 1}
	s.Protocols = []string{"tcp"}

	ports, err := s.Scan()
	if err != nil {
		t.Fatalf("Scan() returned error: %v", err)
	}
	if len(ports) != 0 {
		t.Errorf("expected 0 open ports on port 1, got %d", len(ports))
	}
}

func TestPortStateString(t *testing.T) {
	p := scanner.PortState{Protocol: "tcp", Port: 8080, Address: "127.0.0.1"}
	got := p.String()
	want := "127.0.0.1:8080 (tcp)"
	if got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}
