// Package interpolator resolves variable references embedded in environment
// variable values using ${VAR} or $VAR syntax.
//
// It operates on a map[string]string representing the merged environment and
// performs a single-pass substitution. References that cannot be resolved
// within the map are optionally looked up in the host OS environment when
// the fallbackToOS option is enabled.
//
// Example:
//
//	env := map[string]string{
//		"HOME":    "/home/alice",
//		"CONFIG":  "${HOME}/.config/app",
//	}
//
//	i := interpolator.New(false)
//	resolved := i.Interpolate(env)
//	// resolved["CONFIG"] == "/home/alice/.config/app"
//
// Use Validate to detect self-referencing keys before interpolation.
package interpolator
