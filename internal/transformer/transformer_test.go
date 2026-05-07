package transformer_test

import (
	"sort"
	"testing"

	"envlayer/internal/transformer"
)

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func TestApply_NoOptions(t *testing.T) {
	tr := transformer.New()
	env := map[string]string{"FOO": "bar", "BAZ": "qux"}
	out := tr.Apply(env)
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
}

func TestApply_PrefixFilter(t *testing.T) {
	tr := transformer.New(transformer.WithPrefixFilter("APP_"))
	env := map[string]string{"APP_HOST": "localhost", "APP_PORT": "8080", "DB_URL": "postgres"}
	out := tr.Apply(env)
	if len(out) != 2 {
		t.Fatalf("expected 2 keys after filter, got %d", len(out))
	}
	if _, ok := out["DB_URL"]; ok {
		t.Error("DB_URL should have been filtered out")
	}
}

func TestApply_StripPrefix(t *testing.T) {
	tr := transformer.New(
		transformer.WithPrefixFilter("APP_"),
		transformer.WithStripPrefix(true),
	)
	env := map[string]string{"APP_HOST": "localhost", "APP_PORT": "8080"}
	out := tr.Apply(env)
	if out["HOST"] != "localhost" {
		t.Errorf("expected HOST=localhost, got %q", out["HOST"])
	}
	if out["PORT"] != "8080" {
		t.Errorf("expected PORT=8080, got %q", out["PORT"])
	}
}

func TestApply_KeyCaseUpper(t *testing.T) {
	tr := transformer.New(transformer.WithKeyCase("upper"))
	env := map[string]string{"foo": "bar", "baz": "qux"}
	out := tr.Apply(env)
	if _, ok := out["FOO"]; !ok {
		t.Error("expected key FOO")
	}
}

func TestApply_KeyCaseLower(t *testing.T) {
	tr := transformer.New(transformer.WithKeyCase("lower"))
	env := map[string]string{"FOO": "bar"}
	out := tr.Apply(env)
	if _, ok := out["foo"]; !ok {
		t.Error("expected key foo")
	}
}

func TestApply_ValueCaseUpper(t *testing.T) {
	tr := transformer.New(transformer.WithValueCase("upper"))
	env := map[string]string{"MODE": "production"}
	out := tr.Apply(env)
	if out["MODE"] != "PRODUCTION" {
		t.Errorf("expected PRODUCTION, got %q", out["MODE"])
	}
}

func TestApply_EmptyEnv(t *testing.T) {
	tr := transformer.New(transformer.WithPrefixFilter("X_"))
	out := tr.Apply(map[string]string{})
	if len(out) != 0 {
		t.Errorf("expected empty map, got %d keys", len(out))
	}
}

func TestKeys_ReturnsList(t *testing.T) {
	tr := transformer.New(transformer.WithPrefixFilter("APP_"), transformer.WithStripPrefix(true))
	env := map[string]string{"APP_A": "1", "APP_B": "2", "OTHER": "3"}
	keys := tr.Keys(env)
	sort.Strings(keys)
	if len(keys) != 2 || keys[0] != "A" || keys[1] != "B" {
		t.Errorf("unexpected keys: %v", keys)
	}
}
