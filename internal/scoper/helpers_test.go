package scoper_test

import "os"

// writeEnvFileScoper writes content to path for integration tests.
func writeEnvFileScoper(path, content string) error {
	return os.WriteFile(path, []byte(content), 0o644)
}
