// Package rotator provides key rotation utilities for environment variable maps.
// It supports renaming keys according to a rotation plan and tracking which keys
// have been rotated, allowing safe migration between naming conventions.
package rotator

import "fmt"

// RotationPlan maps old key names to new key names.
type RotationPlan map[string]string

// Result holds the outcome of a rotation operation.
type Result struct {
	Env     map[string]string
	Rotated []string
	Skipped []string
}

// Rotator applies key rotation plans to environment maps.
type Rotator struct {
	plan        RotationPlan
	keepOld     bool
	failMissing bool
}

// Option configures a Rotator.
type Option func(*Rotator)

// WithKeepOld retains the old key alongside the new key after rotation.
func WithKeepOld() Option {
	return func(r *Rotator) { r.keepOld = true }
}

// WithFailOnMissing causes Rotate to return an error if a planned old key is absent.
func WithFailOnMissing() Option {
	return func(r *Rotator) { r.failMissing = true }
}

// New creates a Rotator with the given plan and options.
func New(plan RotationPlan, opts ...Option) *Rotator {
	r := &Rotator{plan: plan}
	for _, o := range opts {
		o(r)
	}
	return r
}

// Rotate applies the rotation plan to env, returning a Result.
func (r *Rotator) Rotate(env map[string]string) (*Result, error) {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = v
	}

	result := &Result{Env: out}

	for oldKey, newKey := range r.plan {
		val, exists := out[oldKey]
		if !exists {
			if r.failMissing {
				return nil, fmt.Errorf("rotator: key %q not found in env", oldKey)
			}
			result.Skipped = append(result.Skipped, oldKey)
			continue
		}
		out[newKey] = val
		if !r.keepOld {
			delete(out, oldKey)
		}
		result.Rotated = append(result.Rotated, oldKey)
	}

	return result, nil
}

// Plan returns the rotation plan used by this Rotator.
func (r *Rotator) Plan() RotationPlan {
	copy := make(RotationPlan, len(r.plan))
	for k, v := range r.plan {
		copy[k] = v
	}
	return copy
}
