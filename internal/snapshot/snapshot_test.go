package snapshot_test

import (
	"os"
	"testing"

	"github.com/envlayer/envlayer/internal/snapshot"
)

func TestSave_AndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	m, err := snapshot.New(dir)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	env := map[string]string{"APP_ENV": "production", "DB_HOST": "localhost"}
	snap, err := m.Save("prod-snapshot", env)
	if err != nil {
		t.Fatalf("Save: %v", err)
	}
	if snap.Label != "prod-snapshot" {
		t.Errorf("expected label prod-snapshot, got %s", snap.Label)
	}

	loaded, err := m.Load("prod-snapshot")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded.Env["APP_ENV"] != "production" {
		t.Errorf("expected APP_ENV=production, got %s", loaded.Env["APP_ENV"])
	}
	if loaded.Env["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %s", loaded.Env["DB_HOST"])
	}
}

func TestLoad_NotFound(t *testing.T) {
	dir := t.TempDir()
	m, _ := snapshot.New(dir)

	_, err := m.Load("nonexistent")
	if err == nil {
		t.Fatal("expected error for missing snapshot, got nil")
	}
}

func TestDelete_RemovesSnapshot(t *testing.T) {
	dir := t.TempDir()
	m, _ := snapshot.New(dir)

	env := map[string]string{"KEY": "value"}
	_, err := m.Save("temp", env)
	if err != nil {
		t.Fatalf("Save: %v", err)
	}

	if err := m.Delete("temp"); err != nil {
		t.Fatalf("Delete: %v", err)
	}

	if _, err := m.Load("temp"); err == nil {
		t.Fatal("expected error after delete, got nil")
	}
}

func TestNew_CreatesDir(t *testing.T) {
	dir := t.TempDir() + "/nested/snapshots"
	_, err := snapshot.New(dir)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Error("expected directory to be created")
	}
}

func TestSave_OverwritesExisting(t *testing.T) {
	dir := t.TempDir()
	m, _ := snapshot.New(dir)

	_, _ = m.Save("overwrite", map[string]string{"A": "1"})
	_, err := m.Save("overwrite", map[string]string{"A": "2"})
	if err != nil {
		t.Fatalf("second Save: %v", err)
	}

	loaded, _ := m.Load("overwrite")
	if loaded.Env["A"] != "2" {
		t.Errorf("expected A=2 after overwrite, got %s", loaded.Env["A"])
	}
}
