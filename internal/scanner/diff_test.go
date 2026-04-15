package scanner_test

import (
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

var (
	port80  = scanner.PortState{Protocol: "tcp", Port: 80, Address: "127.0.0.1"}
	port443 = scanner.PortState{Protocol: "tcp", Port: 443, Address: "127.0.0.1"}
	port22  = scanner.PortState{Protocol: "tcp", Port: 22, Address: "127.0.0.1"}
)

func TestCompareNoChanges(t *testing.T) {
	prev := []scanner.PortState{port80, port443}
	curr := []scanner.PortState{port80, port443}

	d := scanner.Compare(prev, curr)
	if d.HasChanges() {
		t.Errorf("expected no changes, got opened=%v closed=%v", d.Opened, d.Closed)
	}
}

func TestCompareNewPort(t *testing.T) {
	prev := []scanner.PortState{port80}
	curr := []scanner.PortState{port80, port443}

	d := scanner.Compare(prev, curr)
	if !d.HasChanges() {
		t.Fatal("expected changes but got none")
	}
	if len(d.Opened) != 1 || d.Opened[0].Port != 443 {
		t.Errorf("expected port 443 opened, got %v", d.Opened)
	}
	if len(d.Closed) != 0 {
		t.Errorf("expected no closed ports, got %v", d.Closed)
	}
}

func TestCompareClosedPort(t *testing.T) {
	prev := []scanner.PortState{port80, port22}
	curr := []scanner.PortState{port80}

	d := scanner.Compare(prev, curr)
	if !d.HasChanges() {
		t.Fatal("expected changes but got none")
	}
	if len(d.Closed) != 1 || d.Closed[0].Port != 22 {
		t.Errorf("expected port 22 closed, got %v", d.Closed)
	}
	if len(d.Opened) != 0 {
		t.Errorf("expected no opened ports, got %v", d.Opened)
	}
}

func TestCompareFromEmpty(t *testing.T) {
	prev := []scanner.PortState{}
	curr := []scanner.PortState{port80, port443}

	d := scanner.Compare(prev, curr)
	if len(d.Opened) != 2 {
		t.Errorf("expected 2 opened ports, got %d", len(d.Opened))
	}
}
