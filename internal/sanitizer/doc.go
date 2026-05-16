// Package sanitizer provides a Sanitizer that cleans environment variable
// keys and values before they are used downstream.
//
// Keys are sanitized by replacing characters outside [A-Za-z0-9_] with a
// configurable replacement character (default "_"). Leading and trailing
// replacement characters are trimmed. Keys that reduce to an empty string
// after sanitization are dropped from the output map.
//
// Values can optionally have ASCII control characters (0x00–0x1F, 0x7F)
// stripped, which is useful when loading env files that may contain
// embedded null bytes or other non-printable content.
//
// Usage:
//
//	s := sanitizer.New(
//		sanitizer.WithUpperKeys(),
//		sanitizer.WithStripControlChars(),
//	)
//	clean := s.Apply(rawEnv)
package sanitizer
