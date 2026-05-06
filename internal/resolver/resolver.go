package resolver

import (
	"fmt"
	"os"

	"github.com/envlayer/envlayer/internal/loader"
	"github.com/envlayer/envlayer/internal/merger"
)

// Resolved holds the final merged environment variables and metadata.
type Resolved struct {
	Vars   map[string]string
	Layers []string
}

// Resolver combines the merger and loader to produce a final env map.
type Resolver struct {
	merger *merger.Merger
}

// New creates a new Resolver rooted at the given base directory.
func New(baseDir string) *Resolver {
	return &Resolver{
		merger: merger.New(baseDir),
	}
}

// Resolve loads and merges all .env layers for the given context (e.g. "production").
// Later layers override earlier ones. Missing layer files are silently skipped.
func (r *Resolver) Resolve(context string) (*Resolved, error) {
	layers := r.merger.LayersForContext(context)
	if len(layers) == 0 {
		return nil, fmt.Errorf("resolver: no layers defined for context %q", context)
	}

	merged := make(map[string]string)
	loaded := make([]string, 0, len(layers))

	for _, path := range layers {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		}

		vars, err := loader.LoadFile(path)
		if err != nil {
			return nil, fmt.Errorf("resolver: loading %s: %w", path, err)
		}

		for k, v := range vars {
			merged[k] = v
		}
		loaded = append(loaded, path)
	}

	return &Resolved{
		Vars:   merged,
		Layers: loaded,
	}, nil
}

// ResolveToEnv resolves the context and returns a slice of "KEY=VALUE" strings
// suitable for use with exec.Cmd.Env.
func (r *Resolver) ResolveToEnv(context string) ([]string, error) {
	resolved, err := r.Resolve(context)
	if err != nil {
		return nil, err
	}

	env := make([]string, 0, len(resolved.Vars))
	for k, v := range resolved.Vars {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}
	return env, nil
}
