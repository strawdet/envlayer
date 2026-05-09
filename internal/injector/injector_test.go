package injector_test

import (
	"strings"
	"testing"

	"github.com/yourusername/envlayer/internal/injector"
)

func findKey(env []string, key string) (string, bool) {
	prefix := key + "="
	for _, kv := range env {
		if strings.HasPrefix(kv, prefix) {
			return strings.TrimPrefix(kv, prefix), true
		}
	}
	return "", false
}

func TestBuildEnv_InjectsNewKeys(t *testing.T) {
	base := []string{"EXISTING=yes", "PATH=/usr/bin"}
	inj := injector.New(injector.WithBaseEnv(base))

	resolved := map[string]string{"APP_ENV": "production", "DB_HOST": "localhost"}
	env := inj.BuildEnv(resolved)

	if v, ok := findKey(env, "APP_ENV"); !ok || v != "production" {
		t.Errorf("expected APP_ENV=production, got %q", v)
	}
	if v, ok := findKey(env, "DB_HOST"); !ok || v != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", v)
	}
	if v, ok := findKey(env, "EXISTING"); !ok || v != "yes" {
		t.Errorf("expected EXISTING=yes, got %q", v)
	}
}

func TestBuildEnv_NoOverride_ExistingWins(t *testing.T) {
	base := []string{"APP_ENV=staging"}
	inj := injector.New(injector.WithBaseEnv(base))

	resolved := map[string]string{"APP_ENV": "production"}
	env := inj.BuildEnv(resolved)

	if v, ok := findKey(env, "APP_ENV"); !ok || v != "staging" {
		t.Errorf("expected APP_ENV=staging (no override), got %q", v)
	}
}

func TestBuildEnv_WithOverride_ResolvedWins(t *testing.T) {
	base := []string{"APP_ENV=staging"}
	inj := injector.New(injector.WithBaseEnv(base), injector.WithOverride())

	resolved := map[string]string{"APP_ENV": "production"}
	env := inj.BuildEnv(resolved)

	if v, ok := findKey(env, "APP_ENV"); !ok || v != "production" {
		t.Errorf("expected APP_ENV=production (override), got %q", v)
	}
}

func TestBuildEnv_EmptyResolved(t *testing.T) {
	base := []string{"FOO=bar"}
	inj := injector.New(injector.WithBaseEnv(base))

	env := inj.BuildEnv(map[string]string{})

	if v, ok := findKey(env, "FOO"); !ok || v != "bar" {
		t.Errorf("expected FOO=bar, got %q", v)
	}
}

func TestBuildEnv_EmptyBase(t *testing.T) {
	inj := injector.New(injector.WithBaseEnv([]string{}))

	resolved := map[string]string{"ONLY_KEY": "only_val"}
	env := inj.BuildEnv(resolved)

	if v, ok := findKey(env, "ONLY_KEY"); !ok || v != "only_val" {
		t.Errorf("expected ONLY_KEY=only_val, got %q", v)
	}
	if len(env) != 1 {
		t.Errorf("expected 1 entry, got %d", len(env))
	}
}
