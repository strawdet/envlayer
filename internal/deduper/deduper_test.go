package deduper_test

import (
	"sort"
	"testing"

	"github.com/user/envlayer/internal/deduper"
)

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func TestMerge_KeepLast_Default(t *testing.T) {
	d := deduper.New()
	a := map[string]string{"FOO": "a", "BAR": "1"}
	b := map[string]string{"FOO": "b", "BAZ": "2"}
	out := d.Merge(a, b)
	if out["FOO"] != "b" {
		t.Errorf("expected FOO=b, got %s", out["FOO"])
	}
	if out["BAR"] != "1" {
		t.Errorf("expected BAR=1, got %s", out["BAR"])
	}
	if out["BAZ"] != "2" {
		t.Errorf("expected BAZ=2, got %s", out["BAZ"])
	}
}

func TestMerge_KeepFirst(t *testing.T) {
	d := deduper.New(deduper.WithStrategy(deduper.KeepFirst))
	a := map[string]string{"FOO": "original"}
	b := map[string]string{"FOO": "override"}
	out := d.Merge(a, b)
	if out["FOO"] != "original" {
		t.Errorf("expected FOO=original, got %s", out["FOO"])
	}
}

func TestMerge_NoOverlap(t *testing.T) {
	d := deduper.New()
	a := map[string]string{"A": "1"}
	b := map[string]string{"B": "2"}
	out := d.Merge(a, b)
	if len(out) != 2 {
		t.Errorf("expected 2 keys, got %d", len(out))
	}
}

func TestApply_ReturnsCopy(t *testing.T) {
	d := deduper.New()
	env := map[string]string{"X": "1", "Y": "2"}
	out := d.Apply(env)
	out["Z"] = "3"
	if _, ok := env["Z"]; ok {
		t.Error("Apply should return a copy, original was mutated")
	}
}

func TestDuplicates_DetectsOverlap(t *testing.T) {
	a := map[string]string{"FOO": "1", "BAR": "2"}
	b := map[string]string{"FOO": "3", "BAZ": "4"}
	c := map[string]string{"BAR": "5", "QUX": "6"}
	dups := deduper.Duplicates(a, b, c)
	sort.Strings(dups)
	if len(dups) != 2 || dups[0] != "BAR" || dups[1] != "FOO" {
		t.Errorf("unexpected duplicates: %v", dups)
	}
}

func TestDuplicates_NoneFound(t *testing.T) {
	a := map[string]string{"A": "1"}
	b := map[string]string{"B": "2"}
	dups := deduper.Duplicates(a, b)
	if len(dups) != 0 {
		t.Errorf("expected no duplicates, got %v", dups)
	}
}

func TestMerge_EmptyLayers(t *testing.T) {
	d := deduper.New()
	out := d.Merge()
	if len(out) != 0 {
		t.Errorf("expected empty map, got %v", out)
	}
}
