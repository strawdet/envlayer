// Package filter provides key-based filtering of environment variable maps.
// It supports inclusion lists, exclusion lists, and glob-style prefix matching.
package filter

import "strings"

// Option configures a Filter.
type Option func(*Filter)

// Filter selectively includes or excludes keys from an env map.
type Filter struct {
	include []string
	exclude []string
	prefixes []string
}

// WithInclude restricts output to only the specified keys.
func WithInclude(keys ...string) Option {
	return func(f *Filter) {
		f.include = append(f.include, keys...)
	}
}

// WithExclude removes the specified keys from output.
func WithExclude(keys ...string) Option {
	return func(f *Filter) {
		f.exclude = append(f.exclude, keys...)
	}
}

// WithPrefixInclude includes only keys that start with one of the given prefixes.
func WithPrefixInclude(prefixes ...string) Option {
	return func(f *Filter) {
		f.prefixes = append(f.prefixes, prefixes...)
	}
}

// New creates a Filter with the given options.
func New(opts ...Option) *Filter {
	f := &Filter{}
	for _, o := range opts {
		o(f)
	}
	return f
}

// Apply returns a new map containing only the keys that pass all filter rules.
// Precedence: include list > prefix include > exclude list.
func (f *Filter) Apply(env map[string]string) map[string]string {
	out := make(map[string]string, len(env))

	for k, v := range env {
		if len(f.include) > 0 && !contains(f.include, k) {
			continue
		}
		if len(f.prefixes) > 0 && !hasPrefix(f.prefixes, k) {
			continue
		}
		if contains(f.exclude, k) {
			continue
		}
		out[k] = v
	}
	return out
}

func contains(list []string, key string) bool {
	for _, item := range list {
		if item == key {
			return true
		}
	}
	return false
}

func hasPrefix(prefixes []string, key string) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(key, p) {
			return true
		}
	}
	return false
}
