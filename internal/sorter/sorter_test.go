package sorter_test

import (
	"testing"

	"github.com/envlayer/envlayer/internal/sorter"
)

func TestSortedKeys_Ascending(t *testing.T) {
	env := map[string]string{"ZEBRA": "1", "APPLE": "2", "MANGO": "3"}
	s := sorter.New()
	keys := s.SortedKeys(env)
	want := []string{"APPLE", "MANGO", "ZEBRA"}
	for i, k := range keys {
		if k != want[i] {
			t.Errorf("pos %d: got %q, want %q", i, k, want[i])
		}
	}
}

func TestSortedKeys_Descending(t *testing.T) {
	env := map[string]string{"ZEBRA": "1", "APPLE": "2", "MANGO": "3"}
	s := sorter.New(sorter.WithOrder(sorter.Descending))
	keys := s.SortedKeys(env)
	want := []string{"ZEBRA", "MANGO", "APPLE"}
	for i, k := range keys {
		if k != want[i] {
			t.Errorf("pos %d: got %q, want %q", i, k, want[i])
		}
	}
}

func TestSortedKeys_ByValue(t *testing.T) {
	env := map[string]string{"A": "zoo", "B": "ant", "C": "moo"}
	s := sorter.New(sorter.WithSortByValue())
	keys := s.SortedKeys(env)
	want := []string{"B", "C", "A"} // ant, moo, zoo
	for i, k := range keys {
		if k != want[i] {
			t.Errorf("pos %d: got %q, want %q", i, k, want[i])
		}
	}
}

func TestSortedKeys_CaseInsensitive(t *testing.T) {
	env := map[string]string{"beta": "1", "ALPHA": "2", "Gamma": "3"}
	s := sorter.New(sorter.WithCaseInsensitive())
	keys := s.SortedKeys(env)
	want := []string{"ALPHA", "beta", "Gamma"}
	for i, k := range keys {
		if k != want[i] {
			t.Errorf("pos %d: got %q, want %q", i, k, want[i])
		}
	}
}

func TestSorted_ReturnsPairs(t *testing.T) {
	env := map[string]string{"B": "two", "A": "one"}
	s := sorter.New()
	pairs := s.Sorted(env)
	if len(pairs) != 2 {
		t.Fatalf("expected 2 pairs, got %d", len(pairs))
	}
	if pairs[0][0] != "A" || pairs[0][1] != "one" {
		t.Errorf("unexpected first pair: %v", pairs[0])
	}
	if pairs[1][0] != "B" || pairs[1][1] != "two" {
		t.Errorf("unexpected second pair: %v", pairs[1])
	}
}

func TestSortedKeys_EmptyMap(t *testing.T) {
	s := sorter.New()
	keys := s.SortedKeys(map[string]string{})
	if len(keys) != 0 {
		t.Errorf("expected empty slice, got %v", keys)
	}
}
