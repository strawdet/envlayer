// Package scoper provides namespace-based scoping for environment variables,
// allowing keys to be grouped, filtered, and extracted by a namespace prefix.
package scoper

import (
	"fmt"
	"strings"
)

// Scoper manages namespaced environment variable sets.
type Scoper struct {
	separator string
}

// Option configures a Scoper.
type Option func(*Scoper)

// WithSeparator sets the delimiter between namespace and key (default "_").
func WithSeparator(sep string) Option {
	return func(s *Scoper) {
		s.separator = sep
	}
}

// New creates a new Scoper with the given options.
func New(opts ...Option) *Scoper {
	s := &Scoper{separator: "_"}
	for _, o := range opts {
		o(s)
	}
	return s
}

// Scope returns only the entries whose keys begin with namespace+separator,
// stripping the prefix from the returned map keys.
func (s *Scoper) Scope(env map[string]string, namespace string) map[string]string {
	prefix := strings.ToUpper(namespace) + s.separator
	out := make(map[string]string)
	for k, v := range env {
		if strings.HasPrefix(strings.ToUpper(k), prefix) {
			stripped := k[len(prefix):]
			out[stripped] = v
		}
	}
	return out
}

// Namespace prefixes all keys in env with namespace+separator and returns the
// new map, leaving the original unchanged.
func (s *Scoper) Namespace(env map[string]string, namespace string) map[string]string {
	prefix := strings.ToUpper(namespace) + s.separator
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[fmt.Sprintf("%s%s", prefix, k)] = v
	}
	return out
}

// Namespaces returns the distinct namespace prefixes found in env.
func (s *Scoper) Namespaces(env map[string]string) []string {
	seen := make(map[string]struct{})
	for k := range env {
		parts := strings.SplitN(k, s.separator, 2)
		if len(parts) == 2 && parts[0] != "" {
			seen[parts[0]] = struct{}{}
		}
	}
	result := make([]string, 0, len(seen))
	for ns := range seen {
		result = append(result, ns)
	}
	return result
}
