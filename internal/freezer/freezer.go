// Package freezer provides functionality to lock (freeze) a resolved
// environment map, preventing further modification and enabling
// change detection against a frozen baseline.
package freezer

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

// ErrFrozen is returned when an attempt is made to modify a frozen environment.
var ErrFrozen = errors.New("environment is frozen and cannot be modified")

// Freezer holds a frozen snapshot of an environment map.
type Freezer struct {
	env    map[string]string
	frozen bool
}

// New creates a new Freezer. The environment is not frozen until Freeze is called.
func New() *Freezer {
	return &Freezer{
		env: make(map[string]string),
	}
}

// Set stores a key-value pair. Returns ErrFrozen if the environment is frozen.
func (f *Freezer) Set(key, value string) error {
	if f.frozen {
		return fmt.Errorf("%w: cannot set %q", ErrFrozen, key)
	}
	f.env[key] = value
	return nil
}

// Load populates the freezer from an existing map. Returns ErrFrozen if already frozen.
func (f *Freezer) Load(env map[string]string) error {
	if f.frozen {
		return ErrFrozen
	}
	for k, v := range env {
		f.env[k] = v
	}
	return nil
}

// Freeze locks the environment, preventing further modifications.
func (f *Freezer) Freeze() {
	f.frozen = true
}

// IsFrozen reports whether the environment has been frozen.
func (f *Freezer) IsFrozen() bool {
	return f.frozen
}

// Get retrieves a value by key. Returns the value and whether it was found.
func (f *Freezer) Get(key string) (string, bool) {
	v, ok := f.env[key]
	return v, ok
}

// Snapshot returns a copy of the frozen environment. Safe to call at any time.
func (f *Freezer) Snapshot() map[string]string {
	out := make(map[string]string, len(f.env))
	for k, v := range f.env {
		out[k] = v
	}
	return out
}

// Diff compares the frozen environment against another map and returns
// a human-readable summary of added, removed, and changed keys.
func (f *Freezer) Diff(other map[string]string) []string {
	var lines []string

	for k, v := range other {
		if orig, ok := f.env[k]; !ok {
			lines = append(lines, fmt.Sprintf("+ %s=%s", k, v))
		} else if orig != v {
			lines = append(lines, fmt.Sprintf("~ %s: %q -> %q", k, orig, v))
		}
	}

	for k, v := range f.env {
		if _, ok := other[k]; !ok {
			lines = append(lines, fmt.Sprintf("- %s=%s", k, v))
		}
	}

	sort.Slice(lines, func(i, j int) bool {
		return strings.TrimLeft(lines[i], "+-~ ") < strings.TrimLeft(lines[j], "+-~ ")
	})
	return lines
}
