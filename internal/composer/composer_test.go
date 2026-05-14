package composer_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/envlayer/envlayer/internal/composer"
)

func writeLayer(t *testing.T, dir, name, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
		t.Fatalf("writeLayer: %v", err)
	}
}

func TestRun_BasicMergeAndExpand(t *testing.T) {
	dir := t.TempDir()
	writeLayer(t, dir, ".env", "APP=base\nSECRET=s3cr3t\n")
	writeLayer(t, dir, ".env.production", "APP=prod\nURL=https://${APP}.example.com\n")

	c := composer.New(composer.Options{
		BaseDir: dir,
		Context: "production",
	})
	res, err := c.Run()
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if res.Env["APP"] != "prod" {
		t.Errorf("APP = %q, want %q", res.Env["APP"], "prod")
	}
	if res.Env["URL"] != "https://prod.example.com" {
		t.Errorf("URL = %q, want expanded value", res.Env["URL"])
	}
}

func TestRun_RequiredKeyMissing_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	writeLayer(t, dir, ".env", "FOO=bar\n")

	c := composer.New(composer.Options{
		BaseDir:      dir,
		Context:      "",
		RequiredKeys: []string{"DATABASE_URL"},
	})
	_, err := c.Run()
	if err == nil {
		t.Fatal("expected error for missing required key, got nil")
	}
}

func TestRun_StripPrefixAndUpperCase(t *testing.T) {
	dir := t.TempDir()
	writeLayer(t, dir, ".env", "APP_host=localhost\nAPP_port=5432\n")

	c := composer.New(composer.Options{
		BaseDir:     dir,
		Context:     "",
		StripPrefix: "APP_",
		KeyCase:     "upper",
	})
	res, err := c.Run()
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if _, ok := res.Env["HOST"]; !ok {
		t.Errorf("expected key HOST after strip+upper, got keys: %v", res.Env)
	}
}

func TestRun_OutputToFile(t *testing.T) {
	dir := t.TempDir()
	writeLayer(t, dir, ".env", "KEY=value\n")
	out := filepath.Join(dir, "out.env")

	c := composer.New(composer.Options{
		BaseDir:      dir,
		Context:      "",
		OutputFormat: "dotenv",
		OutputPath:   out,
	})
	if _, err := c.Run(); err != nil {
		t.Fatalf("Run: %v", err)
	}
	if _, err := os.Stat(out); os.IsNotExist(err) {
		t.Errorf("expected output file %s to exist", out)
	}
}
