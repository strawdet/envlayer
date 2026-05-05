package loader

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("writeTemp: %v", err)
	}
	return path
}

func TestLoadFile_Basic(t *testing.T) {
	path := writeTemp(t, "APP_ENV=production\nPORT=8080\n")
	env, err := LoadFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["APP_ENV"] != "production" {
		t.Errorf("APP_ENV: got %q, want %q", env["APP_ENV"], "production")
	}
	if env["PORT"] != "8080" {
		t.Errorf("PORT: got %q, want %q", env["PORT"], "8080")
	}
}

func TestLoadFile_CommentsAndBlanks(t *testing.T) {
	path := writeTemp(t, "# comment\n\nKEY=value\n")
	env, err := LoadFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(env) != 1 {
		t.Errorf("expected 1 entry, got %d", len(env))
	}
}

func TestLoadFile_QuotedValues(t *testing.T) {
	path := writeTemp(t, `DB_URL="postgres://localhost/mydb"` + "\nSECRET='abc123'\n")
	env, err := LoadFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["DB_URL"] != "postgres://localhost/mydb" {
		t.Errorf("DB_URL: got %q", env["DB_URL"])
	}
	if env["SECRET"] != "abc123" {
		t.Errorf("SECRET: got %q", env["SECRET"])
	}
}

func TestLoadFile_MissingFile(t *testing.T) {
	_, err := LoadFile("/nonexistent/.env")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadFile_InvalidLine(t *testing.T) {
	path := writeTemp(t, "BADLINE\n")
	_, err := LoadFile(path)
	if err == nil {
		t.Fatal("expected error for invalid line")
	}
}
