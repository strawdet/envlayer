// Package archiver provides gzip-compressed tar archive support for envlayer
// env snapshots.
//
// An archive can hold multiple named Entry values, each containing an env map
// and a creation timestamp. This makes it easy to bundle several layers (e.g.
// base, staging, production) into a single portable file that can be shared,
// stored in CI artefact storage, or committed to a secrets vault.
//
// Basic usage:
//
//	a := archiver.New()
//
//	entries := []archiver.Entry{
//		{Name: "base", Env: map[string]string{"APP": "myapp"}, CreatedAt: time.Now()},
//		{Name: "prod", Env: map[string]string{"PORT": "443"}, CreatedAt: time.Now()},
//	}
//
//	if err := a.Write("envs.tar.gz", entries); err != nil { ... }
//
//	loaded, err := a.Read("envs.tar.gz")
package archiver
