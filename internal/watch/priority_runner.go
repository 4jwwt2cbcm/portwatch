package watch

import (
	"context"
	"time"
)

// PriorityTask is a unit of work with an associated priority.
type PriorityTask struct {
	Name     string
	Priority int
	Fn       func(ctx context.Context) error
}

// PriorityRunner drains a PriorityQueue of tasks in priority order.
type PriorityRunner struct {
	queue    *PriorityQueue[PriorityTask]
	pollInterval time.Duration
	logger   Logger
}

// NewPriorityRunner creates a PriorityRunner backed by the given queue.
// pollInterval controls how often the queue is checked when empty.
func NewPriorityRunner(q *PriorityQueue[PriorityTask], pollInterval time.Duration, logger Logger) *PriorityRunner {
	if logger == nil {
		logger = NewLogger(nil)
	}
	if pollInterval <= 0 {
		pollInterval = 100 * time.Millisecond
	}
	return &PriorityRunner{queue: q, pollInterval: pollInterval, logger: logger}
}

// Run processes tasks from the queue until the context is cancelled.
func (r *PriorityRunner) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		task, _, ok := r.queue.Pop()
		if !ok {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(r.pollInterval):
				continue
			}
		}

		r.logger.Info("running task: " + task.Name)
		if err := task.Fn(ctx); err != nil {
			r.logger.Error("task failed: "+task.Name, err)
		}
	}
}
