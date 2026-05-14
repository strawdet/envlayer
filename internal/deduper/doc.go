// Package deduper detects and resolves duplicate keys across multiple
// environment variable layers.
//
// When merging .env files from different contexts (base, environment,
// profile, local), the same key may appear in more than one layer.
// Deduper provides explicit control over which value wins.
//
// Strategies:
//
//	KeepFirst – the value from the earliest layer is preserved.
//	KeepLast  – the value from the latest layer wins (default).
//
// Example:
//
//	d := deduper.New(deduper.WithStrategy(deduper.KeepFirst))
//	result := d.Merge(baseEnv, overrideEnv)
//
// Use Duplicates to inspect which keys collide before merging:
//
//	keys := deduper.Duplicates(layerA, layerB, layerC)
package deduper
