package interpolator_test

import (
	"os"
	"testing"

	"github.com/envlayer/envlayer/internal/interpolator"
)

func TestInterpolate_BraceStyle(t *testing.T) {
	i := interpolator.New(false)
	env := map[string]string{
		"BASE": "/home/user",
		"PATH": "${BASE}/bin",
	}
	result := i.Interpolate(env)
	if result["PATH"] != "/home/user/bin" {
		t.Errorf("expected /home/user/bin, got %q", result["PATH"])
	}
}

func TestInterpolate_DollarStyle(t *testing.T) {
	i := interpolator.New(false)
	env := map[string]string{
		"HOST": "localhost",
		"DSN":  "postgres://$HOST/db",
	}
	result := i.Interpolate(env)
	if result["DSN"] != "postgres://localhost/db" {
		t.Errorf("unexpected DSN: %q", result["DSN"])
	}
}

func TestInterpolate_UnresolvedNoFallback(t *testing.T) {
	i := interpolator.New(false)
	env := map[string]string{"VAL": "${MISSING}"}
	result := i.Interpolate(env)
	if result["VAL"] != "" {
		t.Errorf("expected empty string, got %q", result["VAL"])
	}
}

func TestInterpolate_FallbackToOS(t *testing.T) {
	os.Setenv("OS_VAR", "from-os")
	defer os.Unsetenv("OS_VAR")

	i := interpolator.New(true)
	env := map[string]string{"VAL": "${OS_VAR}"}
	result := i.Interpolate(env)
	if result["VAL"] != "from-os" {
		t.Errorf("expected from-os, got %q", result["VAL"])
	}
}

func TestInterpolate_NoReferences(t *testing.T) {
	i := interpolator.New(false)
	env := map[string]string{"PLAIN": "hello world"}
	result := i.Interpolate(env)
	if result["PLAIN"] != "hello world" {
		t.Errorf("unexpected value: %q", result["PLAIN"])
	}
}

func TestHasReferences(t *testing.T) {
	if !interpolator.HasReferences("${FOO}") {
		t.Error("expected true for ${FOO}")
	}
	if !interpolator.HasReferences("$BAR") {
		t.Error("expected true for $BAR")
	}
	if interpolator.HasReferences("plain") {
		t.Error("expected false for plain string")
	}
}

func TestValidate_SelfReference(t *testing.T) {
	env := map[string]string{"A": "${A}/extra"}
	if err := interpolator.Validate(env); err == nil {
		t.Error("expected self-reference error")
	}
}

func TestValidate_Valid(t *testing.T) {
	env := map[string]string{"A": "hello", "B": "${A}/world"}
	if err := interpolator.Validate(env); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
