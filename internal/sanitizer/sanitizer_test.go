package sanitizer_test

import (
	"testing"

	"github.com/your-org/envlayer/internal/sanitizer"
)

func TestApply_CleanInput_Unchanged(t *testing.T) {
	s := sanitizer.New()
	env := map[string]string{"APP_HOST": "localhost", "PORT": "8080"}
	out := s.Apply(env)
	if out["APP_HOST"] != "localhost" || out["PORT"] != "8080" {
		t.Fatalf("expected unchanged output, got %v", out)
	}
}

func TestApply_InvalidKeyChars_Replaced(t *testing.T) {
	s := sanitizer.New()
	out := s.Apply(map[string]string{"my-key": "val", "foo bar": "baz"})
	if _, ok := out["my_key"]; !ok {
		t.Errorf("expected my_key, got keys: %v", out)
	}
	if _, ok := out["foo_bar"]; !ok {
		t.Errorf("expected foo_bar, got keys: %v", out)
	}
}

func TestApply_UpperKeys(t *testing.T) {
	s := sanitizer.New(sanitizer.WithUpperKeys())
	out := s.Apply(map[string]string{"app_name": "envlayer"})
	if out["APP_NAME"] != "envlayer" {
		t.Errorf("expected APP_NAME, got %v", out)
	}
}

func TestApply_StripControlChars(t *testing.T) {
	s := sanitizer.New(sanitizer.WithStripControlChars())
	out := s.Apply(map[string]string{"KEY": "val\x00ue\x1F"})
	if out["KEY"] != "value" {
		t.Errorf("expected 'value', got %q", out["KEY"])
	}
}

func TestApply_CustomReplaceChar(t *testing.T) {
	s := sanitizer.New(sanitizer.WithReplaceChar("-"))
	out := s.Apply(map[string]string{"my.key": "v"})
	if _, ok := out["my-key"]; !ok {
		t.Errorf("expected my-key, got %v", out)
	}
}

func TestApply_EmptyKeyDropped(t *testing.T) {
	s := sanitizer.New()
	out := s.Apply(map[string]string{"---": "val"})
	if len(out) != 0 {
		t.Errorf("expected empty map, got %v", out)
	}
}

func TestSanitizeKey_Standalone(t *testing.T) {
	s := sanitizer.New(sanitizer.WithUpperKeys())
	if got := s.SanitizeKey("foo-bar"); got != "FOO_BAR" {
		t.Errorf("expected FOO_BAR, got %s", got)
	}
}

func TestSanitizeValue_Standalone(t *testing.T) {
	s := sanitizer.New(sanitizer.WithStripControlChars())
	if got := s.SanitizeValue("hello\nworld"); got != "helloworld" {
		t.Errorf("expected 'helloworld', got %q", got)
	}
}
