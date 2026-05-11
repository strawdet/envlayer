package scoper_test

import (
	"testing"

	"github.com/yourorg/envlayer/internal/merger"
	"github.com/yourorg/envlayer/internal/scoper"
)

// TestScoper_WithMerger verifies that scoper correctly partitions a merged
// environment produced by the merger into per-namespace sub-maps.
func TestScoper_WithMerger(t *testing.T) {
	dir := t.TempDir()
	writeScopedFile(t, dir+"/.env", "APP_HOST=base\nDB_URL=postgres://base\n")
	writeScopedFile(t, dir+"/.env.production", "APP_HOST=prod\nAPP_PORT=443\n")

	m := merger.New(dir)
	layers := merger.LayersForContext("production")
	merged, err := m.Merge(layers)
	if err != nil {
		t.Fatalf("merge failed: %v", err)
	}

	s := scoper.New()

	appEnv := s.Scope(merged, "APP")
	if appEnv["HOST"] != "prod" {
		t.Errorf("expected APP_HOST=prod after merge, got %q", appEnv["HOST"])
	}
	if appEnv["PORT"] != "443" {
		t.Errorf("expected APP_PORT=443, got %q", appEnv["PORT"])
	}

	dbEnv := s.Scope(merged, "DB")
	if dbEnv["URL"] != "postgres://base" {
		t.Errorf("expected DB_URL=postgres://base, got %q", dbEnv["URL"])
	}

	ns := s.Namespaces(merged)
	found := map[string]bool{}
	for _, n := range ns {
		found[n] = true
	}
	if !found["APP"] || !found["DB"] {
		t.Errorf("expected APP and DB namespaces, got %v", ns)
	}
}

func writeScopedFile(t *testing.T, path, content string) {
	t.Helper()
	if err := writeFileScoper(path, content); err != nil {
		t.Fatalf("writeScopedFile: %v", err)
	}
}

func writeFileScoper(path, content string) error {
	import_os := func() interface{} { return nil } // placeholder
	_ = import_os
	return writeEnvFileScoper(path, content)
}
