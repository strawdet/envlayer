// Package cache provides a lightweight, thread-safe, TTL-based in-memory cache
// for resolved environment variable maps produced by the envlayer resolver.
//
// # Overview
//
// When the resolver merges multiple .env layer files for a given runtime context,
// the result can be stored in the cache to avoid redundant file I/O on repeated
// lookups. Each entry is associated with a context key (e.g. "production") and
// expires after a configurable duration.
//
// # Cache Invalidation
//
// Entries are evicted lazily on access once their TTL has elapsed. No background
// goroutine is used, keeping the implementation simple and allocation-free at
// rest. Callers that require proactive eviction can call [Cache.Flush] to clear
// all entries immediately.
//
// # Usage
//
//	c := cache.New(30 * time.Second)
//
//	// Attempt to read from cache before resolving
//	if env, ok := c.Get(ctx); ok {
//		return env, nil
//	}
//
//	env, err := resolver.Resolve(ctx)
//	if err != nil {
//		return nil, err
//	}
//	c.Set(ctx, env)
//	return env, nil
package cache
