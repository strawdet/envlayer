package differ_test

import (
	"testing"

	"envlayer/internal/differ"
)

func TestCompare_AddedKeys(t *testing.T) {
	d := differ.New()
	base := map[string]string{"FOO": "bar"}
	next := map[string]string{"FOO": "bar", "NEW_KEY": "value"}

	changes := d.Compare(base, next)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Type != differ.Added || changes[0].Key != "NEW_KEY" {
		t.Errorf("unexpected change: %+v", changes[0])
	}
}

func TestCompare_RemovedKeys(t *testing.T) {
	d := differ.New()
	base := map[string]string{"FOO": "bar", "OLD": "gone"}
	next := map[string]string{"FOO": "bar"}

	changes := d.Compare(base, next)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Type != differ.Removed || changes[0].Key != "OLD" {
		t.Errorf("unexpected change: %+v", changes[0])
	}
	if changes[0].OldVal != "gone" {
		t.Errorf("expected OldVal 'gone', got %q", changes[0].OldVal)
	}
}

func TestCompare_ModifiedKeys(t *testing.T) {
	d := differ.New()
	base := map[string]string{"DB_HOST": "localhost"}
	next := map[string]string{"DB_HOST": "prod.db.internal"}

	changes := d.Compare(base, next)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	c := changes[0]
	if c.Type != differ.Modified || c.OldVal != "localhost" || c.NewVal != "prod.db.internal" {
		t.Errorf("unexpected change: %+v", c)
	}
}

func TestCompare_NoChanges(t *testing.T) {
	d := differ.New()
	env := map[string]string{"A": "1", "B": "2"}
	changes := d.Compare(env, env)
	if len(changes) != 0 {
		t.Errorf("expected no changes, got %d", len(changes))
	}
}

func TestHasChanges(t *testing.T) {
	d := differ.New()
	base := map[string]string{"X": "1"}
	next := map[string]string{"X": "2"}
	if !d.HasChanges(base, next) {
		t.Error("expected HasChanges to return true")
	}
	if d.HasChanges(base, base) {
		t.Error("expected HasChanges to return false for identical maps")
	}
}

func TestSummary(t *testing.T) {
	d := differ.New()
	base := map[string]string{"A": "1", "B": "old", "C": "keep"}
	next := map[string]string{"B": "new", "C": "keep", "D": "added"}

	changes := d.Compare(base, next)
	summary := d.Summary(changes)

	if summary[differ.Added] != 1 {
		t.Errorf("expected 1 added, got %d", summary[differ.Added])
	}
	if summary[differ.Removed] != 1 {
		t.Errorf("expected 1 removed, got %d", summary[differ.Removed])
	}
	if summary[differ.Modified] != 1 {
		t.Errorf("expected 1 modified, got %d", summary[differ.Modified])
	}
}
