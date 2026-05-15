package rotator_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/envlayer/internal/loader"
	"github.com/your-org/envlayer/internal/rotator"
)

func writeRotatorEnv(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("writeRotatorEnv: %v", err)
	}
	return p
}

func TestRotator_WithLoader_RenamesKeys(t *testing.T) {
	dir := t.TempDir()
	path := writeRotatorEnv(t, dir, ".env", "DB_HOST=localhost\nDB_PORT=5432\nAPP_ENV=production\n")

	env, err := loader.LoadFile(path)
	if err != nil {
		t.Fatalf("LoadFile: %v", err)
	}

	plan := rotator.RotationPlan{
		"DB_HOST": "DATABASE_HOST",
		"DB_PORT": "DATABASE_PORT",
	}
	r := rotator.New(plan)
	res, err := r.Rotate(env)
	if err != nil {
		t.Fatalf("Rotate: %v", err)
	}

	if res.Env["DATABASE_HOST"] != "localhost" {
		t.Errorf("DATABASE_HOST: got %q, want %q", res.Env["DATABASE_HOST"], "localhost")
	}
	if res.Env["DATABASE_PORT"] != "5432" {
		t.Errorf("DATABASE_PORT: got %q, want %q", res.Env["DATABASE_PORT"], "5432")
	}
	if res.Env["APP_ENV"] != "production" {
		t.Errorf("APP_ENV should be preserved, got %q", res.Env["APP_ENV"])
	}
	if _, ok := res.Env["DB_HOST"]; ok {
		t.Error("old key DB_HOST should have been removed")
	}
	if len(res.Rotated) != 2 {
		t.Errorf("expected 2 rotated, got %d", len(res.Rotated))
	}
}

func TestRotator_WithLoader_KeepOldAndMissing(t *testing.T) {
	dir := t.TempDir()
	path := writeRotatorEnv(t, dir, ".env", "API_KEY=abc123\n")

	env, err := loader.LoadFile(path)
	if err != nil {
		t.Fatalf("LoadFile: %v", err)
	}

	plan := rotator.RotationPlan{
		"API_KEY":    "SERVICE_API_KEY",
		"OLD_SECRET": "NEW_SECRET",
	}
	r := rotator.New(plan, rotator.WithKeepOld())
	res, err := r.Rotate(env)
	if err != nil {
		t.Fatalf("Rotate: %v", err)
	}

	if res.Env["API_KEY"] != "abc123" {
		t.Error("expected old API_KEY to be kept with WithKeepOld")
	}
	if res.Env["SERVICE_API_KEY"] != "abc123" {
		t.Error("expected SERVICE_API_KEY to be set")
	}
	if len(res.Skipped) != 1 {
		t.Errorf("expected 1 skipped key, got %d", len(res.Skipped))
	}
}
