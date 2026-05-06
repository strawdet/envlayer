package exporter_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"envlayer/internal/exporter"
)

var sampleVars = map[string]string{
	"APP_ENV":  "production",
	"DB_HOST":  "localhost",
	"LOG_LEVEL": "info",
}

func TestWrite_DotEnvFormat(t *testing.T) {
	var buf bytes.Buffer
	e := exporter.New(sampleVars, exporter.FormatDotEnv)
	if err := e.Write(&buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, key := range []string{"APP_ENV", "DB_HOST", "LOG_LEVEL"} {
		if !strings.Contains(out, key+"=") {
			t.Errorf("expected key %q in output, got:\n%s", key, out)
		}
	}
}

func TestWrite_ExportFormat(t *testing.T) {
	var buf bytes.Buffer
	e := exporter.New(sampleVars, exporter.FormatExport)
	if err := e.Write(&buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "export APP_ENV=") {
		t.Errorf("expected 'export APP_ENV=' in output, got:\n%s", out)
	}
}

func TestWrite_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	e := exporter.New(sampleVars, exporter.FormatJSON)
	if err := e.Write(&buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.HasPrefix(out, "{") || !strings.HasSuffix(strings.TrimSpace(out), "}") {
		t.Errorf("expected JSON object, got:\n%s", out)
	}
	if !strings.Contains(out, `"APP_ENV"`) {
		t.Errorf("expected APP_ENV key in JSON output, got:\n%s", out)
	}
}

func TestWrite_UnsupportedFormat(t *testing.T) {
	var buf bytes.Buffer
	e := exporter.New(sampleVars, exporter.Format("xml"))
	if err := e.Write(&buf); err == nil {
		t.Error("expected error for unsupported format, got nil")
	}
}

func TestWriteToFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "output.env")
	e := exporter.New(sampleVars, exporter.FormatDotEnv)
	if err := e.WriteToFile(path); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("could not read output file: %v", err)
	}
	if !strings.Contains(string(data), "DB_HOST=") {
		t.Errorf("expected DB_HOST in file output, got:\n%s", string(data))
	}
}

func TestWrite_SortedOutput(t *testing.T) {
	var buf bytes.Buffer
	e := exporter.New(sampleVars, exporter.FormatDotEnv)
	if err := e.Write(&buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "APP_ENV") {
		t.Errorf("expected first line to be APP_ENV (sorted), got: %s", lines[0])
	}
}
