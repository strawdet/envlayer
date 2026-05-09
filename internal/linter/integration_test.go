package linter_test

import (
	"testing"

	"github.com/yourorg/envlayer/internal/linter"
	"github.com/yourorg/envlayer/internal/merger"
)

func TestLinter_WithMergerLayers(t *testing.T) {
	dir := t.TempDir()

	writeEnvFile(t, dir+"/.env", "APP_HOST=localhost\nAPP_PORT=8080\nDEBUG=\n")
	writeEnvFile(t, dir+"/.env.production", "APP_HOST=prod.example.com\nAPP_SECRET=supersecret\n")

	m := merger.New(dir)
	layers, err := m.LayersForContext("production")
	if err != nil {
		t.Fatalf("LayersForContext: %v", err)
	}

	l := linter.New(
		linter.WithDuplicateCheck(),
		linter.WithEmptyValueCheck(),
	)

	dupeFindings := l.LintLayers(layers)
	if len(dupeFindings) == 0 {
		t.Error("expected at least one duplicate finding for APP_HOST")
	}

	merged := make(map[string]string)
	for _, layer := range layers {
		for k, v := range layer {
			merged[k] = v
		}
	}

	valueFindings := l.Lint(merged)
	foundEmpty := false
	for _, f := range valueFindings {
		if f.Key == "DEBUG" && f.Severity == linter.SeverityWarning {
			foundEmpty = true
		}
	}
	if !foundEmpty {
		t.Error("expected empty value warning for DEBUG key")
	}
}

func writeEnvFile(t *testing.T, path, content string) {
	t.Helper()
	if err := writeFile(path, content); err != nil {
		t.Fatalf("writeEnvFile: %v", err)
	}
}

func writeFile(path, content string) error {
	import_os := func() error {
		import "os"
		return os.WriteFile(path, []byte(content), 0644)
	}
	return import_os()
}
