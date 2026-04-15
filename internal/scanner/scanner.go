package scanner

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

// PortState represents the state of a single open port.
type PortState struct {
	Protocol string
	Port     int
	Address  string
}

func (p PortState) String() string {
	return fmt.Sprintf("%s:%d (%s)", p.Address, p.Port, p.Protocol)
}

// Scanner scans for open ports on the local machine.
type Scanner struct {
	PortRange [2]int
	Protocols []string
}

// NewScanner creates a Scanner with sensible defaults.
func NewScanner() *Scanner {
	return &Scanner{
		PortRange: [2]int{1, 65535},
		Protocols: []string{"tcp", "udp"},
	}
}

// Scan returns all currently open ports within the configured range.
func (s *Scanner) Scan() ([]PortState, error) {
	var open []PortState
	for _, proto := range s.Protocols {
		for port := s.PortRange[0]; port <= s.PortRange[1]; port++ {
			addr := net.JoinHostPort("127.0.0.1", strconv.Itoa(port))
			conn, err := net.Dial(proto, addr)
			if err != nil {
				continue
			}
			conn.Close()
			parts := strings.Split(conn.LocalAddr().String(), ":")
			open = append(open, PortState{
				Protocol: proto,
				Port:     port,
				Address:  "127.0.0.1",
			})
		}
	}
	return open, nil
}
