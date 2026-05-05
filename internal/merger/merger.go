package merger

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/yourorg/envlayer/internal/loader"
)

// Merger resolves and merges .env files in priority order.
type Merger struct {
	baseDir string
}

// New creates a Merger rooted at baseDir.
func New(baseDir string) *Merger {
	return &Merger{baseDir: baseDir}
}

// Merge loads .env files for the given layers (lowest to highest priority)
// and returns a single merged EnvMap. Later layers override earlier ones.
// Layers correspond to file names like ".env", ".env.production", etc.
func (m *Merger) Merge(layers ...string) (loader.EnvMap, error) {
	result := make(loader.EnvMap)

	for _, layer := range layers {
		path := filepath.Join(m.baseDir, layer)

		if _, err := os.Stat(path); os.IsNotExist(err) {
			// Missing layers are silently skipped.
			continue
		}

		env, err := loader.LoadFile(path)
		if err != nil {
			return nil, fmt.Errorf("merger: layer %q: %w", layer, err)
		}

		for k, v := range env {
			result[k] = v
		}
	}

	return result, nil
}

// LayersForContext returns the ordered list of .env file names to load
// for a given runtime context (e.g. "production", "staging").
func LayersForContext(context string) []string {
	layers := []string{".env"}
	if context != "" {
		layers = append(layers, ".env."+context)
		layers = append(layers, ".env."+context+".local")
	}
	layers = append(layers, ".env.local")
	return layers
}
