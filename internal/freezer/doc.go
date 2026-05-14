// Package freezer provides a write-once environment store for envlayer.
//
// A Freezer is populated with key-value pairs from a resolved environment map
// and then locked (frozen) to prevent further modification. Once frozen, the
// environment can be safely shared across goroutines without synchronisation.
//
// Typical usage:
//
//	f := freezer.New()
//	_ = f.Load(resolvedEnv)
//	f.Freeze()
//
//	// Safe read-only access
//	v, ok := f.Get("DATABASE_URL")
//
//	// Detect drift from the frozen baseline
//	diff := f.Diff(currentEnv)
//	for _, line := range diff {
//		fmt.Println(line)
//	}
//
// The Diff method returns lines prefixed with '+' (added), '-' (removed),
// or '~' (modified), sorted alphabetically by key for deterministic output.
package freezer
