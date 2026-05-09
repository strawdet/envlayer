package linter_test

import (
	"os"
	"testing"

	"github.com/yourorg/envlayer/internal/linter"
)

func TestLint_EmptyValue(t *testing.T) {
	l := linter.New(linter.WithEmptyValueCheck())
	env := map[string]string{
		"KEY_A": "value",
		"KEY_B": "",
	}
	findings := l.Lint(env)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Key != "KEY_B" {
		t.Errorf("expected KEY_B, got %s", findings[0].Key)
	}
	if findings[0].Severity != linter.SeverityWarning {
		t.Errorf("expected warning severity")
	}
}

func TestLint_OSShadow(t *testing.T) {
	os.Setenv("SHADOW_KEY", "original")
	t.Cleanup(func() { os.Unsetenv("SHADOW_KEY") })

	l := linter.New(linter.WithOSShadowCheck())
	env := map[string]string{
		"SHADOW_KEY": "overridden",
		"OTHER_KEY":  "value",
	}
	findings := l.Lint(env)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Key != "SHADOW_KEY" {
		t.Errorf("expected SHADOW_KEY, got %s", findings[0].Key)
	}
	if findings[0].Severity != linter.SeverityInfo {
		t.Errorf("expected info severity")
	}
}

func TestLint_NoFindings(t *testing.T) {
	l := linter.New(linter.WithEmptyValueCheck(), linter.WithOSShadowCheck())
	env := map[string]string{
		"CLEAN_KEY": "clean_value",
	}
	findings := l.Lint(env)
	if len(findings) != 0 {
		t.Errorf("expected no findings, got %d", len(findings))
	}
}

func TestLintLayers_DuplicateKeys(t *testing.T) {
	l := linter.New(linter.WithDuplicateCheck())
	layers := []map[string]string{
		{"APP_HOST": "localhost", "APP_PORT": "8080"},
		{"APP_HOST": "prod.example.com", "APP_SECRET": "abc123"},
	}
	findings := l.LintLayers(layers)
	if len(findings) != 1 {
		t.Fatalf("expected 1 duplicate finding, got %d", len(findings))
	}
	if findings[0].Key != "APP_HOST" {
		t.Errorf("expected APP_HOST duplicate, got %s", findings[0].Key)
	}
}

func TestLintLayers_NoDuplicates(t *testing.T) {
	l := linter.New(linter.WithDuplicateCheck())
	layers := []map[string]string{
		{"KEY_A": "1"},
		{"KEY_B": "2"},
	}
	findings := l.LintLayers(layers)
	if len(findings) != 0 {
		t.Errorf("expected no findings, got %d", len(findings))
	}
}

func TestFinding_String(t *testing.T) {
	f := linter.Finding{Key: "MY_KEY", Message: "some issue", Severity: linter.SeverityError}
	got := f.String()
	expected := "[error] MY_KEY: some issue"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}
