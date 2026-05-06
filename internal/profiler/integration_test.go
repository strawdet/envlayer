package profiler_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envlayer/internal/merger"
	"github.com/user/envlayer/internal/profiler"
)

// TestProfileIntegration_WithMerger verifies that a loaded profile's context
// and layers can be used directly with the merger to resolve env files.
func TestProfileIntegration_WithMerger(t *testing.T) {
	dir := t.TempDir()

	// Write .env files
	_ = os.WriteFile(filepath.Join(dir, ".env"), []byte("APP=base\nDEBUG=false\n"), 0o644)
	_ = os.WriteFile(filepath.Join(dir, ".env.dev"), []byte("DEBUG=true\n"), 0o644)

	// Write profile config
	profPath := writeProfileConfig(t, dir, []profiler.Profile{
		{
			Name:    "dev",
			Context: "dev",
			Layers:  []string{".env", ".env.dev"},
		},
	})

	p, err := profiler.New(profPath)
	if err != nil {
		t.Fatalf("profiler.New: %v", err)
	}

	prof, ok := p.Get("dev")
	if !ok {
		t.Fatal("expected 'dev' profile")
	}

	m := merger.New(dir)
	result, err := m.Merge(prof.Layers)
	if err != nil {
		t.Fatalf("merger.Merge: %v", err)
	}

	if result["APP"] != "base" {
		t.Errorf("expected APP=base, got %q", result["APP"])
	}
	if result["DEBUG"] != "true" {
		t.Errorf("expected DEBUG=true after dev override, got %q", result["DEBUG"])
	}
}
