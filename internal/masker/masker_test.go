package masker_test

import (
	"testing"

	"github.com/user/envlayer/internal/masker"
)

func TestIsSensitive_DefaultKeys(t *testing.T) {
	m := masker.New()
	cases := []struct {
		key       string
		wantSens  bool
	}{
		{"DB_PASSWORD", true},
		{"API_KEY", true},
		{"AUTH_TOKEN", true},
		{"APP_NAME", false},
		{"PORT", false},
		{"PRIVATE_KEY_PATH", true},
	}
	for _, tc := range cases {
		got := m.IsSensitive(tc.key)
		if got != tc.wantSens {
			t.Errorf("IsSensitive(%q) = %v, want %v", tc.key, got, tc.wantSens)
		}
	}
}

func TestMask_SensitiveValue(t *testing.T) {
	m := masker.New()
	got := m.Mask("DB_PASSWORD", "supersecret")
	if got != "****" {
		t.Errorf("expected masked value, got %q", got)
	}
}

func TestMask_NonSensitiveValue(t *testing.T) {
	m := masker.New()
	got := m.Mask("APP_ENV", "production")
	if got != "production" {
		t.Errorf("expected plain value, got %q", got)
	}
}

func TestMask_RevealChars(t *testing.T) {
	m := masker.New(masker.WithRevealChars(3))
	got := m.Mask("API_KEY", "abcdefgh")
	if got != "abc****" {
		t.Errorf("expected partial reveal, got %q", got)
	}
}

func TestMask_RevealChars_ShortValue(t *testing.T) {
	m := masker.New(masker.WithRevealChars(10))
	got := m.Mask("SECRET", "abc")
	// value shorter than revealChars — fully masked
	if got != "****" {
		t.Errorf("expected full mask for short value, got %q", got)
	}
}

func TestMaskMap(t *testing.T) {
	m := masker.New()
	env := map[string]string{
		"APP_ENV":     "staging",
		"DB_PASSWORD": "s3cr3t",
		"API_KEY":     "key-abc",
	}
	masked := m.MaskMap(env)
	if masked["APP_ENV"] != "staging" {
		t.Errorf("APP_ENV should be unchanged")
	}
	if masked["DB_PASSWORD"] != "****" {
		t.Errorf("DB_PASSWORD should be masked")
	}
	if masked["API_KEY"] != "****" {
		t.Errorf("API_KEY should be masked")
	}
	// original map must be unchanged
	if env["DB_PASSWORD"] != "s3cr3t" {
		t.Errorf("original map should not be mutated")
	}
}

func TestWithCustomMaskChar(t *testing.T) {
	m := masker.New(masker.WithMaskChar("[REDACTED]"))
	got := m.Mask("AUTH_TOKEN", "tok_xyz")
	if got != "[REDACTED]" {
		t.Errorf("expected custom mask char, got %q", got)
	}
}

func TestWithCustomSensitiveKeys(t *testing.T) {
	m := masker.New(masker.WithSensitiveKeys([]string{"INTERNAL"}))
	if !m.IsSensitive("INTERNAL_ID") {
		t.Error("INTERNAL_ID should be sensitive with custom keys")
	}
	if m.IsSensitive("API_KEY") {
		t.Error("API_KEY should NOT be sensitive with custom keys override")
	}
}
