package freezer_test

import (
	"errors"
	"testing"

	"envlayer/internal/freezer"
)

func TestSet_AndGet(t *testing.T) {
	f := freezer.New()
	if err := f.Set("KEY", "value"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v, ok := f.Get("KEY")
	if !ok || v != "value" {
		t.Errorf("expected value=%q ok=true, got %q %v", "value", v, ok)
	}
}

func TestFreeze_PreventsSet(t *testing.T) {
	f := freezer.New()
	_ = f.Set("A", "1")
	f.Freeze()

	if !f.IsFrozen() {
		t.Fatal("expected IsFrozen to return true")
	}

	err := f.Set("B", "2")
	if !errors.Is(err, freezer.ErrFrozen) {
		t.Errorf("expected ErrFrozen, got %v", err)
	}
}

func TestLoad_PopulatesEnv(t *testing.T) {
	f := freezer.New()
	src := map[string]string{"X": "10", "Y": "20"}
	if err := f.Load(src); err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	for k, want := range src {
		if got, ok := f.Get(k); !ok || got != want {
			t.Errorf("key %q: expected %q, got %q", k, want, got)
		}
	}
}

func TestLoad_FrozenReturnsError(t *testing.T) {
	f := freezer.New()
	f.Freeze()
	err := f.Load(map[string]string{"Z": "99"})
	if !errors.Is(err, freezer.ErrFrozen) {
		t.Errorf("expected ErrFrozen, got %v", err)
	}
}

func TestSnapshot_ReturnsCopy(t *testing.T) {
	f := freezer.New()
	_ = f.Load(map[string]string{"A": "1"})
	f.Freeze()

	snap := f.Snapshot()
	snap["A"] = "mutated"

	v, _ := f.Get("A")
	if v != "1" {
		t.Errorf("snapshot mutation affected original: got %q", v)
	}
}

func TestDiff_DetectsChanges(t *testing.T) {
	f := freezer.New()
	_ = f.Load(map[string]string{"A": "1", "B": "2", "C": "3"})
	f.Freeze()

	newEnv := map[string]string{"A": "1", "B": "changed", "D": "new"}
	diff := f.Diff(newEnv)

	findings := map[string]bool{}
	for _, line := range diff {
		findings[line] = true
	}

	if !findings["+ D=new"] {
		t.Error("expected added key D in diff")
	}
	if !findings["- C=3"] {
		t.Error("expected removed key C in diff")
	}
	found := false
	for _, l := range diff {
		if l == `~ B: "2" -> "changed"` {
			found = true
		}
	}
	if !found {
		t.Errorf("expected modified key B in diff, got: %v", diff)
	}
}

func TestDiff_NoChanges(t *testing.T) {
	f := freezer.New()
	_ = f.Load(map[string]string{"A": "1"})
	f.Freeze()

	diff := f.Diff(map[string]string{"A": "1"})
	if len(diff) != 0 {
		t.Errorf("expected empty diff, got: %v", diff)
	}
}
