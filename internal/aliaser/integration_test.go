package aliaser_test

import (
	"os"
	"path/filepath"
	"testing"

	"envlayer/internal/aliaser"
	"envlayer/internal/loader"
)

func writeAliasEnv(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeAliasEnv: %v", err)
	}
	return p
}

func TestAliaser_WithLoader_RenamesKey(t *testing.T) {
	dir := t.TempDir()
	p := writeAliasEnv(t, dir, ".env", "DB_HOST=postgres\nAPP_PORT=8080\n")

	env, err := loader.LoadFile(p)
	if err != nil {
		t.Fatalf("LoadFile: %v", err)
	}

	a := aliaser.New(
		aliaser.WithAliases(map[string][]string{
			"DB_HOST": {"DATABASE_HOST"},
		}),
		aliaser.WithKeepOriginal(false),
	)
	out, err := a.Apply(env)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if _, ok := out["DB_HOST"]; ok {
		t.Error("original DB_HOST should have been removed")
	}
	if out["DATABASE_HOST"] != "postgres" {
		t.Errorf("expected DATABASE_HOST=postgres, got %q", out["DATABASE_HOST"])
	}
	if out["APP_PORT"] != "8080" {
		t.Error("unrelated keys should be preserved")
	}
}

func TestAliaser_WithLoader_KeepOriginal(t *testing.T) {
	dir := t.TempDir()
	p := writeAliasEnv(t, dir, ".env", "SECRET_KEY=abc123\n")

	env, err := loader.LoadFile(p)
	if err != nil {
		t.Fatalf("LoadFile: %v", err)
	}

	a := aliaser.New(
		aliaser.WithAliases(map[string][]string{
			"SECRET_KEY": {"APP_SECRET", "SERVICE_TOKEN"},
		}),
	)
	out, err := a.Apply(env)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	for _, k := range []string{"SECRET_KEY", "APP_SECRET", "SERVICE_TOKEN"} {
		if out[k] != "abc123" {
			t.Errorf("expected %s=abc123, got %q", k, out[k])
		}
	}
}
