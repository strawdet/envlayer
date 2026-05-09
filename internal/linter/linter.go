// Package linter provides static analysis for .env files,
// checking for common issues like duplicate keys, suspicious values,
// and keys that shadow OS environment variables.
package linter

import (
	"fmt"
	"os"
	"strings"
)

// Severity represents the level of a lint finding.
type Severity string

const (
	SeverityError   Severity = "error"
	SeverityWarning Severity = "warning"
	SeverityInfo    Severity = "info"
)

// Finding represents a single lint result.
type Finding struct {
	Key      string
	Message  string
	Severity Severity
}

func (f Finding) String() string {
	return fmt.Sprintf("[%s] %s: %s", f.Severity, f.Key, f.Message)
}

// Linter analyses a map of env vars and returns findings.
type Linter struct {
	checkDuplicates  bool
	checkOSShadow    bool
	checkEmptyValues bool
}

// Option configures the Linter.
type Option func(*Linter)

// WithDuplicateCheck enables detection of duplicate keys across a slice of maps.
func WithDuplicateCheck() Option { return func(l *Linter) { l.checkDuplicates = true } }

// WithOSShadowCheck warns when a key overrides an existing OS env var.
func WithOSShadowCheck() Option { return func(l *Linter) { l.checkOSShadow = true } }

// WithEmptyValueCheck warns when a key has an empty value.
func WithEmptyValueCheck() Option { return func(l *Linter) { l.checkEmptyValues = true } }

// New creates a Linter with the given options.
func New(opts ...Option) *Linter {
	l := &Linter{}
	for _, o := range opts {
		o(l)
	}
	return l
}

// Lint analyses the provided env map and returns all findings.
func (l *Linter) Lint(env map[string]string) []Finding {
	var findings []Finding

	for k, v := range env {
		if l.checkEmptyValues && strings.TrimSpace(v) == "" {
			findings = append(findings, Finding{
				Key:      k,
				Message:  "value is empty",
				Severity: SeverityWarning,
			})
		}
		if l.checkOSShadow {
			if osVal, exists := os.LookupEnv(k); exists && osVal != v {
				findings = append(findings, Finding{
					Key:      k,
					Message:  "shadows an existing OS environment variable",
					Severity: SeverityInfo,
				})
			}
		}
	}
	return findings
}

// LintLayers detects duplicate keys across multiple env layers.
func (l *Linter) LintLayers(layers []map[string]string) []Finding {
	var findings []Finding
	if !l.checkDuplicates {
		return findings
	}
	seen := map[string]int{}
	for i, layer := range layers {
		for k := range layer {
			if prev, ok := seen[k]; ok {
				findings = append(findings, Finding{
					Key:      k,
					Message:  fmt.Sprintf("duplicate key found in layer %d (first seen in layer %d)", i, prev),
					Severity: SeverityWarning,
				})
			} else {
				seen[k] = i
			}
		}
	}
	return findings
}
