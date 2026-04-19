package watch

import (
	"context"
	"fmt"
	"net"
	"time"
)

// ProbePolicy configures the TCP probe behaviour.
type ProbePolicy struct {
	Timeout time.Duration
	Retries int
}

// DefaultProbePolicy returns sensible probe defaults.
func DefaultProbePolicy() ProbePolicy {
	return ProbePolicy{
		Timeout: 2 * time.Second,
		Retries: 2,
	}
}

// Probe checks whether a TCP endpoint is reachable.
type Probe struct {
	policy ProbePolicy
}

// NewProbe creates a Probe with the given policy.
func NewProbe(p ProbePolicy) *Probe {
	if p.Timeout <= 0 {
		p.Timeout = DefaultProbePolicy().Timeout
	}
	if p.Retries < 0 {
		p.Retries = 0
	}
	return &Probe{policy: p}
}

// Check attempts to connect to host:port, retrying up to policy.Retries times.
// Returns nil if the port is reachable, otherwise the last error.
func (pr *Probe) Check(ctx context.Context, host string, port int) error {
	addr := fmt.Sprintf("%s:%d", host, port)
	var lastErr error
	for i := 0; i <= pr.policy.Retries; i++ {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		d := net.Dialer{Timeout: pr.policy.Timeout}
		conn, err := d.DialContext(ctx, "tcp", addr)
		if err == nil {
			conn.Close()
			return nil
		}
		lastErr = err
	}
	return lastErr
}

// Reachable is a convenience wrapper that returns a bool.
func (pr *Probe) Reachable(ctx context.Context, host string, port int) bool {
	return pr.Check(ctx, host, port) == nil
}
