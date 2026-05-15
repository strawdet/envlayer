package cloner_test

import (
	"testing"

	"github.com/yourorg/envlayer/internal/cloner"
)

func TestClone_DeepCopy(t *testing.T) {
	original := map[string]string{"A": "1", "B": "2"}
	c := cloner.New()
	got := c.Clone(original)
	original["A"] = "mutated"
	if got["A"] != "1" {
		t.Errorf("expected deep copy, got mutated value")
	}
}

func TestClone_PrefixFilter(t *testing.T) {
	env := map[string]string{"APP_HOST": "localhost", "APP_PORT": "8080", "DB_URL": "postgres://"}
	c := cloner.New(cloner.WithPrefixFilter("APP_"))
	got := c.Clone(env)
	if len(got) != 2 {
		t.Errorf("expected 2 keys, got %d", len(got))
	}
	if _, ok := got["DB_URL"]; ok {
		t.Error("DB_URL should have been filtered out")
	}
}

func TestClone_StripPrefix(t *testing.T) {
	env := map[string]string{"APP_HOST": "localhost", "APP_PORT": "8080"}
	c := cloner.New(cloner.WithPrefixFilter("APP_"), cloner.WithStripPrefix())
	got := c.Clone(env)
	if got["HOST"] != "localhost" {
		t.Errorf("expected HOST=localhost, got %q", got["HOST"])
	}
	if got["PORT"] != "8080" {
		t.Errorf("expected PORT=8080, got %q", got["PORT"])
	}
}

func TestClone_OmitEmpty(t *testing.T) {
	env := map[string]string{"A": "value", "B": "", "C": "other"}
	c := cloner.New(cloner.WithOmitEmpty())
	got := c.Clone(env)
	if _, ok := got["B"]; ok {
		t.Error("expected empty key B to be omitted")
	}
	if len(got) != 2 {
		t.Errorf("expected 2 keys, got %d", len(got))
	}
}

func TestMerge_SrcOverridesDst(t *testing.T) {
	dst := map[string]string{"A": "old", "B": "keep"}
	src := map[string]string{"A": "new", "C": "added"}
	c := cloner.New()
	got := c.Merge(dst, src)
	if got["A"] != "new" {
		t.Errorf("expected A=new, got %q", got["A"])
	}
	if got["B"] != "keep" {
		t.Errorf("expected B=keep, got %q", got["B"])
	}
	if got["C"] != "added" {
		t.Errorf("expected C=added, got %q", got["C"])
	}
}

func TestMerge_DoesNotMutateDst(t *testing.T) {
	dst := map[string]string{"A": "original"}
	src := map[string]string{"A": "override"}
	c := cloner.New()
	c.Merge(dst, src)
	if dst["A"] != "original" {
		t.Error("Merge must not mutate dst")
	}
}
