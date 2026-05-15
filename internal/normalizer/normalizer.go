// Package normalizer provides key and value normalization for environment maps.
// It supports trimming, case normalization, and character replacement to
// produce consistent, well-formed environment variable entries.
package normalizer

import (
	"strings"
)

// Option configures the Normalizer.
type Option func(*Normalizer)

// Normalizer applies normalization rules to an environment map.
type Normalizer struct {
	upperKeys      bool
	lowerKeys      bool
	replaceHyphens bool
	trimValues     bool
}

// WithUpperKeys converts all keys to UPPER_CASE.
func WithUpperKeys() Option {
	return func(n *Normalizer) { n.upperKeys = true }
}

// WithLowerKeys converts all keys to lower_case.
func WithLowerKeys() Option {
	return func(n *Normalizer) { n.lowerKeys = true }
}

// WithReplaceHyphens replaces hyphens in keys with underscores.
func WithReplaceHyphens() Option {
	return func(n *Normalizer) { n.replaceHyphens = true }
}

// WithTrimValues trims leading and trailing whitespace from all values.
func WithTrimValues() Option {
	return func(n *Normalizer) { n.trimValues = true }
}

// New creates a Normalizer with the given options.
func New(opts ...Option) *Normalizer {
	n := &Normalizer{}
	for _, o := range opts {
		o(n)
	}
	return n
}

// Apply returns a new map with normalization rules applied.
// If both WithUpperKeys and WithLowerKeys are set, upper takes precedence.
func (n *Normalizer) Apply(env map[string]string) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		nk := n.normalizeKey(k)
		nv := n.normalizeValue(v)
		out[nk] = nv
	}
	return out
}

func (n *Normalizer) normalizeKey(k string) string {
	if n.replaceHyphens {
		k = strings.ReplaceAll(k, "-", "_")
	}
	switch {
	case n.upperKeys:
		return strings.ToUpper(k)
	case n.lowerKeys:
		return strings.ToLower(k)
	}
	return k
}

func (n *Normalizer) normalizeValue(v string) string {
	if n.trimValues {
		return strings.TrimSpace(v)
	}
	return v
}
