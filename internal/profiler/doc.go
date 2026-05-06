// Package profiler provides named environment profile management for envlayer.
//
// A profile associates a name with a context string and an ordered list of
// layer file patterns. Profiles are stored as a JSON array in a config file
// (e.g. .envlayer/profiles.json) and can be loaded, queried, and persisted
// at runtime.
//
// Example profiles.json:
//
//	[
//	  {
//	    "name": "dev",
//	    "context": "development",
//	    "layers": [".env", ".env.dev", ".env.local"]
//	  },
//	  {
//	    "name": "prod",
//	    "context": "production",
//	    "layers": [".env", ".env.prod"]
//	  }
//	]
//
// Usage:
//
//	p, err := profiler.New(".envlayer/profiles.json")
//	prof, ok := p.Get("dev")
package profiler
