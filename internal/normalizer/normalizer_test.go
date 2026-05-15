package normalizer_test

import (
	"testing"

	"github.com/your-org/envlayer/internal/normalizer"
)

func TestApply_NoOptions(t *testing.T) {
	n := normalizer.New()
	env := map[string]string{"My-Key": "  hello  ", "OTHER": "world"}
	out := n.Apply(env)
	if out["My-Key"] != "  hello  " {
		t.Errorf("expected unchanged value, got %q", out["My-Key"])
	}
	if out["OTHER"] != "world" {
		t.Errorf("expected unchanged value, got %q", out["OTHER"])
	}
}

func TestApply_UpperKeys(t *testing.T) {
	n := normalizer.New(normalizer.WithUpperKeys())
	out := n.Apply(map[string]string{"my_key": "val", "Another": "x"})
	if _, ok := out["MY_KEY"]; !ok {
		t.Error("expected MY_KEY to exist")
	}
	if _, ok := out["ANOTHER"]; !ok {
		t.Error("expected ANOTHER to exist")
	}
}

func TestApply_LowerKeys(t *testing.T) {
	n := normalizer.New(normalizer.WithLowerKeys())
	out := n.Apply(map[string]string{"MY_KEY": "val"})
	if _, ok := out["my_key"]; !ok {
		t.Error("expected my_key to exist")
	}
}

func TestApply_UpperTakesPrecedenceOverLower(t *testing.T) {
	n := normalizer.New(normalizer.WithLowerKeys(), normalizer.WithUpperKeys())
	out := n.Apply(map[string]string{"Mixed": "v"})
	if _, ok := out["MIXED"]; !ok {
		t.Error("expected MIXED (upper wins)")
	}
}

func TestApply_ReplaceHyphens(t *testing.T) {
	n := normalizer.New(normalizer.WithReplaceHyphens())
	out := n.Apply(map[string]string{"my-key": "val", "no-dash-here": "x"})
	if _, ok := out["my_key"]; !ok {
		t.Error("expected my_key after hyphen replacement")
	}
	if _, ok := out["no_dash_here"]; !ok {
		t.Error("expected no_dash_here after hyphen replacement")
	}
}

func TestApply_TrimValues(t *testing.T) {
	n := normalizer.New(normalizer.WithTrimValues())
	out := n.Apply(map[string]string{"KEY": "  spaced  ", "OTHER": "\t\nnewline\n"})
	if out["KEY"] != "spaced" {
		t.Errorf("expected trimmed value, got %q", out["KEY"])
	}
	if out["OTHER"] != "newline" {
		t.Errorf("expected trimmed value, got %q", out["OTHER"])
	}
}

func TestApply_CombinedOptions(t *testing.T) {
	n := normalizer.New(
		normalizer.WithUpperKeys(),
		normalizer.WithReplaceHyphens(),
		normalizer.WithTrimValues(),
	)
	out := n.Apply(map[string]string{"my-service-url": "  http://localhost  "})
	if v, ok := out["MY_SERVICE_URL"]; !ok || v != "http://localhost" {
		t.Errorf("expected MY_SERVICE_URL=http://localhost, got %q=%q", "MY_SERVICE_URL", v)
	}
}

func TestApply_EmptyMap(t *testing.T) {
	n := normalizer.New(normalizer.WithUpperKeys())
	out := n.Apply(map[string]string{})
	if len(out) != 0 {
		t.Errorf("expected empty map, got %d entries", len(out))
	}
}
