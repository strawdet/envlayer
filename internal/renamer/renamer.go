// Package renamer provides key renaming functionality for environment variable maps.
// It supports explicit key mappings and regex-based pattern renaming.
package renamer

import (
	"fmt"
	"regexp"
)

// Option configures the Renamer.
type Option func(*Renamer)

// Renamer renames keys in an environment map.
type Renamer struct {
	mappings map[string]string
	patterns []patternRule
}

type patternRule struct {
	re          *regexp.Regexp
	replacement string
}

// WithMapping adds an explicit old->new key rename.
func WithMapping(oldKey, newKey string) Option {
	return func(r *Renamer) {
		r.mappings[oldKey] = newKey
	}
}

// WithPattern adds a regex pattern rename rule.
// The replacement string may use $1, $2, etc. for capture groups.
func WithPattern(pattern, replacement string) (Option, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("renamer: invalid pattern %q: %w", pattern, err)
	}
	return func(r *Renamer) {
		r.patterns = append(r.patterns, patternRule{re: re, replacement: replacement})
	}, nil
}

// New creates a new Renamer with the given options.
func New(opts ...Option) *Renamer {
	r := &Renamer{
		mappings: make(map[string]string),
	}
	for _, o := range opts {
		o(r)
	}
	return r
}

// Apply returns a new map with keys renamed according to configured rules.
// Explicit mappings take precedence over pattern rules.
// If a rename would cause a collision, the original key is kept.
func (r *Renamer) Apply(env map[string]string) map[string]string {
	result := make(map[string]string, len(env))
	for k, v := range env {
		newKey := r.resolveKey(k)
		if _, exists := result[newKey]; exists {
			// collision: keep original key
			result[k] = v
		} else {
			result[newKey] = v
		}
	}
	return result
}

// resolveKey returns the renamed key for k, or k itself if no rule matches.
func (r *Renamer) resolveKey(k string) string {
	if newKey, ok := r.mappings[k]; ok {
		return newKey
	}
	for _, p := range r.patterns {
		if p.re.MatchString(k) {
			return p.re.ReplaceAllString(k, p.replacement)
		}
	}
	return k
}
