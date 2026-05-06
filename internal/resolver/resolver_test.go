package resolver_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/envlayer/envlayer/internal/resolver"
)

func setupResolverDir(t *testing.T, files map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	for name, content := range files {
		path := filepath.Join(dir, name)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("setup: write %s: %v", name, err)
		}
	}
	return dir
}

func TestResolve_BasicMerge(t *testing.T) {
	dir := setupResolverDir(t, map[string]string{
		".env":            "APP_NAME=myapp\nDEBUG=false",
		".env.production": "DEBUG=true\nDB_HOST=prod-db",
	})

	r := resolver.New(dir)
	resolved, err := r.Resolve("production")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resolved.Vars["APP_NAME"] != "myapp" {
		t.Errorf("expected APP_NAME=myapp, got %q", resolved.Vars["APP_NAME"])
	}
	if resolved.Vars["DEBUG"] != "true" {
		t.Errorf("expected DEBUG=true (overridden), got %q", resolved.Vars["DEBUG"])
	}
	if resolved.Vars["DB_HOST"] != "prod-db" {
		t.Errorf("expected DB_HOST=prod-db, got %q", resolved.Vars["DB_HOST"])
	}

	if len(resolved.Layers) != 2 {
		t.Errorf("expected 2 loaded layers, got %d", len(resolved.Layers))
	}
}

func TestResolve_MissingContextFileSkipped(t *testing.T) {
	dir := setupResolverDir(t, map[string]string{
		".env": "BASE=1",
		// .env.staging intentionally absent
	})

	r := resolver.New(dir)
	resolved, err := r.Resolve("staging")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resolved.Vars["BASE"] != "1" {
		t.Errorf("expected BASE=1, got %q", resolved.Vars["BASE"])
	}
	if len(resolved.Layers) != 1 {
		t.Errorf("expected 1 loaded layer, got %d", len(resolved.Layers))
	}
}

func TestResolveToEnv_Format(t *testing.T) {
	dir := setupResolverDir(t, map[string]string{
		".env": "FOO=bar\nBAZ=qux",
	})

	r := resolver.New(dir)
	env, err := r.ResolveToEnv("development")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, entry := range env {
		if !strings.Contains(entry, "=") {
			t.Errorf("env entry missing '=': %q", entry)
		}
	}
}
