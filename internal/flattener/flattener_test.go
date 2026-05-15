package flattener_test

import (
	"sort"
	"testing"

	"github.com/your-org/envlayer/internal/flattener"
)

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func TestFlatten_DefaultSeparator(t *testing.T) {
	f := flattener.New()
	env := map[string]string{
		"APP_HOST": "localhost",
		"APP_PORT": "8080",
		"DB_HOST":  "db.local",
	}
	out := f.Flatten(env)
	if out["HOST"] != "db.local" && out["HOST"] != "localhost" {
		t.Errorf("expected HOST to be one of the values, got %q", out["HOST"])
	}
	if out["PORT"] != "8080" {
		t.Errorf("expected PORT=8080, got %q", out["PORT"])
	}
}

func TestFlatten_WithPrefix(t *testing.T) {
	f := flattener.New(flattener.WithPrefix("APP_"))
	env := map[string]string{
		"APP_HOST": "localhost",
		"APP_PORT": "8080",
		"DB_HOST":  "db.local",
	}
	out := f.Flatten(env)
	if _, ok := out["HOST"]; !ok {
		t.Error("expected HOST key from APP_ prefix")
	}
	if _, ok := out["DB_HOST"]; ok {
		t.Error("DB_HOST should have been excluded")
	}
	if len(out) != 2 {
		t.Errorf("expected 2 keys, got %d", len(out))
	}
}

func TestFlatten_LowerKeys(t *testing.T) {
	f := flattener.New(flattener.WithLowerKeys())
	env := map[string]string{"APP_HOST": "localhost"}
	out := f.Flatten(env)
	if out["host"] != "localhost" {
		t.Errorf("expected lower-case key 'host', got %v", sortedKeys(out))
	}
}

func TestFlatten_CustomSeparator(t *testing.T) {
	f := flattener.New(flattener.WithSeparator("."))
	env := map[string]string{
		"app.host": "localhost",
		"app.port": "9000",
	}
	out := f.Flatten(env)
	if out["host"] != "localhost" {
		t.Errorf("expected host=localhost, got %q", out["host"])
	}
	if out["port"] != "9000" {
		t.Errorf("expected port=9000, got %q", out["port"])
	}
}

func TestSegments_ReturnsLeadingParts(t *testing.T) {
	f := flattener.New()
	env := map[string]string{
		"APP_HOST": "localhost",
		"APP_PORT": "8080",
		"DB_HOST":  "db.local",
		"PLAIN":    "value",
	}
	segs := f.Segments(env)
	sort.Strings(segs)
	if len(segs) != 2 || segs[0] != "APP" || segs[1] != "DB" {
		t.Errorf("unexpected segments: %v", segs)
	}
}
