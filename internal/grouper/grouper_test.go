package grouper_test

import (
	"reflect"
	"sort"
	"testing"

	"github.com/envlayer/envlayer/internal/grouper"
)

func TestGroup_ByPrefix(t *testing.T) {
	g := grouper.New()
	env := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"APP_ENV":  "production",
		"NOPREFIX": "value",
	}
	groups := g.Group(env)

	if groups["DB"]["HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", groups["DB"]["HOST"])
	}
	if groups["DB"]["PORT"] != "5432" {
		t.Errorf("expected DB_PORT=5432, got %q", groups["DB"]["PORT"])
	}
	if groups["APP"]["ENV"] != "production" {
		t.Errorf("expected APP_ENV=production")
	}
	if groups[""]["NOPREFIX"] != "value" {
		t.Errorf("expected ungrouped NOPREFIX=value")
	}
}

func TestGroup_CustomSeparator(t *testing.T) {
	g := grouper.New(grouper.WithSeparator("."))
	env := map[string]string{
		"db.host": "localhost",
		"db.port": "5432",
	}
	groups := g.Group(env)
	if groups["db"]["host"] != "localhost" {
		t.Errorf("expected db.host=localhost")
	}
}

func TestPrefixes_Sorted(t *testing.T) {
	g := grouper.New()
	env := map[string]string{
		"Z_KEY":   "1",
		"A_KEY":   "2",
		"M_KEY":   "3",
		"NOGROUP": "4",
	}
	prefixes := g.Prefixes(env)
	expected := []string{"", "A", "M", "Z"}
	sort.Strings(expected)
	if !reflect.DeepEqual(prefixes, expected) {
		t.Errorf("expected %v, got %v", expected, prefixes)
	}
}

func TestFlatten_RoundTrip(t *testing.T) {
	g := grouper.New()
	original := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"APP_ENV": "staging",
	}
	grouped := g.Group(original)
	flat := g.Flatten(grouped)

	if !reflect.DeepEqual(flat, original) {
		t.Errorf("round-trip mismatch\ngot:  %v\nwant: %v", flat, original)
	}
}

func TestGroup_EmptyEnv(t *testing.T) {
	g := grouper.New()
	groups := g.Group(map[string]string{})
	if len(groups) != 0 {
		t.Errorf("expected empty groups, got %v", groups)
	}
}
