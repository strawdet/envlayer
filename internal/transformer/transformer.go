// Package transformer provides key/value transformation utilities for
// environment variable maps, such as prefix filtering, key renaming,
// and value case normalization.
package transformer

import (
	"strings"
)

// Option configures a Transformer.
type Option func(*Transformer)

// Transformer applies a chain of transformations to an env map.
type Transformer struct {
	prefixFilter string
	stripPrefix  bool
	keyCase      string // "upper", "lower", or ""
	valueCase    string // "upper", "lower", or ""
}

// WithPrefixFilter retains only keys that start with the given prefix.
func WithPrefixFilter(prefix string) Option {
	return func(t *Transformer) { t.prefixFilter = prefix }
}

// WithStripPrefix removes the prefix from matching keys after filtering.
func WithStripPrefix(strip bool) Option {
	return func(t *Transformer) { t.stripPrefix = strip }
}

// WithKeyCase normalises all keys to "upper" or "lower" case.
func WithKeyCase(c string) Option {
	return func(t *Transformer) { t.keyCase = c }
}

// WithValueCase normalises all values to "upper" or "lower" case.
func WithValueCase(c string) Option {
	return func(t *Transformer) { t.valueCase = c }
}

// New creates a Transformer with the supplied options.
func New(opts ...Option) *Transformer {
	t := &Transformer{}
	for _, o := range opts {
		o(t)
	}
	return t
}

// Apply runs all configured transformations on a copy of env and returns it.
func (t *Transformer) Apply(env map[string]string) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		if t.prefixFilter != "" && !strings.HasPrefix(k, t.prefixFilter) {
			continue
		}
		if t.stripPrefix && t.prefixFilter != "" {
			k = strings.TrimPrefix(k, t.prefixFilter)
		}
		switch t.keyCase {
		case "upper":
			k = strings.ToUpper(k)
		case "lower":
			k = strings.ToLower(k)
		}
		switch t.valueCase {
		case "upper":
			v = strings.ToUpper(v)
		case "lower":
			v = strings.ToLower(v)
		}
		out[k] = v
	}
	return out
}

// Keys returns the sorted list of keys present after applying transformations.
func (t *Transformer) Keys(env map[string]string) []string {
	applied := t.Apply(env)
	keys := make([]string, 0, len(applied))
	for k := range applied {
		keys = append(keys, k)
	}
	return keys
}
