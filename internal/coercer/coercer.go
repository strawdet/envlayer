// Package coercer provides type coercion for environment variable values.
// It converts string values from env maps into typed Go values such as
// int, float64, bool, and duration, with optional default fallbacks.
package coercer

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Coercer converts string environment values to typed Go values.
type Coercer struct {
	env map[string]string
}

// New creates a new Coercer backed by the provided env map.
func New(env map[string]string) *Coercer {
	copy := make(map[string]string, len(env))
	for k, v := range env {
		copy[k] = v
	}
	return &Coercer{env: copy}
}

// String returns the raw string value for key, or def if not found.
func (c *Coercer) String(key, def string) string {
	if v, ok := c.env[key]; ok {
		return v
	}
	return def
}

// Int returns the integer value for key, or def if not found or unparseable.
func (c *Coercer) Int(key string, def int) (int, error) {
	v, ok := c.env[key]
	if !ok {
		return def, nil
	}
	n, err := strconv.Atoi(strings.TrimSpace(v))
	if err != nil {
		return def, fmt.Errorf("coercer: key %q value %q is not a valid int", key, v)
	}
	return n, nil
}

// Float64 returns the float64 value for key, or def if not found or unparseable.
func (c *Coercer) Float64(key string, def float64) (float64, error) {
	v, ok := c.env[key]
	if !ok {
		return def, nil
	}
	f, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
	if err != nil {
		return def, fmt.Errorf("coercer: key %q value %q is not a valid float64", key, v)
	}
	return f, nil
}

// Bool returns the boolean value for key, or def if not found or unparseable.
// Accepts: 1, t, true, yes, on (case-insensitive) as true.
func (c *Coercer) Bool(key string, def bool) (bool, error) {
	v, ok := c.env[key]
	if !ok {
		return def, nil
	}
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "1", "t", "true", "yes", "on":
		return true, nil
	case "0", "f", "false", "no", "off":
		return false, nil
	}
	return def, fmt.Errorf("coercer: key %q value %q is not a valid bool", key, v)
}

// Duration returns the time.Duration value for key, or def if not found or unparseable.
func (c *Coercer) Duration(key string, def time.Duration) (time.Duration, error) {
	v, ok := c.env[key]
	if !ok {
		return def, nil
	}
	d, err := time.ParseDuration(strings.TrimSpace(v))
	if err != nil {
		return def, fmt.Errorf("coercer: key %q value %q is not a valid duration", key, v)
	}
	return d, nil
}
