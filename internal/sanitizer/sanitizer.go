// Package sanitizer provides utilities for sanitizing environment variable
// keys and values by removing or replacing disallowed characters.
package sanitizer

import (
	"regexp"
	"strings"
)

var (
	defaultKeyPattern   = regexp.MustCompile(`[^A-Za-z0-9_]`)
	defaultValuePattern = regexp.MustCompile(`[\x00-\x1F\x7F]`)
)

// Option configures a Sanitizer.
type Option func(*Sanitizer)

// WithStripControlChars enables stripping of control characters from values.
func WithStripControlChars() Option {
	return func(s *Sanitizer) { s.stripControl = true }
}

// WithUpperKeys normalizes all keys to uppercase after sanitizing.
func WithUpperKeys() Option {
	return func(s *Sanitizer) { s.upperKeys = true }
}

// WithReplaceChar sets the replacement character used when an invalid
// character is found in a key (default: "_").
func WithReplaceChar(c string) Option {
	return func(s *Sanitizer) { s.replaceChar = c }
}

// Sanitizer cleans environment variable maps.
type Sanitizer struct {
	stripControl bool
	upperKeys   bool
	replaceChar string
}

// New creates a new Sanitizer with the given options.
func New(opts ...Option) *Sanitizer {
	s := &Sanitizer{replaceChar: "_"}
	for _, o := range opts {
		o(s)
	}
	return s
}

// Apply returns a sanitized copy of env.
func (s *Sanitizer) Apply(env map[string]string) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		cleanKey := s.sanitizeKey(k)
		cleanVal := s.sanitizeValue(v)
		if cleanKey != "" {
			out[cleanKey] = cleanVal
		}
	}
	return out
}

// SanitizeKey returns a sanitized version of a single key.
func (s *Sanitizer) SanitizeKey(k string) string { return s.sanitizeKey(k) }

// SanitizeValue returns a sanitized version of a single value.
func (s *Sanitizer) SanitizeValue(v string) string { return s.sanitizeValue(v) }

func (s *Sanitizer) sanitizeKey(k string) string {
	k = defaultKeyPattern.ReplaceAllString(k, s.replaceChar)
	k = strings.Trim(k, s.replaceChar)
	if s.upperKeys {
		k = strings.ToUpper(k)
	}
	return k
}

func (s *Sanitizer) sanitizeValue(v string) string {
	if s.stripControl {
		v = defaultValuePattern.ReplaceAllString(v, "")
	}
	return v
}
