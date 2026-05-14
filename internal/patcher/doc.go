// Package patcher provides fine-grained, non-destructive patching for
// environment variable maps.
//
// A Patch is a single atomic operation — either a set (add/update) or a
// delete — applied to a copy of an existing env map.  The original map is
// never mutated.
//
// Basic usage:
//
//	p := patcher.New()
//	updated, err := p.Apply(base, []patcher.Patch{
//		{Op: patcher.OpSet,    Key: "APP_ENV", Value: "production"},
//		{Op: patcher.OpDelete, Key: "DEBUG"},
//	})
//
// You can also compute the patches needed to transform one env into another:
//
//	patches := patcher.Diff(envA, envB)
//
// Options:
//
//	WithIgnoreUnknownDeletes — silently skip delete ops for keys that do
//	not exist in the base map instead of returning an error.
package patcher
