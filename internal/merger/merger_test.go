package merger

import (
	"os"
	"path/filepath"
	"testing"
)

func setupDir(t *testing.T, files map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	for name, content := range files {
		path := filepath.Join(dir, name)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("setupDir: %v", err)
		}
	}
	return dir
}

func TestMerge_BasicOverride(t *testing.T) {
	dir := setupDir(t, map[string]string{
		".env":            "PORT=3000\nAPP_ENV=development\n",
		".env.production": "APP_ENV=production\nDEBUG=false\n",
	})

	m := New(dir)
	env, err := m.Merge(".env", ".env.production")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["APP_ENV"] != "production" {
		t.Errorf("APP_ENV: got %q, want %q", env["APP_ENV"], "production")
	}
	if env["PORT"] != "3000" {
		t.Errorf("PORT: got %q, want %q", env["PORT"], "3000")
	}
	if env["DEBUG"] != "false" {
		t.Errorf("DEBUG: got %q, want %q", env["DEBUG"], "false")
	}
}

func TestMerge_MissingLayerSkipped(t *testing.T) {
	dir := setupDir(t, map[string]string{
		".env": "KEY=base\n",
	})

	m := New(dir)
	env, err := m.Merge(".env", ".env.missing")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["KEY"] != "base" {
		t.Errorf("KEY: got %q", env["KEY"])
	}
}

func TestLayersForContext(t *testing.T) {
	layers := LayersForContext("staging")
	expected := []string{".env", ".env.staging", ".env.staging.local", ".env.local"}
	if len(layers) != len(expected) {
		t.Fatalf("layers length: got %d, want %d", len(layers), len(expected))
	}
	for i, l := range layers {
		if l != expected[i] {
			t.Errorf("layer[%d]: got %q, want %q", i, l, expected[i])
		}
	}
}

func TestLayersForContext_Empty(t *testing.T) {
	layers := LayersForContext("")
	if len(layers) != 2 || layers[0] != ".env" || layers[1] != ".env.local" {
		t.Errorf("unexpected layers for empty context: %v", layers)
	}
}
