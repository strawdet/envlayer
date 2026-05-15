// Package deduper provides utilities for detecting and removing duplicate
// key-value pairs across one or more environment maps.
package deduper

// Strategy controls how duplicates are resolved when merging layers.
type Strategy int

const (
	// KeepFirst retains the first occurrence of a duplicate key.
	KeepFirst Strategy = iota
	// KeepLast retains the last occurrence of a duplicate key.
	KeepLast
)

// Deduper removes duplicate keys from environment maps.
type Deduper struct {
	strategy Strategy
}

// Option is a functional option for Deduper.
type Option func(*Deduper)

// WithStrategy sets the deduplication strategy.
func WithStrategy(s Strategy) Option {
	return func(d *Deduper) {
		d.strategy = s
	}
}

// New creates a new Deduper with the given options.
// Default strategy is KeepLast.
func New(opts ...Option) *Deduper {
	d := &Deduper{strategy: KeepLast}
	for _, o := range opts {
		o(d)
	}
	return d
}

// Apply removes duplicate keys from a single env map.
// For a single map there are no cross-key duplicates; this is a no-op
// but is provided for interface consistency.
func (d *Deduper) Apply(env map[string]string) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = v
	}
	return out
}

// Merge combines multiple env layers, resolving duplicate keys according
// to the configured strategy.
func (d *Deduper) Merge(layers ...map[string]string) map[string]string {
	out := make(map[string]string)
	for _, layer := range layers {
		for k, v := range layer {
			_, exists := out[k]
			if !exists || d.strategy == KeepLast {
				out[k] = v
			}
		}
	}
	return out
}

// Duplicates returns all keys that appear in more than one layer.
func Duplicates(layers ...map[string]string) []string {
	seen := make(map[string]int)
	for _, layer := range layers {
		for k := range layer {
			seen[k]++
		}
	}
	var dups []string
	for k, count := range seen {
		if count > 1 {
			dups = append(dups, k)
		}
	}
	return dups
}

// DuplicateCount returns a map of keys to the number of layers they appear in,
// including only keys that appear in more than one layer.
func DuplicateCount(layers ...map[string]string) map[string]int {
	seen := make(map[string]int)
	for _, layer := range layers {
		for k := range layer {
			seen[k]++
		}
	}
	counts := make(map[string]int)
	for k, count := range seen {
		if count > 1 {
			counts[k] = count
		}
	}
	return counts
}
