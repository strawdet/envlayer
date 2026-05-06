// Package differ provides utilities for comparing two environment variable
// maps and identifying what changed between them.
//
// It is useful for detecting drift between environment layers (e.g. base vs
// production overrides), auditing changes between reloads, or surfacing
// meaningful diffs to the user before applying a new configuration.
//
// Basic usage:
//
//	d := differ.New()
//	changes := d.Compare(previousEnv, currentEnv)
//	for _, c := range changes {
//		fmt.Printf("%s %s: %q -> %q\n", c.Type, c.Key, c.OldVal, c.NewVal)
//	}
//
// Change types:
//   - Added:    key present in next but not in base
//   - Removed:  key present in base but not in next
//   - Modified: key present in both but with different values
package differ
