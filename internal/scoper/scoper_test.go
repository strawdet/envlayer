package scoper_test

import (
	"sort"
	"testing"

	"github.com/yourorg/envlayer/internal/scoper"
)

func TestScope_ExtractsNamespace(t *testing.T) {
	s := scoper.New()
	env := map[string]string{
		"APP_HOST": "localhost",
		"APP_PORT": "8080",
		"DB_HOST":  "db.local",
	}
	got := s.Scope(env, "APP")
	if got["HOST"] != "localhost" {
		t.Errorf("expected HOST=localhost, got %q", got["HOST"])
	}
	if got["PORT"] != "8080" {
		t.Errorf("expected PORT=8080, got %q", got["PORT"])
	}
	if _, ok := got["DB_HOST"]; ok {
		t.Error("DB_HOST should not appear in APP scope")
	}
}

func TestScope_EmptyNamespace(t *testing.T) {
	s := scoper.New()
	env := map[string]string{"FOO": "bar"}
	got := s.Scope(env, "MISSING")
	if len(got) != 0 {
		t.Errorf("expected empty map, got %v", got)
	}
}

func TestNamespace_PrefixesKeys(t *testing.T) {
	s := scoper.New()
	env := map[string]string{"HOST": "localhost", "PORT": "9000"}
	got := s.Namespace(env, "SVC")
	if got["SVC_HOST"] != "localhost" {
		t.Errorf("expected SVC_HOST=localhost, got %q", got["SVC_HOST"])
	}
	if got["SVC_PORT"] != "9000" {
		t.Errorf("expected SVC_PORT=9000, got %q", got["SVC_PORT"])
	}
}

func TestNamespace_CustomSeparator(t *testing.T) {
	s := scoper.New(scoper.WithSeparator("."))
	env := map[string]string{"KEY": "val"}
	got := s.Namespace(env, "NS")
	if got["NS.KEY"] != "val" {
		t.Errorf("expected NS.KEY=val, got %v", got)
	}
}

func TestNamespaces_ReturnsDistinct(t *testing.T) {
	s := scoper.New()
	env := map[string]string{
		"APP_HOST": "h",
		"APP_PORT": "p",
		"DB_URL":   "u",
		"NOPREFIX": "x",
	}
	ns := s.Namespaces(env)
	sort.Strings(ns)
	if len(ns) != 2 || ns[0] != "APP" || ns[1] != "DB" {
		t.Errorf("unexpected namespaces: %v", ns)
	}
}
