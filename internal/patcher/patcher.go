// Package patcher provides utilities for applying partial updates to
// an existing environment map — adding, updating, or removing keys
// without replacing the entire set.
package patcher

import "fmt"

// Op represents the type of patch operation.
type Op int

const (
	OpSet    Op = iota // Add or update a key.
	OpDelete           // Remove a key.
)

// Patch describes a single change to apply to an env map.
type Patch struct {
	Op    Op
	Key   string
	Value string // ignored for OpDelete
}

// Option configures a Patcher.
type Option func(*Patcher)

// WithIgnoreUnknownDeletes silences errors when deleting a key that
// does not exist in the base env.
func WithIgnoreUnknownDeletes() Option {
	return func(p *Patcher) { p.ignoreUnknownDeletes = true }
}

// Patcher applies a slice of Patch operations to an env map.
type Patcher struct {
	ignoreUnknownDeletes bool
}

// New creates a Patcher with the supplied options.
func New(opts ...Option) *Patcher {
	p := &Patcher{}
	for _, o := range opts {
		o(p)
	}
	return p
}

// Apply returns a new map that is the result of applying patches to base.
// The original base map is never mutated.
func (p *Patcher) Apply(base map[string]string, patches []Patch) (map[string]string, error) {
	out := make(map[string]string, len(base))
	for k, v := range base {
		out[k] = v
	}

	for _, patch := range patches {
		if patch.Key == "" {
			return nil, fmt.Errorf("patcher: patch key must not be empty")
		}
		switch patch.Op {
		case OpSet:
			out[patch.Key] = patch.Value
		case OpDelete:
			if _, exists := out[patch.Key]; !exists && !p.ignoreUnknownDeletes {
				return nil, fmt.Errorf("patcher: delete of unknown key %q", patch.Key)
			}
			delete(out, patch.Key)
		default:
			return nil, fmt.Errorf("patcher: unknown op %d for key %q", patch.Op, patch.Key)
		}
	}
	return out, nil
}

// Diff returns the patches needed to transform src into dst.
// Useful for computing the delta between two env snapshots.
func Diff(src, dst map[string]string) []Patch {
	var patches []Patch
	for k, v := range dst {
		if existing, ok := src[k]; !ok || existing != v {
			patches = append(patches, Patch{Op: OpSet, Key: k, Value: v})
		}
	}
	for k := range src {
		if _, ok := dst[k]; !ok {
			patches = append(patches, Patch{Op: OpDelete, Key: k})
		}
	}
	return patches
}
