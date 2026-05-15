// Package cloner provides utilities for deep-copying environment variable maps
// with optional key transformation and filtering during the clone operation.
package cloner

import "strings"

// Option configures the Cloner.
type Option func(*Cloner)

// Cloner deep-copies env maps with optional transformations.
type Cloner struct {
	prefixFilter string
	stripPrefix  bool
	omitEmpty    bool
}

// WithPrefixFilter restricts cloning to keys that begin with the given prefix.
func WithPrefixFilter(prefix string) Option {
	return func(c *Cloner) { c.prefixFilter = prefix }
}

// WithStripPrefix removes the prefix from keys during cloning (requires WithPrefixFilter).
func WithStripPrefix() Option {
	return func(c *Cloner) { c.stripPrefix = true }
}

// WithOmitEmpty skips keys whose values are empty strings.
func WithOmitEmpty() Option {
	return func(c *Cloner) { c.omitEmpty = true }
}

// New creates a Cloner with the given options.
func New(opts ...Option) *Cloner {
	c := &Cloner{}
	for _, o := range opts {
		o(c)
	}
	return c
}

// Clone returns a deep copy of env, applying any configured transformations.
func (c *Cloner) Clone(env map[string]string) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		if c.omitEmpty && v == "" {
			continue
		}
		if c.prefixFilter != "" && !strings.HasPrefix(k, c.prefixFilter) {
			continue
		}
		key := k
		if c.stripPrefix && c.prefixFilter != "" {
			key = strings.TrimPrefix(k, c.prefixFilter)
		}
		out[key] = v
	}
	return out
}

// Merge clones src into dst, returning a new map. Keys in src override dst.
func (c *Cloner) Merge(dst, src map[string]string) map[string]string {
	base := c.Clone(dst)
	for k, v := range c.Clone(src) {
		base[k] = v
	}
	return base
}
