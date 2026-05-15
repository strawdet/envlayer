package flattener_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/envlayer/internal/flattener"
	"github.com/your-org/envlayer/internal/loader"
)

func writeLayerFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeLayerFile: %v", err)
	}
	return p
}

func TestFlattener_WithLoader_ExtractsShortKeys(t *testing.T) {
	dir := t.TempDir()
	p := writeLayerFile(t, dir, ".env", "APP_HOST=localhost\nAPP_PORT=8080\nDB_URL=postgres://localhost/test\n")

	env, err := loader.LoadFile(p)
	if err != nil {
		t.Fatalf("LoadFile: %v", err)
	}

	f := flattener.New(flattener.WithPrefix("APP_"))
	out := f.Flatten(env)

	if out["HOST"] != "localhost" {
		t.Errorf("expected HOST=localhost, got %q", out["HOST"])
	}
	if out["PORT"] != "8080" {
		t.Errorf("expected PORT=8080, got %q", out["PORT"])
	}
	if _, ok := out["DB_URL"]; ok {
		t.Error("DB_URL should have been filtered by prefix")
	}
}

func TestFlattener_WithLoader_LowerAndCustomSep(t *testing.T) {
	dir := t.TempDir()
	p := writeLayerFile(t, dir, ".env", "svc.name=auth\nsvc.port=443\nother.key=ignored\n")

	env, err := loader.LoadFile(p)
	if err != nil {
		t.Fatalf("LoadFile: %v", err)
	}

	f := flattener.New(
		flattener.WithSeparator("."),
		flattener.WithPrefix("svc."),
		flattener.WithLowerKeys(),
	)
	out := f.Flatten(env)

	if out["name"] != "auth" {
		t.Errorf("expected name=auth, got %q", out["name"])
	}
	if out["port"] != "443" {
		t.Errorf("expected port=443, got %q", out["port"])
	}
	if len(out) != 2 {
		t.Errorf("expected 2 entries, got %d", len(out))
	}
}
