// Package cloner provides deep-copy utilities for environment variable maps.
//
// It supports optional transformations during cloning:
//
//   - WithPrefixFilter: restrict cloned keys to those matching a prefix
//   - WithStripPrefix:  remove the matched prefix from keys in the output
//   - WithOmitEmpty:    exclude keys with empty string values
//
// Example:
//
//	c := cloner.New(
//		cloner.WithPrefixFilter("APP_"),
//		cloner.WithStripPrefix(),
//		cloner.WithOmitEmpty(),
//	)
//	result := c.Clone(env)
//
// The Merge method combines two maps, with source keys taking precedence
// over destination keys, without mutating either input.
package cloner
