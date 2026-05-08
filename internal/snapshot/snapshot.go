// Package snapshot provides functionality to capture and restore
// environment variable states at a point in time.
package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Snapshot holds a captured state of environment variables.
type Snapshot struct {
	Label     string            `json:"label"`
	CapturedAt time.Time        `json:"captured_at"`
	Env       map[string]string `json:"env"`
}

// Manager manages saving and loading snapshots.
type Manager struct {
	dir string
}

// New creates a new snapshot Manager that stores snapshots in dir.
func New(dir string) (*Manager, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("snapshot: create dir: %w", err)
	}
	return &Manager{dir: dir}, nil
}

// Save captures the given env map under the provided label.
func (m *Manager) Save(label string, env map[string]string) (*Snapshot, error) {
	snap := &Snapshot{
		Label:      label,
		CapturedAt: time.Now().UTC(),
		Env:        env,
	}
	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("snapshot: marshal: %w", err)
	}
	path := m.pathFor(label)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return nil, fmt.Errorf("snapshot: write file: %w", err)
	}
	return snap, nil
}

// Load retrieves a previously saved snapshot by label.
func (m *Manager) Load(label string) (*Snapshot, error) {
	data, err := os.ReadFile(m.pathFor(label))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("snapshot: %q not found", label)
		}
		return nil, fmt.Errorf("snapshot: read file: %w", err)
	}
	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, fmt.Errorf("snapshot: unmarshal: %w", err)
	}
	return &snap, nil
}

// Delete removes a snapshot by label.
func (m *Manager) Delete(label string) error {
	if err := os.Remove(m.pathFor(label)); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("snapshot: %q not found", label)
		}
		return fmt.Errorf("snapshot: delete: %w", err)
	}
	return nil
}

// List returns the labels of all snapshots stored in the manager's directory.
func (m *Manager) List() ([]string, error) {
	entries, err := os.ReadDir(m.dir)
	if err != nil {
		return nil, fmt.Errorf("snapshot: list: %w", err)
	}
	var labels []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		const suffix = ".snap.json"
		if len(name) > len(suffix) && name[len(name)-len(suffix):] == suffix {
			labels = append(labels, name[:len(name)-len(suffix)])
		}
	}
	return labels, nil
}

func (m *Manager) pathFor(label string) string {
	return fmt.Sprintf("%s/%s.snap.json", m.dir, label)
}
