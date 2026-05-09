package printer_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourusername/envlayer/internal/printer"
)

func TestPrint_KVFormat(t *testing.T) {
	p := printer.New()
	env := map[string]string{"APP_ENV": "production", "PORT": "8080"}
	var buf bytes.Buffer
	if err := p.Print(&buf, env, printer.FormatKV); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "APP_ENV=production") {
		t.Errorf("expected APP_ENV=production in output, got: %s", out)
	}
	if !strings.Contains(out, "PORT=8080") {
		t.Errorf("expected PORT=8080 in output, got: %s", out)
	}
}

func TestPrint_TableFormat(t *testing.T) {
	p := printer.New()
	env := map[string]string{"DB_HOST": "localhost"}
	var buf bytes.Buffer
	if err := p.Print(&buf, env, printer.FormatTable); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "KEY") || !strings.Contains(out, "VALUE") {
		t.Errorf("expected table header in output, got: %s", out)
	}
	if !strings.Contains(out, "DB_HOST") || !strings.Contains(out, "localhost") {
		t.Errorf("expected DB_HOST/localhost in output, got: %s", out)
	}
}

func TestPrint_JSONFormat(t *testing.T) {
	p := printer.New()
	env := map[string]string{"FOO": "bar"}
	var buf bytes.Buffer
	if err := p.Print(&buf, env, printer.FormatJSON); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `"FOO"`) || !strings.Contains(out, `"bar"`) {
		t.Errorf("expected JSON with FOO/bar, got: %s", out)
	}
}

func TestPrint_UnsupportedFormat(t *testing.T) {
	p := printer.New()
	var buf bytes.Buffer
	err := p.Print(&buf, map[string]string{"X": "1"}, "yaml")
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
	if !strings.Contains(err.Error(), "unsupported format") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestPrint_WithMasker(t *testing.T) {
	masker := func(key, value string) string {
		if strings.Contains(strings.ToLower(key), "secret") {
			return "***"
		}
		return value
	}
	p := printer.New(printer.WithMasker(masker))
	env := map[string]string{"APP_SECRET": "s3cr3t", "APP_ENV": "staging"}
	var buf bytes.Buffer
	if err := p.Print(&buf, env, printer.FormatKV); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if strings.Contains(out, "s3cr3t") {
		t.Errorf("expected secret to be masked, got: %s", out)
	}
	if !strings.Contains(out, "APP_ENV=staging") {
		t.Errorf("expected unmasked APP_ENV in output, got: %s", out)
	}
}

func TestPrint_SortedOutput(t *testing.T) {
	p := printer.New()
	env := map[string]string{"Z_KEY": "z", "A_KEY": "a", "M_KEY": "m"}
	var buf bytes.Buffer
	if err := p.Print(&buf, env, printer.FormatKV); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "A_KEY") {
		t.Errorf("expected first line to be A_KEY, got: %s", lines[0])
	}
	if !strings.HasPrefix(lines[2], "Z_KEY") {
		t.Errorf("expected last line to be Z_KEY, got: %s", lines[2])
	}
}
