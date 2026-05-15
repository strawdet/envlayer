// Package caster provides type-casting utilities for environment variable maps.
// It converts raw string values into typed Go values with configurable defaults
// and optional strict mode that returns errors on conversion failures.
package caster

import (
	"fmt"
	"strconv"
	"strings"
)

// Caster holds the env map and options.
type Caster struct {
	env    map[string]string
	strict bool
}

// Option configures a Caster.
type Option func(*Caster)

// WithStrict causes Cast to return errors instead of using defaults on bad input.
func WithStrict() Option {
	return func(c *Caster) { c.strict = true }
}

// New creates a Caster for the given env map.
func New(env map[string]string, opts ...Option) *Caster {
	c := &Caster{env: env}
	for _, o := range opts {
		o(c)
	}
	return c
}

// String returns the string value for key, or def if absent.
func (c *Caster) String(key, def string) string {
	if v, ok := c.env[key]; ok {
		return v
	}
	return def
}

// Int returns the integer value for key, or def on absence/parse error.
// In strict mode a parse error returns an error.
func (c *Caster) Int(key string, def int) (int, error) {
	v, ok := c.env[key]
	if !ok {
		return def, nil
	}
	n, err := strconv.Atoi(strings.TrimSpace(v))
	if err != nil {
		if c.strict {
			return 0, fmt.Errorf("caster: key %q: cannot cast %q to int", key, v)
		}
		return def, nil
	}
	return n, nil
}

// Bool returns the boolean value for key, or def on absence/parse error.
func (c *Caster) Bool(key string, def bool) (bool, error) {
	v, ok := c.env[key]
	if !ok {
		return def, nil
	}
	b, err := strconv.ParseBool(strings.TrimSpace(v))
	if err != nil {
		if c.strict {
			return false, fmt.Errorf("caster: key %q: cannot cast %q to bool", key, v)
		}
		return def, nil
	}
	return b, nil
}

// Float64 returns the float64 value for key, or def on absence/parse error.
func (c *Caster) Float64(key string, def float64) (float64, error) {
	v, ok := c.env[key]
	if !ok {
		return def, nil
	}
	f, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
	if err != nil {
		if c.strict {
			return 0, fmt.Errorf("caster: key %q: cannot cast %q to float64", key, v)
		}
		return def, nil
	}
	return f, nil
}

// Strings splits the value for key by sep and returns the parts.
// Returns def if the key is absent.
func (c *Caster) Strings(key, sep string, def []string) []string {
	v, ok := c.env[key]
	if !ok {
		return def
	}
	parts := strings.Split(v, sep)
	for i, p := range parts {
		parts[i] = strings.TrimSpace(p)
	}
	return parts
}
