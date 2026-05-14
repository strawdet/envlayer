package redactor_test

import (
	"testing"

	"envlayer/internal/redactor"
)

func TestRedact_ByKey(t *testing.T) {
	r := redactor.New(redactor.WithKeys("SECRET", "PASSWORD"))
	env := map[string]string{
		"SECRET":   "super-secret",
		"PASSWORD": "hunter2",
		"APP_NAME": "envlayer",
	}
	out := r.Redact(env)
	if out["SECRET"] != "[REDACTED]" {
		t.Errorf("expected SECRET to be redacted, got %q", out["SECRET"])
	}
	if out["PASSWORD"] != "[REDACTED]" {
		t.Errorf("expected PASSWORD to be redacted, got %q", out["PASSWORD"])
	}
	if out["APP_NAME"] != "envlayer" {
		t.Errorf("expected APP_NAME to be unchanged, got %q", out["APP_NAME"])
	}
}

func TestRedact_ByPattern(t *testing.T) {
	r := redactor.New(redactor.WithPattern(`^sk_live_`))
	env := map[string]string{
		"STRIPE_KEY": "sk_live_abc123",
		"OTHER_KEY":  "pk_test_xyz",
	}
	out := r.Redact(env)
	if out["STRIPE_KEY"] != "[REDACTED]" {
		t.Errorf("expected STRIPE_KEY to be redacted, got %q", out["STRIPE_KEY"])
	}
	if out["OTHER_KEY"] != "pk_test_xyz" {
		t.Errorf("expected OTHER_KEY unchanged, got %q", out["OTHER_KEY"])
	}
}

func TestRedact_CustomText(t *testing.T) {
	r := redactor.New(
		redactor.WithKeys("TOKEN"),
		redactor.WithRedactedText("***"),
	)
	env := map[string]string{"TOKEN": "abc"}
	out := r.Redact(env)
	if out["TOKEN"] != "***" {
		t.Errorf("expected custom redacted text, got %q", out["TOKEN"])
	}
}

func TestRedact_NoOptions_NoChange(t *testing.T) {
	r := redactor.New()
	env := map[string]string{"FOO": "bar", "BAZ": "qux"}
	out := r.Redact(env)
	for k, v := range env {
		if out[k] != v {
			t.Errorf("expected %q=%q unchanged, got %q", k, v, out[k])
		}
	}
}

func TestRedactValue_PatternMatch(t *testing.T) {
	r := redactor.New(redactor.WithPattern(`^Bearer `))
	if got := r.RedactValue("Bearer token123"); got != "[REDACTED]" {
		t.Errorf("expected redacted, got %q", got)
	}
	if got := r.RedactValue("plain-value"); got != "plain-value" {
		t.Errorf("expected unchanged, got %q", got)
	}
}

func TestRedact_KeyCaseInsensitive(t *testing.T) {
	r := redactor.New(redactor.WithKeys("api_key"))
	env := map[string]string{"API_KEY": "sensitive"}
	out := r.Redact(env)
	if out["API_KEY"] != "[REDACTED]" {
		t.Errorf("expected case-insensitive key match, got %q", out["API_KEY"])
	}
}
