// Package validator provides validation for environment variable keys and values.
package validator

import (
	"fmt"
	"regexp"
	"strings"
)

var validKeyPattern = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

// Result holds the outcome of a validation pass.
type Result struct {
	Errors   []string
	Warnings []string
}

// IsValid returns true when there are no errors.
func (r Result) IsValid() bool {
	return len(r.Errors) == 0
}

// Validator checks environment variable maps for common problems.
type Validator struct {
	RequiredKeys []string
}

// New creates a Validator. requiredKeys may be empty.
func New(requiredKeys []string) *Validator {
	return &Validator{RequiredKeys: requiredKeys}
}

// Validate inspects env and returns a Result describing any issues found.
func (v *Validator) Validate(env map[string]string) Result {
	var res Result

	for k, val := range env {
		if !validKeyPattern.MatchString(k) {
			res.Errors = append(res.Errors, fmt.Sprintf("invalid key %q: must match [A-Za-z_][A-Za-z0-9_]*", k))
		}
		if strings.Contains(val, "\n") {
			res.Warnings = append(res.Warnings, fmt.Sprintf("key %q has a multi-line value", k))
		}
		if val == "" {
			res.Warnings = append(res.Warnings, fmt.Sprintf("key %q has an empty value", k))
		}
	}

	for _, req := range v.RequiredKeys {
		if _, ok := env[req]; !ok {
			res.Errors = append(res.Errors, fmt.Sprintf("required key %q is missing", req))
		}
	}

	return res
}
