// Package redactor provides pattern-based value redaction for environment variables.
// It replaces matching values with a configurable redaction string before output
// or logging, ensuring sensitive data is never exposed unintentionally.
package redactor

import (
	"regexp"
	"strings"
)

const defaultRedactedText = "[REDACTED]"

// Option configures a Redactor.
type Option func(*Redactor)

// Redactor replaces sensitive env values matched by pattern or key name.
type Redactor struct {
	patterns    []*regexp.Regexp
	keys        map[string]struct{}
	redactedText string
}

// WithPattern adds a regex pattern; any value matching it will be redacted.
func WithPattern(pattern string) Option {
	return func(r *Redactor) {
		if re, err := regexp.Compile(pattern); err == nil {
			r.patterns = append(r.patterns, re)
		}
	}
}

// WithKeys registers specific env keys whose values should always be redacted.
func WithKeys(keys ...string) Option {
	return func(r *Redactor) {
		for _, k := range keys {
			r.keys[strings.ToUpper(k)] = struct{}{}
		}
	}
}

// WithRedactedText sets the replacement string (default: "[REDACTED]").
func WithRedactedText(text string) Option {
	return func(r *Redactor) {
		r.redactedText = text
	}
}

// New creates a Redactor with the given options.
func New(opts ...Option) *Redactor {
	r := &Redactor{
		keys:         make(map[string]struct{}),
		redactedText: defaultRedactedText,
	}
	for _, o := range opts {
		o(r)
	}
	return r
}

// Redact returns a copy of env with sensitive values replaced.
func (r *Redactor) Redact(env map[string]string) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		if r.isSensitiveKey(k) || r.matchesPattern(v) {
			out[k] = r.redactedText
		} else {
			out[k] = v
		}
	}
	return out
}

// RedactValue redacts a single value if it matches any registered pattern.
func (r *Redactor) RedactValue(v string) string {
	if r.matchesPattern(v) {
		return r.redactedText
	}
	return v
}

func (r *Redactor) isSensitiveKey(k string) bool {
	_, ok := r.keys[strings.ToUpper(k)]
	return ok
}

func (r *Redactor) matchesPattern(v string) bool {
	for _, re := range r.patterns {
		if re.MatchString(v) {
			return true
		}
	}
	return false
}
