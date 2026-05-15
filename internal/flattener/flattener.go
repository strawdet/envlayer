// Package flattener provides utilities for flattening nested key structures
// in environment variable maps using a configurable separator.
package flattener

import (
	"strings"
)

// Option is a functional option for configuring a Flattener.
type Option func(*Flattener)

// Flattener collapses nested key prefixes into a single-level map.
type Flattener struct {
	separator string
	prefix    string
	lowerKeys bool
}

// WithSeparator sets the separator used to join nested key segments.
func WithSeparator(sep string) Option {
	return func(f *Flattener) { f.separator = sep }
}

// WithPrefix restricts flattening to keys that begin with the given prefix.
func WithPrefix(prefix string) Option {
	return func(f *Flattener) { f.prefix = prefix }
}

// WithLowerKeys normalises all output keys to lower-case.
func WithLowerKeys() Option {
	return func(f *Flattener) { f.lowerKeys = true }
}

// New creates a Flattener with the supplied options.
func New(opts ...Option) *Flattener {
	f := &Flattener{separator: "_"}
	for _, o := range opts {
		o(f)
	}
	return f
}

// Flatten takes a map whose keys may contain the separator and returns a new
// map where each key is the last segment after splitting on the separator.
// Collisions are resolved by keeping the last value encountered.
func (f *Flattener) Flatten(env map[string]string) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		if f.prefix != "" && !strings.HasPrefix(k, f.prefix) {
			continue
		}
		parts := strings.Split(k, f.separator)
		short := parts[len(parts)-1]
		if f.lowerKeys {
			short = strings.ToLower(short)
		}
		out[short] = v
	}
	return out
}

// Segments returns all unique leading segments (i.e. everything before the
// first separator) present in the provided map.
func (f *Flattener) Segments(env map[string]string) []string {
	seen := map[string]struct{}{}
	for k := range env {
		parts := strings.SplitN(k, f.separator, 2)
		if len(parts) > 1 {
			seen[parts[0]] = struct{}{}
		}
	}
	out := make([]string, 0, len(seen))
	for s := range seen {
		out = append(out, s)
	}
	return out
}
