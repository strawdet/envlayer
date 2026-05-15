// Package aliaser provides key aliasing for environment variable maps.
// It allows defining alternate names for existing keys, copying or moving
// values under new names without modifying the original entries.
package aliaser

import "fmt"

// Option configures an Aliaser.
type Option func(*Aliaser)

// Aliaser maps source keys to one or more alias keys.
type Aliaser struct {
	aliases  map[string][]string // source -> []alias
	keepOrig bool
}

// WithKeepOriginal controls whether the original key is retained after aliasing.
// Default is true.
func WithKeepOriginal(keep bool) Option {
	return func(a *Aliaser) {
		a.keepOrig = keep
	}
}

// WithAliases registers a mapping of source key to alias keys.
func WithAliases(aliases map[string][]string) Option {
	return func(a *Aliaser) {
		for src, targets := range aliases {
			a.aliases[src] = append(a.aliases[src], targets...)
		}
	}
}

// New creates a new Aliaser with the given options.
func New(opts ...Option) *Aliaser {
	a := &Aliaser{
		aliases:  make(map[string][]string),
		keepOrig: true,
	}
	for _, o := range opts {
		o(a)
	}
	return a
}

// Apply returns a new map with aliases applied to env.
// If a source key is not present in env the alias is silently skipped.
// Returns an error if an alias target key already exists in env.
func (a *Aliaser) Apply(env map[string]string) (map[string]string, error) {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = v
	}
	for src, targets := range a.aliases {
		val, ok := out[src]
		if !ok {
			continue
		}
		for _, t := range targets {
			if _, exists := out[t]; exists {
				return nil, fmt.Errorf("aliaser: target key %q already exists", t)
			}
			out[t] = val
		}
		if !a.keepOrig {
			delete(out, src)
		}
	}
	return out, nil
}

// Aliases returns a copy of the registered alias map.
func (a *Aliaser) Aliases() map[string][]string {
	copy := make(map[string][]string, len(a.aliases))
	for k, v := range a.aliases {
		tmp := make([]string, len(v))
		_ = copy_(tmp, v)
		copy[k] = tmp
	}
	return copy
}

func copy_(dst, src []string) int { return copy(dst, src) }
