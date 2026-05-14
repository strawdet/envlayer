// Package sorter provides utilities for sorting and ordering environment
// variable maps by key, value, or custom criteria.
package sorter

import (
	"cmp"
	"slices"
	"strings"
)

// Order defines the sort direction.
type Order int

const (
	Ascending  Order = iota
	Descending
)

// Option configures the Sorter.
type Option func(*Sorter)

// Sorter sorts environment variable maps.
type Sorter struct {
	order      Order
	byValue    bool
	caseInsens bool
}

// WithOrder sets the sort direction (default: Ascending).
func WithOrder(o Order) Option {
	return func(s *Sorter) { s.order = o }
}

// WithSortByValue sorts entries by value instead of key.
func WithSortByValue() Option {
	return func(s *Sorter) { s.byValue = true }
}

// WithCaseInsensitive enables case-insensitive comparison.
func WithCaseInsensitive() Option {
	return func(s *Sorter) { s.caseInsens = true }
}

// New creates a new Sorter with the given options.
func New(opts ...Option) *Sorter {
	s := &Sorter{order: Ascending}
	for _, o := range opts {
		o(s)
	}
	return s
}

// SortedKeys returns the keys of env sorted according to the Sorter's configuration.
func (s *Sorter) SortedKeys(env map[string]string) []string {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}

	norm := func(v string) string {
		if s.caseInsens {
			return strings.ToLower(v)
		}
		return v
	}

	slices.SortFunc(keys, func(a, b string) int {
		var va, vb string
		if s.byValue {
			va, vb = norm(env[a]), norm(env[b])
		} else {
			va, vb = norm(a), norm(b)
		}
		result := cmp.Compare(va, vb)
		if s.order == Descending {
			return -result
		}
		return result
	})
	return keys
}

// Sorted returns a slice of [2]string pairs {key, value} in sorted order.
func (s *Sorter) Sorted(env map[string]string) [][2]string {
	keys := s.SortedKeys(env)
	pairs := make([][2]string, len(keys))
	for i, k := range keys {
		pairs[i] = [2]string{k, env[k]}
	}
	return pairs
}
