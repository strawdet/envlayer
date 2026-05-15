package aliaser

import (
	"testing"
)

func TestApply_BasicAlias(t *testing.T) {
	a := New(WithAliases(map[string][]string{
		"DB_HOST": {"DATABASE_HOST"},
	}))
	env := map[string]string{"DB_HOST": "localhost", "PORT": "5432"}
	out, err := a.Apply(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DATABASE_HOST"] != "localhost" {
		t.Errorf("expected alias value 'localhost', got %q", out["DATABASE_HOST"])
	}
	if out["DB_HOST"] != "localhost" {
		t.Error("original key should be retained by default")
	}
}

func TestApply_DropOriginal(t *testing.T) {
	a := New(
		WithAliases(map[string][]string{"OLD_KEY": {"NEW_KEY"}}),
		WithKeepOriginal(false),
	)
	out, err := a.Apply(map[string]string{"OLD_KEY": "value"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["OLD_KEY"]; ok {
		t.Error("original key should have been removed")
	}
	if out["NEW_KEY"] != "value" {
		t.Errorf("expected NEW_KEY=value, got %q", out["NEW_KEY"])
	}
}

func TestApply_MultipleTargets(t *testing.T) {
	a := New(WithAliases(map[string][]string{
		"SECRET": {"APP_SECRET", "SERVICE_SECRET"},
	}))
	out, err := a.Apply(map[string]string{"SECRET": "s3cr3t"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, k := range []string{"APP_SECRET", "SERVICE_SECRET", "SECRET"} {
		if out[k] != "s3cr3t" {
			t.Errorf("expected %s=s3cr3t, got %q", k, out[k])
		}
	}
}

func TestApply_MissingSourceSkipped(t *testing.T) {
	a := New(WithAliases(map[string][]string{"MISSING": {"ALIAS"}}))
	out, err := a.Apply(map[string]string{"OTHER": "val"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["ALIAS"]; ok {
		t.Error("alias should not be created when source is absent")
	}
}

func TestApply_TargetAlreadyExists_Error(t *testing.T) {
	a := New(WithAliases(map[string][]string{"SRC": {"EXISTING"}}))
	_, err := a.Apply(map[string]string{"SRC": "a", "EXISTING": "b"})
	if err == nil {
		t.Fatal("expected error when alias target already exists")
	}
}

func TestApply_NoOptions_NoChange(t *testing.T) {
	a := New()
	env := map[string]string{"A": "1", "B": "2"}
	out, err := a.Apply(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != len(env) {
		t.Errorf("expected %d keys, got %d", len(env), len(out))
	}
}

func TestAliases_ReturnsCopy(t *testing.T) {
	a := New(WithAliases(map[string][]string{"X": {"Y"}}))
	copy := a.Aliases()
	copy["X"] = append(copy["X"], "Z")
	if len(a.Aliases()["X"]) != 1 {
		t.Error("Aliases should return an independent copy")
	}
}
