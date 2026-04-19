package watch

import "time"

// EnvelopeMetadata holds routing and tracing metadata for a wrapped value.
type EnvelopeMetadata struct {
	ID        string
	Source    string
	Timestamp time.Time
	Attempt   int
}

// Envelope wraps a value with metadata for tracing and routing through pipelines.
type Envelope[T any] struct {
	Meta    EnvelopeMetadata
	Payload T
	Err     error
}

// NewEnvelope creates an Envelope with the given source and payload.
func NewEnvelope[T any](id, source string, payload T) Envelope[T] {
	return Envelope[T]{
		Meta: EnvelopeMetadata{
			ID:        id,
			Source:    source,
			Timestamp: time.Now(),
			Attempt:   1,
		},
		Payload: payload,
	}
}

// WithError returns a copy of the Envelope with the given error attached.
func (e Envelope[T]) WithError(err error) Envelope[T] {
	e.Err = err
	return e
}

// WithAttempt returns a copy of the Envelope with the attempt count incremented.
func (e Envelope[T]) WithAttempt(n int) Envelope[T] {
	e.Meta.Attempt = n
	return e
}

// OK returns true if the envelope carries no error.
func (e Envelope[T]) OK() bool {
	return e.Err == nil
}
