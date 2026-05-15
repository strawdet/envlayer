// Package aliaser provides key aliasing for environment variable maps.
//
// An Aliaser copies values from source keys to one or more target (alias)
// keys. This is useful when migrating from old key names to new ones while
// maintaining backward compatibility, or when multiple services expect the
// same value under different names.
//
// Example:
//
//	a := aliaser.New(
//		aliaser.WithAliases(map[string][]string{
//			"DB_HOST": {"DATABASE_HOST", "POSTGRES_HOST"},
//		}),
//		aliaser.WithKeepOriginal(true),
//	)
//	out, err := a.Apply(env)
//
// By default the original key is preserved. Pass WithKeepOriginal(false) to
// remove it after aliasing (effectively a rename).
package aliaser
