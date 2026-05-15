// Package flattener provides a Flattener that reduces multi-segment environment
// variable keys (e.g. APP_DB_HOST) to their final segment (HOST) or a chosen
// short form.
//
// This is useful when a downstream component expects flat, un-prefixed keys but
// the resolved environment still carries namespace prefixes from the merger or
// scoper layers.
//
// Basic usage:
//
//	f := flattener.New(
//		flattener.WithPrefix("APP_"),
//		flattener.WithLowerKeys(),
//	)
//	flat := f.Flatten(resolved)
//
// Separator defaults to "_" but can be changed with WithSeparator.
package flattener
