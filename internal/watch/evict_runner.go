package watch

import "context"

// EvictRunner wraps a function with an EvictCache so that repeated calls
// with the same key return a cached result until the TTL expires.
type EvictRunner[T any] struct {
	cache *EvictCache[T]
	fn    func(ctx context.Context, key string) (T, error)
}

// NewEvictRunner returns an EvictRunner backed by the given cache.
// If cache is nil, a default cache is created. If fn is nil, a no-op is used.
func NewEvictRunner[T any](cache *EvictCache[T], fn func(ctx context.Context, key string) (T, error)) *EvictRunner[T] {
	if cache == nil {
		cache = NewEvictCache[T](DefaultEvictPolicy())
	}
	if fn == nil {
		fn = func(_ context.Context, _ string) (T, error) {
			var zero T
			return zero, nil
		}
	}
	return &EvictRunner[T]{cache: cache, fn: fn}
}

// Run returns a cached result for key if one exists and has not expired.
// Otherwise it calls the underlying function, caches the result, and returns it.
func (r *EvictRunner[T]) Run(ctx context.Context, key string) (T, error) {
	if v, ok := r.cache.Get(key); ok {
		return v, nil
	}
	v, err := r.fn(ctx, key)
	if err != nil {
		return v, err
	}
	r.cache.Set(key, v)
	return v, nil
}
