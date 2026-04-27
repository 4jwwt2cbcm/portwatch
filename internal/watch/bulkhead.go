package watch

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
)

// ErrBulkheadFull is returned when the bulkhead has no available slots.
var ErrBulkheadFull = errors.New("bulkhead: no available slots")

// DefaultBulkheadPolicy returns a BulkheadPolicy with sensible defaults.
func DefaultBulkheadPolicy() BulkheadPolicy {
	return BulkheadPolicy{
		MaxConcurrent: 8,
		QueueDepth:    16,
	}
}

// BulkheadPolicy configures a Bulkhead.
type BulkheadPolicy struct {
	MaxConcurrent int
	QueueDepth    int
}

// Bulkhead limits the number of concurrent executions and optionally queues
// overflow up to a fixed depth, shedding load beyond that.
type Bulkhead struct {
	policy  BulkheadPolicy
	sem     chan struct{}
	queue   chan struct{}
	active  atomic.Int64
	queued  atomic.Int64
	shed    atomic.Int64
	mu      sync.Mutex
}

// NewBulkhead creates a Bulkhead with the given policy.
func NewBulkhead(p BulkheadPolicy) *Bulkhead {
	if p.MaxConcurrent <= 0 {
		p.MaxConcurrent = DefaultBulkheadPolicy().MaxConcurrent
	}
	if p.QueueDepth < 0 {
		p.QueueDepth = 0
	}
	return &Bulkhead{
		policy: p,
		sem:    make(chan struct{}, p.MaxConcurrent),
		queue:  make(chan struct{}, p.QueueDepth),
	}
}

// Do runs fn if a slot is available, queuing if capacity allows, or returns
// ErrBulkheadFull when both the concurrency limit and queue are exhausted.
func (b *Bulkhead) Do(ctx context.Context, fn func() error) error {
	select {
	case b.sem <- struct{}{}:
		b.active.Add(1)
		defer func() {
			<-b.sem
			b.active.Add(-1)
		}()
		return fn()
	default:
	}

	select {
	case b.queue <- struct{}{}:
		b.queued.Add(1)
	default:
		b.shed.Add(1)
		return ErrBulkheadFull
	}

	select {
	case <-ctx.Done():
		<-b.queue
		b.queued.Add(-1)
		return ctx.Err()
	case b.sem <- struct{}{}:
		<-b.queue
		b.queued.Add(-1)
		b.active.Add(1)
		defer func() {
			<-b.sem
			b.active.Add(-1)
		}()
		return fn()
	}
}

// Active returns the number of concurrently executing calls.
func (b *Bulkhead) Active() int { return int(b.active.Load()) }

// Queued returns the number of calls waiting in the queue.
func (b *Bulkhead) Queued() int { return int(b.queued.Load()) }

// Shed returns the total number of calls rejected due to a full bulkhead.
func (b *Bulkhead) Shed() int { return int(b.shed.Load()) }

// Reset clears the shed counter. Active and queued state is not affected.
func (b *Bulkhead) Reset() { b.shed.Store(0) }
