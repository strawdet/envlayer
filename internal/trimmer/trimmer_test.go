package trimmer_test

import (
	"testing"

	"github.com/yourorg/envlayer/internal/trimmer"
)

func TestApply_TrimValues_Default(t *testing.T) {
	tr := trimmer.New()
	env := map[string]string{
		"KEY": "  hello  ",
		"OTHER": "\tworld\t",
	}
	out := tr.Apply(env)
	if out["KEY"] != "hello" {
		t.Errorf("expected 'hello', got %q", out["KEY"])
	}
	if out["OTHER"] != "world" {
		t.Errorf("expected 'world', got %q", out["OTHER"])
	}
}

func TestApply_TrimKeys(t *testing.T) {
	tr := trimmer.New(trimmer.WithTrimKeys())
	env := map[string]string{"  SPACED  ": "value"}
	out := tr.Apply(env)
	if _, ok := out["SPACED"]; !ok {
		t.Error("expected trimmed key 'SPACED' to exist")
	}
	if _, ok := out["  SPACED  "]; ok {
		t.Error("expected original spaced key to be gone")
	}
}

func TestApply_CollapseWhitespace(t *testing.T) {
	tr := trimmer.New(trimmer.WithCollapseWhitespace())
	env := map[string]string{"MSG": "hello   world\there"}
	out := tr.Apply(env)
	if out["MSG"] != "hello world here" {
		t.Errorf("unexpected collapsed value: %q", out["MSG"])
	}
}

func TestApply_SkipKeys(t *testing.T) {
	tr := trimmer.New(trimmer.WithSkipKeys("RAW"))
	env := map[string]string{
		"RAW":   "  keep me  ",
		"CLEAN": "  trim me  ",
	}
	out := tr.Apply(env)
	if out["RAW"] != "  keep me  " {
		t.Errorf("expected raw value untouched, got %q", out["RAW"])
	}
	if out["CLEAN"] != "trim me" {
		t.Errorf("expected trimmed value, got %q", out["CLEAN"])
	}
}

func TestApply_NoOptions_StillTrimsValues(t *testing.T) {
	tr := trimmer.New()
	env := map[string]string{"A": " x "}
	out := tr.Apply(env)
	if out["A"] != "x" {
		t.Errorf("expected 'x', got %q", out["A"])
	}
}

func TestApply_EmptyEnv(t *testing.T) {
	tr := trimmer.New()
	out := tr.Apply(map[string]string{})
	if len(out) != 0 {
		t.Errorf("expected empty output, got %d keys", len(out))
	}
}

func TestApply_ReturnsCopy(t *testing.T) {
	tr := trimmer.New()
	original := map[string]string{"K": "  v  "}
	out := tr.Apply(original)
	out["K"] = "mutated"
	if original["K"] != "  v  " {
		t.Error("Apply should not mutate the original map")
	}
}
