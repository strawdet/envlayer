// Package profiler manages named environment profiles, allowing users to
// define and switch between sets of context layers (e.g. "dev", "staging").
package profiler

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Profile represents a named set of context layers.
type Profile struct {
	Name    string   `json:"name"`
	Context string   `json:"context"`
	Layers  []string `json:"layers"`
}

// Profiler loads and stores profiles from a JSON config file.
type Profiler struct {
	configPath string
	profiles   map[string]Profile
}

// New creates a new Profiler backed by the given config file path.
func New(configPath string) (*Profiler, error) {
	p := &Profiler{
		configPath: configPath,
		profiles:   make(map[string]Profile),
	}
	if err := p.load(); err != nil {
		return nil, err
	}
	return p, nil
}

// load reads profiles from the JSON config file.
func (p *Profiler) load() error {
	data, err := os.ReadFile(p.configPath)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("profiler: read config: %w", err)
	}
	var list []Profile
	if err := json.Unmarshal(data, &list); err != nil {
		return fmt.Errorf("profiler: parse config: %w", err)
	}
	for _, prof := range list {
		p.profiles[prof.Name] = prof
	}
	return nil
}

// Get returns a profile by name.
func (p *Profiler) Get(name string) (Profile, bool) {
	prof, ok := p.profiles[name]
	return prof, ok
}

// List returns all profile names.
func (p *Profiler) List() []string {
	names := make([]string, 0, len(p.profiles))
	for name := range p.profiles {
		names = append(names, name)
	}
	return names
}

// Save persists the current profiles to the config file.
func (p *Profiler) Save(prof Profile) error {
	p.profiles[prof.Name] = prof
	list := make([]Profile, 0, len(p.profiles))
	for _, v := range p.profiles {
		list = append(list, v)
	}
	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return fmt.Errorf("profiler: marshal: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(p.configPath), 0o755); err != nil {
		return fmt.Errorf("profiler: mkdir: %w", err)
	}
	return os.WriteFile(p.configPath, data, 0o644)
}
