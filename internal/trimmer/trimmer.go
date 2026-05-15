// Package trimmer provides utilities for trimming whitespace and
// normalizing environment variable values across a resolved env map.
package trimmer

import (
	"strings"
)

// Option configures a Trimmer.
type Option func(*Trimmer)

// Trimmer trims and normalizes env values.
type Trimmer struct {
	trimKeys    bool
	trimValues  bool
	collapseWS  bool
	sensitiveKeys map[string]struct{}
}

// WithTrimKeys enables trimming of leading/trailing whitespace from keys.
func WithTrimKeys() Option {
	return func(t *Trimmer) { t.trimKeys = true }
}

// WithTrimValues enables trimming of leading/trailing whitespace from values.
func WithTrimValues() Option {
	return func(t *Trimmer) { t.trimValues = true }
}

// WithCollapseWhitespace collapses internal runs of whitespace in values to a single space.
func WithCollapseWhitespace() Option {
	return func(t *Trimmer) { t.collapseWS = true }
}

// WithSkipKeys marks specific keys whose values should not be modified.
func WithSkipKeys(keys ...string) Option {
	return func(t *Trimmer) {
		for _, k := range keys {
			t.sensitiveKeys[k] = struct{}{}
		}
	}
}

// New creates a new Trimmer with the given options.
func New(opts ...Option) *Trimmer {
	t := &Trimmer{
		trimValues:    true,
		sensitiveKeys: make(map[string]struct{}),
	}
	for _, o := range opts {
		o(t)
	}
	return t
}

// Apply returns a new map with trimming rules applied.
func (t *Trimmer) Apply(env map[string]string) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		newKey := k
		if t.trimKeys {
			newKey = strings.TrimSpace(k)
		}
		if _, skip := t.sensitiveKeys[newKey]; skip {
			out[newKey] = v
			continue
		}
		newVal := v
		if t.trimValues {
			newVal = strings.TrimSpace(newVal)
		}
		if t.collapseWS {
			newVal = collapseWhitespace(newVal)
		}
		out[newKey] = newVal
	}
	return out
}

// collapseWhitespace replaces runs of whitespace with a single space.
func collapseWhitespace(s string) string {
	var b strings.Builder
	inSpace := false
	for _, r := range s {
		if r == ' ' || r == '\t' {
			if !inSpace {
				b.WriteRune(' ')
				inSpace = true
			}
		} else {
			b.WriteRune(r)
			inSpace = false
		}
	}
	return b.String()
}
