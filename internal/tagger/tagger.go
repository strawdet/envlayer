// Package tagger provides tagging and annotation support for environment variables.
// Tags are arbitrary string labels attached to keys, enabling grouping, filtering,
// and documentation of env vars by category (e.g. "secret", "infra", "app").
package tagger

import "sort"

// Tagger manages tags associated with environment variable keys.
type Tagger struct {
	tags map[string][]string // key -> list of tags
}

// Option configures a Tagger.
type Option func(*Tagger)

// WithTags pre-loads a tag mapping.
func WithTags(tags map[string][]string) Option {
	return func(t *Tagger) {
		for k, v := range tags {
			copy := make([]string, len(v))
			_ = copy[:len(v)]
			for i, tag := range v {
				copy[i] = tag
			}
			t.tags[k] = copy
		}
	}
}

// New creates a new Tagger with the given options.
func New(opts ...Option) *Tagger {
	t := &Tagger{tags: make(map[string][]string)}
	for _, o := range opts {
		o(t)
	}
	return t
}

// Tag adds one or more tags to the given key.
func (t *Tagger) Tag(key string, tags ...string) {
	existing := t.tags[key]
	seen := make(map[string]struct{}, len(existing))
	for _, tg := range existing {
		seen[tg] = struct{}{}
	}
	for _, tg := range tags {
		if _, ok := seen[tg]; !ok {
			existing = append(existing, tg)
			seen[tg] = struct{}{}
		}
	}
	t.tags[key] = existing
}

// Untag removes a specific tag from a key. No-op if tag is not present.
func (t *Tagger) Untag(key, tag string) {
	tags := t.tags[key]
	result := tags[:0]
	for _, tg := range tags {
		if tg != tag {
			result = append(result, tg)
		}
	}
	t.tags[key] = result
}

// Tags returns the tags for the given key, sorted alphabetically.
func (t *Tagger) Tags(key string) []string {
	tags := t.tags[key]
	out := make([]string, len(tags))
	copy(out, tags)
	sort.Strings(out)
	return out
}

// HasTag reports whether the given key has the specified tag.
func (t *Tagger) HasTag(key, tag string) bool {
	for _, tg := range t.tags[key] {
		if tg == tag {
			return true
		}
	}
	return false
}

// FilterByTag returns all keys from env that have the given tag.
func (t *Tagger) FilterByTag(env map[string]string, tag string) map[string]string {
	out := make(map[string]string)
	for k, v := range env {
		if t.HasTag(k, tag) {
			out[k] = v
		}
	}
	return out
}

// AllTags returns a deduplicated, sorted list of all tags across all keys.
func (t *Tagger) AllTags() []string {
	seen := make(map[string]struct{})
	for _, tags := range t.tags {
		for _, tg := range tags {
			seen[tg] = struct{}{}
		}
	}
	out := make([]string, 0, len(seen))
	for tg := range seen {
		out = append(out, tg)
	}
	sort.Strings(out)
	return out
}
