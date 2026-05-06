// Package differ compares two sets of environment variables and reports
// additions, removals, and modifications between them.
package differ

import "sort"

// ChangeType describes the kind of difference detected.
type ChangeType string

const (
	Added    ChangeType = "added"
	Removed  ChangeType = "removed"
	Modified ChangeType = "modified"
)

// Change represents a single difference between two env maps.
type Change struct {
	Key    string
	Type   ChangeType
	OldVal string
	NewVal string
}

// Differ compares environment variable maps.
type Differ struct{}

// New returns a new Differ instance.
func New() *Differ {
	return &Differ{}
}

// Compare returns the list of changes going from base to next.
func (d *Differ) Compare(base, next map[string]string) []Change {
	var changes []Change

	// Check for removed or modified keys.
	for k, oldVal := range base {
		newVal, exists := next[k]
		if !exists {
			changes = append(changes, Change{Key: k, Type: Removed, OldVal: oldVal})
		} else if newVal != oldVal {
			changes = append(changes, Change{Key: k, Type: Modified, OldVal: oldVal, NewVal: newVal})
		}
	}

	// Check for added keys.
	for k, newVal := range next {
		if _, exists := base[k]; !exists {
			changes = append(changes, Change{Key: k, Type: Added, NewVal: newVal})
		}
	}

	// Sort for deterministic output.
	sort.Slice(changes, func(i, j int) bool {
		if changes[i].Key != changes[j].Key {
			return changes[i].Key < changes[j].Key
		}
		return changes[i].Type < changes[j].Type
	})

	return changes
}

// HasChanges returns true if there is at least one difference.
func (d *Differ) HasChanges(base, next map[string]string) bool {
	return len(d.Compare(base, next)) > 0
}

// Summary returns counts of each change type.
func (d *Differ) Summary(changes []Change) map[ChangeType]int {
	summary := map[ChangeType]int{
		Added:    0,
		Removed:  0,
		Modified: 0,
	}
	for _, c := range changes {
		summary[c.Type]++
	}
	return summary
}
