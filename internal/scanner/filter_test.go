package scanner

import (
	"testing"
)

func makePort(proto string, port uint16) PortState {
	return PortState{Protocol: proto, Port: port}
}

func TestFilterAcceptAllByDefault(t *testing.T) {
	f := NewFilter()
	ports := []PortState{
		makePort("tcp", 80),
		makePort("udp", 53),
		makePort("tcp", 443),
	}
	got := f.Apply(ports)
	if len(got) != len(ports) {
		t.Fatalf("expected %d ports, got %d", len(ports), len(got))
	}
}

func TestFilterByProtocol(t *testing.T) {
	f := NewFilter().WithProtocols("tcp")
	ports := []PortState{
		makePort("tcp", 80),
		makePort("udp", 53),
		makePort("tcp", 443),
	}
	got := f.Apply(ports)
	if len(got) != 2 {
		t.Fatalf("expected 2 tcp ports, got %d", len(got))
	}
	for _, p := range got {
		if p.Protocol != "tcp" {
			t.Errorf("unexpected protocol %q", p.Protocol)
		}
	}
}

func TestFilterByPortRange(t *testing.T) {
	f := NewFilter().WithPortRange(1024, 8080)
	ports := []PortState{
		makePort("tcp", 80),
		makePort("tcp", 3000),
		makePort("tcp", 8080),
		makePort("tcp", 9000),
	}
	got := f.Apply(ports)
	if len(got) != 2 {
		t.Fatalf("expected 2 ports in range, got %d", len(got))
	}
}

func TestFilterCombinedProtocolAndRange(t *testing.T) {
	f := NewFilter().WithProtocols("udp").WithPortRange(50, 100)
	ports := []PortState{
		makePort("udp", 53),
		makePort("tcp", 80),
		makePort("udp", 123),
		makePort("udp", 67),
	}
	got := f.Apply(ports)
	if len(got) != 2 {
		t.Fatalf("expected 2 ports, got %d", len(got))
	}
}

func TestFilterEmptyInput(t *testing.T) {
	f := NewFilter().WithProtocols("tcp")
	got := f.Apply([]PortState{})
	if len(got) != 0 {
		t.Fatalf("expected empty result, got %d", len(got))
	}
}

func TestFilterAcceptDirectly(t *testing.T) {
	f := NewFilter().WithPortRange(1, 1024)
	if !f.Accept(makePort("tcp", 443)) {
		t.Error("expected port 443 to be accepted")
	}
	if f.Accept(makePort("tcp", 8443)) {
		t.Error("expected port 8443 to be rejected")
	}
}
