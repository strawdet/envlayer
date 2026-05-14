// Package grouper provides functionality to group environment variables
// by namespace prefix, enabling structured access to related keys.
package grouper

import (
	"sort"
	"strings"
)

// Grouper partitions a flat env map into named groups based on key prefixes.
type Grouper struct {
	separator string
}

// Option is a functional option for Grouper.
type Option func(*Grouper)

// WithSeparator sets the separator used to split prefix from key name.
func WithSeparator(sep string) Option {
	return func(g *Grouper) {
		g.separator = sep
	}
}

// New creates a new Grouper with the given options.
// Default separator is "_".
func New(opts ...Option) *Grouper {
	g := &Grouper{separator: "_"}
	for _, o := range opts {
		o(g)
	}
	return g
}

// Group partitions env into a map of prefix -> (suffix -> value).
// Keys without a separator are placed under the empty-string group.
func (g *Grouper) Group(env map[string]string) map[string]map[string]string {
	result := make(map[string]map[string]string)
	for k, v := range env {
		prefix, suffix := g.split(k)
		if result[prefix] == nil {
			result[prefix] = make(map[string]string)
		}
		result[prefix][suffix] = v
	}
	return result
}

// Prefixes returns a sorted list of distinct prefixes present in env.
func (g *Grouper) Prefixes(env map[string]string) []string {
	seen := make(map[string]struct{})
	for k := range env {
		prefix, _ := g.split(k)
		seen[prefix] = struct{}{}
	}
	out := make([]string, 0, len(seen))
	for p := range seen {
		out = append(out, p)
	}
	sort.Strings(out)
	return out
}

// Flatten converts a grouped map back to a flat env map, re-joining prefix
// and suffix with the configured separator.
func (g *Grouper) Flatten(groups map[string]map[string]string) map[string]string {
	out := make(map[string]string)
	for prefix, kv := range groups {
		for suffix, val := range kv {
			key := g.join(prefix, suffix)
			out[key] = val
		}
	}
	return out
}

func (g *Grouper) split(key string) (prefix, suffix string) {
	idx := strings.Index(key, g.separator)
	if idx < 0 {
		return "", key
	}
	return key[:idx], key[idx+len(g.separator):]
}

func (g *Grouper) join(prefix, suffix string) string {
	if prefix == "" {
		return suffix
	}
	return prefix + g.separator + suffix
}
