// Package cache provides a simple in-memory cache for resolved environment
// variable maps, keyed by context string. This avoids re-reading and re-merging
// layer files on repeated lookups within the same process lifetime.
package cache

import (
	"sync"
	"time"
)

// Entry holds a cached environment map and its expiry time.
type Entry struct {
	Env       map[string]string
	CachedAt  time.Time
	ExpiresAt time.Time
}

// IsExpired reports whether the cache entry has passed its TTL.
func (e Entry) IsExpired() bool {
	return time.Now().After(e.ExpiresAt)
}

// Cache is a thread-safe in-memory store for resolved env maps.
type Cache struct {
	mu      sync.RWMutex
	entries map[string]Entry
	ttl     time.Duration
}

// New creates a Cache with the given TTL for each entry.
func New(ttl time.Duration) *Cache {
	return &Cache{
		entries: make(map[string]Entry),
		ttl:     ttl,
	}
}

// Get returns the env map for the given context key, if present and not expired.
func (c *Cache) Get(ctx string) (map[string]string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	e, ok := c.entries[ctx]
	if !ok || e.IsExpired() {
		return nil, false
	}
	copy := make(map[string]string, len(e.Env))
	for k, v := range e.Env {
		copy[k] = v
	}
	return copy, true
}

// Set stores an env map under the given context key with the configured TTL.
func (c *Cache) Set(ctx string, env map[string]string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now()
	copy := make(map[string]string, len(env))
	for k, v := range env {
		copy[k] = v
	}
	c.entries[ctx] = Entry{
		Env:       copy,
		CachedAt:  now,
		ExpiresAt: now.Add(c.ttl),
	}
}

// Invalidate removes the entry for the given context key.
func (c *Cache) Invalidate(ctx string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, ctx)
}

// Flush removes all entries from the cache.
func (c *Cache) Flush() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries = make(map[string]Entry)
}

// Len returns the number of entries currently in the cache (including expired).
func (c *Cache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.entries)
}
