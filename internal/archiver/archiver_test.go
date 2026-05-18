package archiver_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/nicholasgasior/envlayer/internal/archiver"
)

func tmpArchive(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "test.tar.gz")
}

func TestWrite_AndRead_RoundTrip(t *testing.T) {
	a := archiver.New()
	entries := []archiver.Entry{
		{Name: "base", Env: map[string]string{"APP": "envlayer", "PORT": "8080"}, CreatedAt: time.Now()},
		{Name: "prod", Env: map[string]string{"APP": "envlayer", "PORT": "443", "TLS": "true"}, CreatedAt: time.Now()},
	}
	path := tmpArchive(t)
	if err := a.Write(path, entries); err != nil {
		t.Fatalf("Write: %v", err)
	}
	got, err := a.Read(path)
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if len(got) != len(entries) {
		t.Fatalf("expected %d entries, got %d", len(entries), len(got))
	}
	for i, e := range entries {
		if got[i].Name != e.Name {
			t.Errorf("entry %d name: want %q got %q", i, e.Name, got[i].Name)
		}
		for k, v := range e.Env {
			if got[i].Env[k] != v {
				t.Errorf("entry %d key %q: want %q got %q", i, k, v, got[i].Env[k])
			}
		}
	}
}

func TestWrite_EmptyEntries(t *testing.T) {
	a := archiver.New()
	path := tmpArchive(t)
	if err := a.Write(path, nil); err != nil {
		t.Fatalf("Write empty: %v", err)
	}
	got, err := a.Read(path)
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected 0 entries, got %d", len(got))
	}
}

func TestRead_MissingFile(t *testing.T) {
	a := archiver.New()
	_, err := a.Read("/nonexistent/path.tar.gz")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestWrite_CreatesFile(t *testing.T) {
	a := archiver.New()
	path := tmpArchive(t)
	entry := archiver.Entry{Name: "dev", Env: map[string]string{"ENV": "dev"}, CreatedAt: time.Now()}
	if err := a.Write(path, []archiver.Entry{entry}); err != nil {
		t.Fatalf("Write: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("archive file not created: %v", err)
	}
}
