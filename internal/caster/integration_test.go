package caster_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envlayer/internal/caster"
	"github.com/yourorg/envlayer/internal/loader"
)

func writeCasterEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestCaster_WithLoader_IntAndBool(t *testing.T) {
	p := writeCasterEnv(t, "PORT=9090\nDEBUG=false\nRATIO=2.71\n")
	env, err := loader.LoadFile(p)
	if err != nil {
		t.Fatal(err)
	}

	c := caster.New(env)

	port, err := c.Int("PORT", 0)
	if err != nil || port != 9090 {
		t.Fatalf("port: want 9090/nil, got %d/%v", port, err)
	}

	debug, err := c.Bool("DEBUG", true)
	if err != nil || debug {
		t.Fatalf("debug: want false/nil, got %v/%v", debug, err)
	}

	ratio, err := c.Float64("RATIO", 0)
	if err != nil || ratio != 2.71 {
		t.Fatalf("ratio: want 2.71/nil, got %f/%v", ratio, err)
	}
}

func TestCaster_WithLoader_StrictRejectsGarbage(t *testing.T) {
	p := writeCasterEnv(t, "PORT=not-a-number\n")
	env, err := loader.LoadFile(p)
	if err != nil {
		t.Fatal(err)
	}

	c := caster.New(env, caster.WithStrict())
	_, err = c.Int("PORT", 0)
	if err == nil {
		t.Fatal("expected strict error for non-numeric PORT")
	}
}

func TestCaster_WithLoader_StringsList(t *testing.T) {
	p := writeCasterEnv(t, "ALLOWED_HOSTS=localhost, 127.0.0.1, ::1\n")
	env, err := loader.LoadFile(p)
	if err != nil {
		t.Fatal(err)
	}

	c := caster.New(env)
	hosts := c.Strings("ALLOWED_HOSTS", ",", nil)
	if len(hosts) != 3 {
		t.Fatalf("want 3 hosts, got %d: %v", len(hosts), hosts)
	}
	if hosts[0] != "localhost" {
		t.Fatalf("want localhost, got %q", hosts[0])
	}
}
