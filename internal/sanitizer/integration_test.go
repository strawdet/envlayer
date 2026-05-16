package sanitizer_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/envlayer/internal/loader"
	"github.com/your-org/envlayer/internal/sanitizer"
)

func writeSanitizerEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeSanitizerEnv: %v", err)
	}
	return p
}

func TestSanitizer_WithLoader_CleansKeys(t *testing.T) {
	path := writeSanitizerEnv(t, "my-key=hello\nfoo bar=world\nGOOD_KEY=ok\n")

	env, err := loader.LoadFile(path)
	if err != nil {
		t.Fatalf("LoadFile: %v", err)
	}

	s := sanitizer.New(sanitizer.WithUpperKeys())
	out := s.Apply(env)

	if out["MY_KEY"] != "hello" {
		t.Errorf("expected MY_KEY=hello, got %v", out)
	}
	if out["FOO_BAR"] != "world" {
		t.Errorf("expected FOO_BAR=world, got %v", out)
	}
	if out["GOOD_KEY"] != "ok" {
		t.Errorf("expected GOOD_KEY=ok, got %v", out)
	}
}

func TestSanitizer_WithLoader_StripsControlChars(t *testing.T) {
	path := writeSanitizerEnv(t, "SECRET=abc\x01def\nNAME=envlayer\n")

	env, err := loader.LoadFile(path)
	if err != nil {
		t.Fatalf("LoadFile: %v", err)
	}

	s := sanitizer.New(sanitizer.WithStripControlChars())
	out := s.Apply(env)

	if out["SECRET"] != "abcdef" {
		t.Errorf("expected 'abcdef', got %q", out["SECRET"])
	}
	if out["NAME"] != "envlayer" {
		t.Errorf("expected 'envlayer', got %q", out["NAME"])
	}
}
