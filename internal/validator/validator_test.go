package validator_test

import (
	"testing"

	"github.com/yourorg/envlayer/internal/validator"
)

func TestValidate_ValidEnv(t *testing.T) {
	v := validator.New(nil)
	env := map[string]string{
		"APP_ENV":  "production",
		"PORT":     "8080",
		"_SECRET":  "abc123",
	}
	res := v.Validate(env)
	if !res.IsValid() {
		t.Errorf("expected valid, got errors: %v", res.Errors)
	}
}

func TestValidate_InvalidKey(t *testing.T) {
	v := validator.New(nil)
	env := map[string]string{
		"123BAD": "value",
	}
	res := v.Validate(env)
	if res.IsValid() {
		t.Error("expected errors for invalid key")
	}
	if len(res.Errors) != 1 {
		t.Errorf("expected 1 error, got %d", len(res.Errors))
	}
}

func TestValidate_EmptyValueWarning(t *testing.T) {
	v := validator.New(nil)
	env := map[string]string{
		"EMPTY_KEY": "",
	}
	res := v.Validate(env)
	if !res.IsValid() {
		t.Errorf("unexpected errors: %v", res.Errors)
	}
	if len(res.Warnings) != 1 {
		t.Errorf("expected 1 warning, got %d", len(res.Warnings))
	}
}

func TestValidate_RequiredKeyMissing(t *testing.T) {
	v := validator.New([]string{"DATABASE_URL", "SECRET_KEY"})
	env := map[string]string{
		"APP_ENV": "staging",
	}
	res := v.Validate(env)
	if res.IsValid() {
		t.Error("expected errors for missing required keys")
	}
	if len(res.Errors) != 2 {
		t.Errorf("expected 2 errors, got %d: %v", len(res.Errors), res.Errors)
	}
}

func TestValidate_MultiLineWarning(t *testing.T) {
	v := validator.New(nil)
	env := map[string]string{
		"CERT": "line1\nline2",
	}
	res := v.Validate(env)
	if !res.IsValid() {
		t.Errorf("unexpected errors: %v", res.Errors)
	}
	if len(res.Warnings) == 0 {
		t.Error("expected multi-line warning")
	}
}
