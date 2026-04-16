package watch

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func makePort(port int, proto string) scanner.PortState {
	return scanner.PortState{Port: port, Protocol: proto, Open: true}
}

func TestSnapshotStoreEmptyOnInit(t *testing.T) {
	s := NewSnapshotStore()
	_, ok := s.Get()
	if ok {
		t.Fatal("expected no snapshot on init")
	}
}

func TestSnapshotStoreSetAndGet(t *testing.T) {
	s := NewSnapshotStore()
	ports := []scanner.PortState{makePort(80, "tcp"), makePort(443, "tcp")}
	before := time.Now()
	s.Set(ports)
	snap, ok := s.Get()
	if !ok {
		t.Fatal("expected snapshot after Set")
	}
	if len(snap.Ports) != 2 {
		t.Fatalf("expected 2 ports, got %d", len(snap.Ports))
	}
	if snap.CapturedAt.Before(before) {
		t.Error("CapturedAt should be after test start")
	}
}

func TestSnapshotStoreGetReturnsCopy(t *testing.T) {
	s := NewSnapshotStore()
	s.Set([]scanner.PortState{makePort(22, "tcp")})
	snap, _ := s.Get()
	snap.Ports[0].Port = 9999
	snap2, _ := s.Get()
	if snap2.Ports[0].Port == 9999 {
		t.Error("Get should return a copy, not a reference")
	}
}

func TestSnapshotStoreClear(t *testing.T) {
	s := NewSnapshotStore()
	s.Set([]scanner.PortState{makePort(8080, "tcp")})
	s.Clear()
	_, ok := s.Get()
	if ok {
		t.Error("expected no snapshot after Clear")
	}
}

func TestSnapshotStoreOverwrite(t *testing.T) {
	s := NewSnapshotStore()
	s.Set([]scanner.PortState{makePort(80, "tcp")})
	s.Set([]scanner.PortState{makePort(443, "tcp"), makePort(8443, "tcp")})
	snap, _ := s.Get()
	if len(snap.Ports) != 2 {
		t.Fatalf("expected 2 ports after overwrite, got %d", len(snap.Ports))
	}
}
