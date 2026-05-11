package filter_test

import (
	"testing"

	"envlayer/internal/filter"
)

func baseEnv() map[string]string {
	return map[string]string{
		"APP_HOST":   "localhost",
		"APP_PORT":   "8080",
		"DB_HOST":    "db.local",
		"DB_PASS":    "secret",
		"LOG_LEVEL":  "info",
	}
}

func TestApply_NoOptions(t *testing.T) {
	f := filter.New()
	out := f.Apply(baseEnv())
	if len(out) != 5 {
		t.Errorf("expected 5 keys, got %d", len(out))
	}
}

func TestApply_IncludeList(t *testing.T) {
	f := filter.New(filter.WithInclude("APP_HOST", "LOG_LEVEL"))
	out := f.Apply(baseEnv())
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
	if out["APP_HOST"] != "localhost" {
		t.Errorf("unexpected value for APP_HOST: %s", out["APP_HOST"])
	}
	if _, ok := out["DB_HOST"]; ok {
		t.Error("DB_HOST should be excluded")
	}
}

func TestApply_ExcludeList(t *testing.T) {
	f := filter.New(filter.WithExclude("DB_PASS", "LOG_LEVEL"))
	out := f.Apply(baseEnv())
	if len(out) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(out))
	}
	if _, ok := out["DB_PASS"]; ok {
		t.Error("DB_PASS should have been excluded")
	}
}

func TestApply_PrefixInclude(t *testing.T) {
	f := filter.New(filter.WithPrefixInclude("APP_"))
	out := f.Apply(baseEnv())
	if len(out) != 2 {
		t.Fatalf("expected 2 APP_ keys, got %d", len(out))
	}
	if _, ok := out["APP_HOST"]; !ok {
		t.Error("APP_HOST should be present")
	}
	if _, ok := out["APP_PORT"]; !ok {
		t.Error("APP_PORT should be present")
	}
}

func TestApply_PrefixInclude_MultiplePrefix(t *testing.T) {
	f := filter.New(filter.WithPrefixInclude("APP_", "DB_"))
	out := f.Apply(baseEnv())
	if len(out) != 4 {
		t.Fatalf("expected 4 keys, got %d", len(out))
	}
}

func TestApply_ExcludeOverridesPrefixInclude(t *testing.T) {
	f := filter.New(
		filter.WithPrefixInclude("DB_"),
		filter.WithExclude("DB_PASS"),
	)
	out := f.Apply(baseEnv())
	if _, ok := out["DB_PASS"]; ok {
		t.Error("DB_PASS should be excluded even with prefix match")
	}
	if _, ok := out["DB_HOST"]; !ok {
		t.Error("DB_HOST should be present")
	}
}

func TestApply_EmptyEnv(t *testing.T) {
	f := filter.New(filter.WithInclude("MISSING"))
	out := f.Apply(map[string]string{})
	if len(out) != 0 {
		t.Errorf("expected empty map, got %d keys", len(out))
	}
}
