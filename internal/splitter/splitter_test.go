package splitter_test

import (
	"sort"
	"testing"

	"github.com/yourorg/envlayer/internal/splitter"
)

func sortedStrings(ss []string) []string {
	out := make([]string, len(ss))
	copy(out, ss)
	sort.Strings(out)
	return out
}

func TestSplit_ByPrefix(t *testing.T) {
	s := splitter.New()
	env := map[string]string{
		"DB_HOST":  "localhost",
		"DB_PORT":  "5432",
		"APP_NAME": "envlayer",
	}
	buckets := s.Split(env)
	if buckets["DB"]["HOST"] != "localhost" {
		t.Errorf("expected DB/HOST=localhost, got %q", buckets["DB"]["HOST"])
	}
	if buckets["DB"]["PORT"] != "5432" {
		t.Errorf("expected DB/PORT=5432, got %q", buckets["DB"]["PORT"])
	}
	if buckets["APP"]["NAME"] != "envlayer" {
		t.Errorf("expected APP/NAME=envlayer, got %q", buckets["APP"]["NAME"])
	}
}

func TestSplit_DefaultBucket(t *testing.T) {
	s := splitter.New()
	env := map[string]string{
		"NOPREFIXKEY": "value",
		"DB_HOST":     "localhost",
	}
	buckets := s.Split(env)
	if buckets["_default"]["NOPREFIXKEY"] != "value" {
		t.Errorf("expected _default/NOPREFIXKEY=value, got %q", buckets["_default"]["NOPREFIXKEY"])
	}
	if _, ok := buckets["DB"]; !ok {
		t.Error("expected DB bucket to exist")
	}
}

func TestSplit_CustomSeparator(t *testing.T) {
	s := splitter.New(splitter.WithSeparator("."))
	env := map[string]string{
		"app.port": "8080",
		"app.host": "0.0.0.0",
		"debug":    "true",
	}
	buckets := s.Split(env)
	if buckets["app"]["port"] != "8080" {
		t.Errorf("expected app/port=8080, got %q", buckets["app"]["port"])
	}
	if buckets["_default"]["debug"] != "true" {
		t.Errorf("expected _default/debug=true, got %q", buckets["_default"]["debug"])
	}
}

func TestBucket_ReturnsCopy(t *testing.T) {
	s := splitter.New()
	env := map[string]string{"DB_HOST": "localhost"}
	s.Split(env)
	b := s.Bucket("DB")
	if b == nil {
		t.Fatal("expected non-nil bucket")
	}
	b["HOST"] = "mutated"
	original := s.Bucket("DB")
	if original["HOST"] == "mutated" {
		t.Error("Bucket should return a copy, not a reference")
	}
}

func TestBucket_MissingReturnsNil(t *testing.T) {
	s := splitter.New()
	s.Split(map[string]string{})
	if s.Bucket("NONEXISTENT") != nil {
		t.Error("expected nil for missing bucket")
	}
}

func TestBuckets_ReturnsAllNames(t *testing.T) {
	s := splitter.New()
	env := map[string]string{
		"DB_HOST":  "localhost",
		"APP_NAME": "test",
		"ORPHAN":   "val",
	}
	s.Split(env)
	names := sortedStrings(s.Buckets())
	expected := []string{"APP", "DB", "_default"}
	if len(names) != len(expected) {
		t.Fatalf("expected %v buckets, got %v", expected, names)
	}
	for i, n := range names {
		if n != expected[i] {
			t.Errorf("bucket[%d]: expected %q got %q", i, expected[i], n)
		}
	}
}
