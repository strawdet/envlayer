// Package rotator provides key rotation utilities for environment variable maps.
//
// A Rotator applies a RotationPlan — a mapping of old key names to new key names —
// to an environment map. This is useful when migrating between naming conventions
// or rotating secrets without breaking existing consumers immediately.
//
// Example usage:
//
//	plan := rotator.RotationPlan{
//		"DB_HOST": "DATABASE_HOST",
//		"DB_PASS": "DATABASE_PASSWORD",
//	}
//	r := rotator.New(plan, rotator.WithKeepOld())
//	result, err := r.Rotate(env)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(result.Rotated) // ["DB_HOST", "DB_PASS"]
//
// Options:
//   - WithKeepOld: retains the original key after rotation.
//   - WithFailOnMissing: returns an error if a planned key is absent.
package rotator
