package watch

import (
	"context"
	"errors"
	"time"
)

// ErrLeaseNotHeld is returned when a function is run without holding the lease.
var ErrLeaseNotHeld = errors.New("lease not held")

// LeaseRunner wraps a Lease and runs a function only while the lease is held,
// automatically renewing the lease when ShouldRenew returns true.
type LeaseRunner struct {
	lease  *Lease
	ticker func(time.Duration) *time.Ticker
}

// NewLeaseRunner creates a LeaseRunner backed by the given Lease.
// If lease is nil, a default lease is created.
func NewLeaseRunner(lease *Lease) *LeaseRunner {
	if lease == nil {
		lease = NewLease(DefaultLeasePolicy())
	}
	return &LeaseRunner{
		lease:  lease,
		ticker: time.NewTicker,
	}
}

// Run acquires the lease and calls fn. Returns ErrLeaseNotHeld if acquisition fails.
// Renews the lease in the background at the policy interval until fn returns or ctx is done.
func (r *LeaseRunner) Run(ctx context.Context, fn func(ctx context.Context) error) error {
	if !r.lease.Acquire() {
		return ErrLeaseNotHeld
	}
	defer r.lease.Release()

	renewInterval := time.Duration(float64(r.lease.policy.TTL) * (1 - r.lease.policy.RenewAt))
	if renewInterval <= 0 {
		renewInterval = r.lease.policy.TTL / 4
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		t := r.ticker(renewInterval)
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				if r.lease.ShouldRenew() {
					r.lease.Renew()
				}
			}
		}
	}()

	return fn(ctx)
}

// Held reports whether the underlying lease is currently held.
func (r *LeaseRunner) Held() bool {
	return r.lease.Held()
}
