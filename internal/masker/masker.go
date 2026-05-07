// Package masker provides utilities for masking sensitive environment
// variable values before display, logging, or export.
package masker

import "strings"

// DefaultSensitiveKeys contains common key substrings considered sensitive.
var DefaultSensitiveKeys = []string{
	"PASSWORD", "PASSWD", "SECRET", "TOKEN", "API_KEY",
	"PRIVATE_KEY", "AUTH", "CREDENTIAL", "CERT", "DSN",
}

// Masker masks sensitive environment variable values.
type Masker struct {
	sensitiveKeys []string
	maskChar      string
	revealChars   int
}

// Option configures a Masker.
type Option func(*Masker)

// WithSensitiveKeys overrides the default list of sensitive key substrings.
func WithSensitiveKeys(keys []string) Option {
	return func(m *Masker) { m.sensitiveKeys = keys }
}

// WithMaskChar sets the mask replacement string (default "****").
func WithMaskChar(c string) Option {
	return func(m *Masker) { m.maskChar = c }
}

// WithRevealChars sets how many leading characters to reveal (default 0).
func WithRevealChars(n int) Option {
	return func(m *Masker) { m.revealChars = n }
}

// New creates a Masker with the given options.
func New(opts ...Option) *Masker {
	m := &Masker{
		sensitiveKeys: DefaultSensitiveKeys,
		maskChar:      "****",
		revealChars:   0,
	}
	for _, o := range opts {
		o(m)
	}
	return m
}

// IsSensitive reports whether the key is considered sensitive.
func (m *Masker) IsSensitive(key string) bool {
	upper := strings.ToUpper(key)
	for _, s := range m.sensitiveKeys {
		if strings.Contains(upper, s) {
			return true
		}
	}
	return false
}

// Mask returns the masked form of value if the key is sensitive,
// otherwise it returns the value unchanged.
func (m *Masker) Mask(key, value string) string {
	if !m.IsSensitive(key) {
		return value
	}
	if m.revealChars > 0 && len(value) > m.revealChars {
		return value[:m.revealChars] + m.maskChar
	}
	return m.maskChar
}

// MaskMap returns a copy of env with sensitive values masked.
func (m *Masker) MaskMap(env map[string]string) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = m.Mask(k, v)
	}
	return out
}
