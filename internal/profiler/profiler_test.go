package profiler_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envlayer/internal/profiler"
)

func writeProfileConfig(t *testing.T, dir string, profiles []profiler.Profile) string {
	t.Helper()
	path := filepath.Join(dir, "profiles.json")
	data, _ := json.Marshal(profiles)
	_ = os.WriteFile(path, data, 0o644)
	return path
}

func TestNew_LoadsProfiles(t *testing.T) {
	dir := t.TempDir()
	path := writeProfileConfig(t, dir, []profiler.Profile{
		{Name: "dev", Context: "development", Layers: []string{".env", ".env.dev"}},
		{Name: "prod", Context: "production", Layers: []string{".env", ".env.prod"}},
	})
	p, err := profiler.New(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	dev, ok := p.Get("dev")
	if !ok {
		t.Fatal("expected 'dev' profile")
	}
	if dev.Context != "development" {
		t.Errorf("expected context 'development', got %q", dev.Context)
	}
	if len(dev.Layers) != 2 {
		t.Errorf("expected 2 layers, got %d", len(dev.Layers))
	}
}

func TestNew_MissingFile_NoError(t *testing.T) {
	dir := t.TempDir()
	p, err := profiler.New(filepath.Join(dir, "nonexistent.json"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(p.List()) != 0 {
		t.Errorf("expected empty profile list")
	}
}

func TestGet_UnknownProfile(t *testing.T) {
	dir := t.TempDir()
	p, _ := profiler.New(filepath.Join(dir, "profiles.json"))
	_, ok := p.Get("ghost")
	if ok {
		t.Error("expected not found for unknown profile")
	}
}

func TestSave_PersistsProfile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "profiles.json")
	p, _ := profiler.New(path)

	newProf := profiler.Profile{Name: "staging", Context: "staging", Layers: []string{".env", ".env.staging"}}
	if err := p.Save(newProf); err != nil {
		t.Fatalf("save error: %v", err)
	}

	p2, err := profiler.New(path)
	if err != nil {
		t.Fatalf("reload error: %v", err)
	}
	got, ok := p2.Get("staging")
	if !ok {
		t.Fatal("expected 'staging' profile after reload")
	}
	if got.Context != "staging" {
		t.Errorf("expected context 'staging', got %q", got.Context)
	}
}

func TestList_ReturnsAllNames(t *testing.T) {
	dir := t.TempDir()
	path := writeProfileConfig(t, dir, []profiler.Profile{
		{Name: "a", Context: "ctx-a"},
		{Name: "b", Context: "ctx-b"},
	})
	p, _ := profiler.New(path)
	names := p.List()
	if len(names) != 2 {
		t.Errorf("expected 2 profiles, got %d", len(names))
	}
}
