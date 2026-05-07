// Package snapshot provides point-in-time capture and restoration of
// resolved environment variable states.
//
// A Snapshot records a labelled, timestamped copy of a merged env map
// and persists it as a JSON file on disk. Snapshots can be loaded back
// later for diffing, auditing, or rollback workflows.
//
// Basic usage:
//
//	m, err := snapshot.New(".envlayer/snapshots")
//	if err != nil { ... }
//
//	// Save current resolved env
//	snap, err := m.Save("before-deploy", resolvedEnv)
//
//	// Load it back later
//	snap, err = m.Load("before-deploy")
//	fmt.Println(snap.CapturedAt, snap.Env)
//
// Snapshots integrate naturally with the differ package to compare
// two captured states and surface what changed between deployments.
package snapshot
