package patcher_test

import (
	"testing"

	"github.com/envlayer/envlayer/internal/patcher"
)

func baseEnv() map[string]string {
	return map[string]string{
		"APP_ENV":  "development",
		"DB_HOST":  "localhost",
		"LOG_LEVEL": "info",
	}
}

func TestApply_SetNewKey(t *testing.T) {
	p := patcher.New()
	result, err := p.Apply(baseEnv(), []patcher.Patch{
		{Op: patcher.OpSet, Key: "NEW_KEY", Value: "hello"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["NEW_KEY"] != "hello" {
		t.Errorf("expected NEW_KEY=hello, got %q", result["NEW_KEY"])
	}
	if result["APP_ENV"] != "development" {
		t.Error("base key APP_ENV should be preserved")
	}
}

func TestApply_UpdateExistingKey(t *testing.T) {
	p := patcher.New()
	result, err := p.Apply(baseEnv(), []patcher.Patch{
		{Op: patcher.OpSet, Key: "APP_ENV", Value: "production"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["APP_ENV"] != "production" {
		t.Errorf("expected APP_ENV=production, got %q", result["APP_ENV"])
	}
}

func TestApply_DeleteExistingKey(t *testing.T) {
	p := patcher.New()
	result, err := p.Apply(baseEnv(), []patcher.Patch{
		{Op: patcher.OpDelete, Key: "LOG_LEVEL"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := result["LOG_LEVEL"]; ok {
		t.Error("LOG_LEVEL should have been deleted")
	}
}

func TestApply_DeleteUnknownKey_Error(t *testing.T) {
	p := patcher.New()
	_, err := p.Apply(baseEnv(), []patcher.Patch{
		{Op: patcher.OpDelete, Key: "DOES_NOT_EXIST"},
	})
	if err == nil {
		t.Fatal("expected error deleting unknown key, got nil")
	}
}

func TestApply_DeleteUnknownKey_Ignored(t *testing.T) {
	p := patcher.New(patcher.WithIgnoreUnknownDeletes())
	_, err := p.Apply(baseEnv(), []patcher.Patch{
		{Op: patcher.OpDelete, Key: "DOES_NOT_EXIST"},
	})
	if err != nil {
		t.Fatalf("expected no error with ignore option, got: %v", err)
	}
}

func TestApply_EmptyKeyReturnsError(t *testing.T) {
	p := patcher.New()
	_, err := p.Apply(baseEnv(), []patcher.Patch{
		{Op: patcher.OpSet, Key: "", Value: "oops"},
	})
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestApply_DoesNotMutateBase(t *testing.T) {
	base := baseEnv()
	p := patcher.New()
	_, _ = p.Apply(base, []patcher.Patch{
		{Op: patcher.OpSet, Key: "APP_ENV", Value: "staging"},
		{Op: patcher.OpDelete, Key: "DB_HOST"},
	})
	if base["APP_ENV"] != "development" {
		t.Error("base map was mutated: APP_ENV changed")
	}
	if _, ok := base["DB_HOST"]; !ok {
		t.Error("base map was mutated: DB_HOST was deleted")
	}
}

func TestDiff_DetectsChanges(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2", "C": "3"}
	dst := map[string]string{"A": "1", "B": "changed", "D": "4"}

	patches := patcher.Diff(src, dst)

	opsMap := map[string]patcher.Patch{}
	for _, p := range patches {
		opsMap[p.Key] = p
	}

	if p, ok := opsMap["B"]; !ok || p.Op != patcher.OpSet || p.Value != "changed" {
		t.Error("expected B to be updated")
	}
	if p, ok := opsMap["D"]; !ok || p.Op != patcher.OpSet || p.Value != "4" {
		t.Error("expected D to be added")
	}
	if p, ok := opsMap["C"]; !ok || p.Op != patcher.OpDelete {
		t.Error("expected C to be deleted")
	}
	if _, ok := opsMap["A"]; ok {
		t.Error("A is unchanged; should not appear in diff")
	}
}
