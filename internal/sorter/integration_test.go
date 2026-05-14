package sorter_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/envlayer/envlayer/internal/loader"
	"github.com/envlayer/envlayer/internal/sorter"
)

func writeSorterEnv(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestSorter_WithLoader_SortedOutput(t *testing.T) {
	dir := t.TempDir()
	path := writeSorterEnv(t, dir, ".env", "ZEBRA=stripes\nAPPLE=fruit\nMONKEY=banana\n")

	env, err := loader.LoadFile(path)
	if err != nil {
		t.Fatalf("LoadFile: %v", err)
	}

	s := sorter.New(sorter.WithOrder(sorter.Ascending))
	keys := s.SortedKeys(env)

	if len(keys) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(keys))
	}
	expected := []string{"APPLE", "MONKEY", "ZEBRA"}
	for i, k := range keys {
		if k != expected[i] {
			t.Errorf("index %d: got %q, want %q", i, k, expected[i])
		}
	}
}

func TestSorter_WithLoader_ByValue(t *testing.T) {
	dir := t.TempDir()
	path := writeSorterEnv(t, dir, ".env", "Z=alpha\nA=omega\nM=beta\n")

	env, err := loader.LoadFile(path)
	if err != nil {
		t.Fatalf("LoadFile: %v", err)
	}

	s := sorter.New(sorter.WithSortByValue())
	pairs := s.Sorted(env)

	// expected value order: alpha, beta, omega
	wantValues := []string{"alpha", "beta", "omega"}
	for i, p := range pairs {
		if p[1] != wantValues[i] {
			t.Errorf("index %d: got value %q, want %q", i, p[1], wantValues[i])
		}
	}
}
