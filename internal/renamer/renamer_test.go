package renamer_test

import (
	"testing"

	"github.com/nicholasgasior/envlayer/internal/renamer"
)

func TestApply_NoOptions(t *testing.T) {
	r := renamer.New()
	env := map[string]string{"FOO": "bar", "BAZ": "qux"}
	out := r.Apply(env)
	if out["FOO"] != "bar" || out["BAZ"] != "qux" {
		t.Errorf("expected unchanged map, got %v", out)
	}
}

func TestApply_ExplicitMapping(t *testing.T) {
	r := renamer.New(
		renamer.WithMapping("OLD_KEY", "NEW_KEY"),
	)
	env := map[string]string{"OLD_KEY": "value", "OTHER": "x"}
	out := r.Apply(env)
	if _, ok := out["OLD_KEY"]; ok {
		t.Error("expected OLD_KEY to be removed")
	}
	if out["NEW_KEY"] != "value" {
		t.Errorf("expected NEW_KEY=value, got %q", out["NEW_KEY"])
	}
	if out["OTHER"] != "x" {
		t.Errorf("expected OTHER=x, got %q", out["OTHER"])
	}
}

func TestApply_PatternRename(t *testing.T) {
	opt, err := renamer.WithPattern(`^APP_(.+)$`, `SVC_$1`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	r := renamer.New(opt)
	env := map[string]string{"APP_HOST": "localhost", "APP_PORT": "8080", "DB_URL": "postgres"}
	out := r.Apply(env)
	if out["SVC_HOST"] != "localhost" {
		t.Errorf("expected SVC_HOST=localhost, got %q", out["SVC_HOST"])
	}
	if out["SVC_PORT"] != "8080" {
		t.Errorf("expected SVC_PORT=8080, got %q", out["SVC_PORT"])
	}
	if out["DB_URL"] != "postgres" {
		t.Errorf("expected DB_URL unchanged, got %q", out["DB_URL"])
	}
}

func TestApply_MappingTakesPrecedenceOverPattern(t *testing.T) {
	opt, err := renamer.WithPattern(`^APP_(.+)$`, `SVC_$1`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	r := renamer.New(
		renamer.WithMapping("APP_HOST", "EXPLICIT_HOST"),
		opt,
	)
	env := map[string]string{"APP_HOST": "localhost"}
	out := r.Apply(env)
	if out["EXPLICIT_HOST"] != "localhost" {
		t.Errorf("expected EXPLICIT_HOST=localhost, got %v", out)
	}
	if _, ok := out["SVC_HOST"]; ok {
		t.Error("pattern should not have applied when explicit mapping exists")
	}
}

func TestWithPattern_InvalidRegex(t *testing.T) {
	_, err := renamer.WithPattern(`[invalid`, `X`)
	if err == nil {
		t.Error("expected error for invalid regex")
	}
}

func TestApply_CollisionKeepsOriginal(t *testing.T) {
	r := renamer.New(
		renamer.WithMapping("A", "B"),
	)
	// Both A and B exist; renaming A->B would collide with existing B
	env := map[string]string{"A": "fromA", "B": "fromB"}
	out := r.Apply(env)
	if len(out) != 2 {
		t.Errorf("expected 2 keys, got %d: %v", len(out), out)
	}
}
