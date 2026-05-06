package encryptor_test

import (
	"testing"

	"github.com/user/envlayer/internal/encryptor"
	"github.com/user/envlayer/internal/loader"
	"os"
	"path/filepath"
)

// writeEncEnv writes a simple .env file for the integration test.
func writeEncEnv(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	return p
}

// TestEncryptor_WithLoader verifies that values loaded from a .env file
// can be round-tripped through encryption and decryption.
func TestEncryptor_WithLoader(t *testing.T) {
	dir := t.TempDir()
	path := writeEncEnv(t, dir, ".env", "DB_PASS=hunter2\nAPI_KEY=xyz789\n")

	env, err := loader.LoadFile(path)
	if err != nil {
		t.Fatalf("LoadFile: %v", err)
	}

	e := encryptor.New("integration-passphrase")

	encrypted, err := e.EncryptMap(env)
	if err != nil {
		t.Fatalf("EncryptMap: %v", err)
	}

	// Encrypted values must not equal originals.
	for k, v := range env {
		if encrypted[k] == v {
			t.Errorf("key %s: encrypted value equals plaintext", k)
		}
	}

	decrypted, err := e.DecryptMap(encrypted)
	if err != nil {
		t.Fatalf("DecryptMap: %v", err)
	}

	// Decrypted values must match originals.
	for k, want := range env {
		if got := decrypted[k]; got != want {
			t.Errorf("key %s: want %q, got %q", k, want, got)
		}
	}
}
