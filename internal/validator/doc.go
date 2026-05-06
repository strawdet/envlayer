// Package validator provides utilities for validating merged environment
// variable maps produced by the resolver.
//
// It checks that:
//   - All keys conform to the POSIX environment variable naming convention
//     ([A-Za-z_][A-Za-z0-9_]*).
//   - All keys declared as required are present in the final map.
//   - Values with unusual characteristics (empty, multi-line) are surfaced
//     as warnings rather than hard errors.
//
// Usage:
//
//	v := validator.New([]string{"DATABASE_URL", "SECRET_KEY"})
//	result := v.Validate(env)
//	if !result.IsValid() {
//	    for _, e := range result.Errors {
//	        fmt.Fprintln(os.Stderr, "ERROR:", e)
//	    }
//	}
package validator
