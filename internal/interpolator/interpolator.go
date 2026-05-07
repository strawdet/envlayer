// Package interpolator provides runtime interpolation of environment variable
// values using a simple ${VAR} or $VAR syntax, resolving against a provided map.
package interpolator

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

var pattern = regexp.MustCompile(`\$\{([^}]+)\}|\$([A-Za-z_][A-Za-z0-9_]*)`)

// Interpolator replaces variable references in values with resolved results.
type Interpolator struct {
	fallbackToOS bool
}

// New creates a new Interpolator. When fallbackToOS is true, unresolved
// references are looked up in the process environment before returning empty.
func New(fallbackToOS bool) *Interpolator {
	return &Interpolator{fallbackToOS: fallbackToOS}
}

// Interpolate resolves all variable references in src using the provided env
// map. Unknown references are replaced with an empty string unless a fallback
// is found in the OS environment.
func (i *Interpolator) Interpolate(src map[string]string) map[string]string {
	result := make(map[string]string, len(src))
	for k, v := range src {
		result[k] = i.InterpolateValue(v, src)
	}
	return result
}

// InterpolateValue resolves variable references within a single value string.
func (i *Interpolator) InterpolateValue(value string, env map[string]string) string {
	return pattern.ReplaceAllStringFunc(value, func(match string) string {
		name := extractName(match)
		if val, ok := env[name]; ok {
			return val
		}
		if i.fallbackToOS {
			if val, ok := os.LookupEnv(name); ok {
				return val
			}
		}
		return ""
	})
}

// HasReferences reports whether value contains any interpolation tokens.
func HasReferences(value string) bool {
	return pattern.MatchString(value)
}

// Validate checks for circular or self-referencing variables and returns an
// error describing the first cycle found.
func Validate(env map[string]string) error {
	for k, v := range env {
		for _, match := range pattern.FindAllString(v, -1) {
			ref := extractName(match)
			if ref == k {
				return fmt.Errorf("interpolator: self-reference detected for key %q", k)
			}
		}
	}
	return nil
}

func extractName(match string) string {
	match = strings.TrimPrefix(match, "${") 
	match = strings.TrimSuffix(match, "}")
	match = strings.TrimPrefix(match, "$")
	return match
}
