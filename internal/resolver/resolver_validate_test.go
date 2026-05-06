package resolver_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envlayer/internal/resolver"
	"github.com/yourorg/envlayer/internal/validator"
)

// TestResolveAndValidate_RequiredKeysPresent ensures that a fully resolved env
// passes validation when all required keys exist across layers.
func TestResolveAndValidate_RequiredKeysPresent(t *testing.T) {
	dir := t.TempDir()

	writeFile(t, filepath.Join(dir, ".env"), "APP_ENV=production\nDATABASE_URL=postgres://localhost/db")
	writeFile(t, filepath.Join(dir, ".env.production"), "SECRET_KEY=supersecret")

	r := resolver.New(dir)
	env, err := r.Resolve("production")
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}

	v := validator.New([]string{"DATABASE_URL", "SECRET_KEY"})
	res := v.Validate(env)
	if !res.IsValid() {
		t.Errorf("expected valid env, got errors: %v", res.Errors)
	}
}

// TestResolveAndValidate_MissingRequired ensures that missing required keys are
// caught after resolution.
func TestResolveAndValidate_MissingRequired(t *testing.T) {
	dir := t.TempDir()

	writeFile(t, filepath.Join(dir, ".env"), "APP_ENV=staging")

	r := resolver.New(dir)
	env, err := r.Resolve("staging")
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}

	v := validator.New([]string{"DATABASE_URL"})
	res := v.Validate(env)
	if res.IsValid() {
		t.Error("expected validation failure for missing DATABASE_URL")
	}
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("writeFile %s: %v", path, err)
	}
}
