package scoper_test

import (
	"os"
	"testing"
)

// writeEnvFileScoper writes content to path for integration tests.
func writeEnvFileScoper(path, content string) error {
	return os.WriteFile(path, []byte(content), 0o644)
}

// mustWriteEnvFileScoper writes content to path and fails the test on error.
func mustWriteEnvFileScoper(t *testing.T, path, content string) {
	t.Helper()
	if err := writeEnvFileScoper(path, content); err != nil {
		t.Fatalf("failed to write env file %q: %v", path, err)
	}
}
