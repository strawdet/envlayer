// Package expander provides variable interpolation for environment values.
// It resolves references like ${VAR} or $VAR within env map values.
package expander

import (
	"os"
	"strings"
)

// Expander resolves variable references within an env map.
type Expander struct {
	// fallbackToOS controls whether unresolved vars fall back to OS env.
	fallbackToOS bool
}

// New creates a new Expander.
// If fallbackToOS is true, variables not found in the map are looked up
// in the process environment.
func New(fallbackToOS bool) *Expander {
	return &Expander{fallbackToOS: fallbackToOS}
}

// Expand resolves all variable references in the values of the provided map.
// It returns a new map with interpolated values. The original map is not modified.
// References use the ${VAR} or $VAR syntax.
func (e *Expander) Expand(env map[string]string) map[string]string {
	result := make(map[string]string, len(env))
	for k, v := range env {
		result[k] = e.expandValue(v, env)
	}
	return result
}

// expandValue interpolates variable references within a single string value.
func (e *Expander) expandValue(value string, env map[string]string) string {
	return os.Expand(value, func(key string) string {
		if val, ok := env[key]; ok {
			return val
		}
		if e.fallbackToOS {
			return os.Getenv(key)
		}
		return ""
	})
}

// ExpandValue is a convenience function that expands a single value string
// using the provided env map as the source of variable values.
// Variables not found in env resolve to empty string.
func ExpandValue(value string, env map[string]string) string {
	e := New(false)
	return e.expandValue(value, env)
}

// HasReferences reports whether the given value contains any variable references.
func HasReferences(value string) bool {
	return strings.Contains(value, "$")
}
