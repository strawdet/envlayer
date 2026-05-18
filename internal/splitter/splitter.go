// Package splitter splits a flat env map into multiple named buckets
// based on key prefix, allowing consumers to work with isolated subsets.
package splitter

import "strings"

// Splitter partitions an env map into named buckets by key prefix.
type Splitter struct {
	separator string
	buckets   map[string]map[string]string
}

// Option is a functional option for Splitter.
type Option func(*Splitter)

// WithSeparator sets the separator used to split prefix from key name.
// Defaults to "_".
func WithSeparator(sep string) Option {
	return func(s *Splitter) {
		if sep != "" {
			s.separator = sep
		}
	}
}

// New creates a new Splitter with the provided options.
func New(opts ...Option) *Splitter {
	s := &Splitter{
		separator: "_",
		buckets:   make(map[string]map[string]string),
	}
	for _, o := range opts {
		o(s)
	}
	return s
}

// Split partitions env into buckets keyed by the first prefix segment.
// Keys without a separator are placed in the "_default" bucket.
func (s *Splitter) Split(env map[string]string) map[string]map[string]string {
	result := make(map[string]map[string]string)
	for k, v := range env {
		parts := strings.SplitN(k, s.separator, 2)
		if len(parts) == 2 && parts[0] != "" && parts[1] != "" {
			bucket := parts[0]
			if result[bucket] == nil {
				result[bucket] = make(map[string]string)
			}
			result[bucket][parts[1]] = v
		} else {
			if result["_default"] == nil {
				result["_default"] = make(map[string]string)
			}
			result["_default"][k] = v
		}
	}
	s.buckets = result
	return result
}

// Bucket returns a single named bucket from the last Split call.
// Returns nil if the bucket does not exist.
func (s *Splitter) Bucket(name string) map[string]string {
	b, ok := s.buckets[name]
	if !ok {
		return nil
	}
	out := make(map[string]string, len(b))
	for k, v := range b {
		out[k] = v
	}
	return out
}

// Buckets returns the names of all buckets from the last Split call.
func (s *Splitter) Buckets() []string {
	names := make([]string, 0, len(s.buckets))
	for name := range s.buckets {
		names = append(names, name)
	}
	return names
}
