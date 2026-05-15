package rotator_test

import (
	"testing"

	"github.com/your-org/envlayer/internal/rotator"
)

func baseEnv() map[string]string {
	return map[string]string{
		"DB_HOST": "localhost",
		"DB_PASS": "secret",
		"APP_PORT": "8080",
	}
}

func TestRotate_RenamesKeys(t *testing.T) {
	plan := rotator.RotationPlan{"DB_HOST": "DATABASE_HOST", "DB_PASS": "DATABASE_PASSWORD"}
	r := rotator.New(plan)
	res, err := r.Rotate(baseEnv())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["DATABASE_HOST"] != "localhost" {
		t.Errorf("expected DATABASE_HOST=localhost, got %q", res.Env["DATABASE_HOST"])
	}
	if _, ok := res.Env["DB_HOST"]; ok {
		t.Error("expected old key DB_HOST to be removed")
	}
	if len(res.Rotated) != 2 {
		t.Errorf("expected 2 rotated keys, got %d", len(res.Rotated))
	}
}

func TestRotate_KeepOld(t *testing.T) {
	plan := rotator.RotationPlan{"DB_HOST": "DATABASE_HOST"}
	r := rotator.New(plan, rotator.WithKeepOld())
	res, err := r.Rotate(baseEnv())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["DB_HOST"] != "localhost" {
		t.Error("expected old key DB_HOST to be kept")
	}
	if res.Env["DATABASE_HOST"] != "localhost" {
		t.Error("expected new key DATABASE_HOST to be set")
	}
}

func TestRotate_MissingKeySkipped(t *testing.T) {
	plan := rotator.RotationPlan{"MISSING_KEY": "NEW_KEY"}
	r := rotator.New(plan)
	res, err := r.Rotate(baseEnv())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "MISSING_KEY" {
		t.Errorf("expected MISSING_KEY in skipped, got %v", res.Skipped)
	}
}

func TestRotate_FailOnMissing(t *testing.T) {
	plan := rotator.RotationPlan{"MISSING_KEY": "NEW_KEY"}
	r := rotator.New(plan, rotator.WithFailOnMissing())
	_, err := r.Rotate(baseEnv())
	if err == nil {
		t.Error("expected error for missing key, got nil")
	}
}

func TestRotate_UnchangedKeysPreserved(t *testing.T) {
	plan := rotator.RotationPlan{"DB_HOST": "DATABASE_HOST"}
	r := rotator.New(plan)
	res, err := r.Rotate(baseEnv())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["APP_PORT"] != "8080" {
		t.Error("expected APP_PORT to be preserved")
	}
}

func TestPlan_ReturnsCopy(t *testing.T) {
	plan := rotator.RotationPlan{"A": "B"}
	r := rotator.New(plan)
	p := r.Plan()
	p["X"] = "Y"
	if _, ok := r.Plan()["X"]; ok {
		t.Error("modifying returned plan should not affect rotator")
	}
}
