package templater

import (
	"testing"
)

func TestHasTemplates_True(t *testing.T) {
	if !HasTemplates("hello {{.NAME}}") {
		t.Fatal("expected true")
	}
}

func TestHasTemplates_False(t *testing.T) {
	if HasTemplates("hello $NAME") {
		t.Fatal("expected false")
	}
}

func TestRenderValue_Simple(t *testing.T) {
	tmpl := New(false)
	env := map[string]string{"APP_HOST": "localhost", "APP_PORT": "8080"}
	got, err := tmpl.RenderValue("{{.APP_HOST}}:{{.APP_PORT}}", env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "localhost:8080" {
		t.Errorf("got %q, want %q", got, "localhost:8080")
	}
}

func TestRenderValue_NoTemplate(t *testing.T) {
	tmpl := New(false)
	env := map[string]string{"X": "y"}
	got, err := tmpl.RenderValue("plain-value", env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "plain-value" {
		t.Errorf("got %q", got)
	}
}

func TestRenderValue_MissingKey_NonStrict(t *testing.T) {
	tmpl := New(false)
	env := map[string]string{}
	got, err := tmpl.RenderValue("prefix-{{.MISSING}}-suffix", env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "prefix--suffix" {
		t.Errorf("got %q, want %q", got, "prefix--suffix")
	}
}

func TestRenderValue_MissingKey_Strict(t *testing.T) {
	tmpl := New(true)
	env := map[string]string{}
	_, err := tmpl.RenderValue("{{.MISSING}}", env)
	if err == nil {
		t.Fatal("expected error in strict mode for missing key")
	}
}

func TestRenderAll_MultipleKeys(t *testing.T) {
	tmpl := New(false)
	env := map[string]string{
		"BASE_URL": "https://{{.HOST}}:{{.PORT}}",
		"HOST":     "example.com",
		"PORT":     "443",
		"STATIC":   "no-template",
	}
	out, err := tmpl.RenderAll(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "https://example.com:443"
	if out["BASE_URL"] != want {
		t.Errorf("BASE_URL: got %q, want %q", out["BASE_URL"], want)
	}
	if out["STATIC"] != "no-template" {
		t.Errorf("STATIC should be unchanged")
	}
}

func TestRenderAll_OriginalUnmodified(t *testing.T) {
	tmpl := New(false)
	env := map[string]string{"GREETING": "hello {{.NAME}}", "NAME": "world"}
	_, err := tmpl.RenderAll(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["GREETING"] != "hello {{.NAME}}" {
		t.Error("original env map was mutated")
	}
}
