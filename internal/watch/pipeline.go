package watch

import "context"

// Stage is a function that processes a value and returns a transformed value or error.
type Stage[T any] func(ctx context.Context, in T) (T, error)

// Pipeline chains multiple stages sequentially, passing output of each into the next.
type Pipeline[T any] struct {
	stages []Stage[T]
}

// NewPipeline creates a Pipeline with the given stages.
func NewPipeline[T any](stages ...Stage[T]) *Pipeline[T] {
	return &Pipeline[T]{stages: stages}
}

// Run executes all stages in order. If any stage returns an error, execution stops.
func (p *Pipeline[T]) Run(ctx context.Context, in T) (T, error) {
	var err error
	for _, stage := range p.stages {
		if ctx.Err() != nil {
			return in, ctx.Err()
		}
		in, err = stage(ctx, in)
		if err != nil {
			return in, err
		}
	}
	return in, nil
}

// Len returns the number of stages in the pipeline.
func (p *Pipeline[T]) Len() int {
	return len(p.stages)
}
