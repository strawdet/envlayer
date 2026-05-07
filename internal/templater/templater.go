// Package templater provides simple template rendering for env values,
// allowing values to contain Go-style {{.KEY}} placeholders that are
// substituted from the resolved environment map.
package templater

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"text/template"
)

// Templater renders template expressions embedded in env values.
type Templater struct {
	strict bool
}

// New returns a Templater. When strict is true, referencing an unknown key
// returns an error; otherwise the placeholder is replaced with an empty string.
func New(strict bool) *Templater {
	return &Templater{strict: strict}
}

var placeholderRe = regexp.MustCompile(`\{\{\s*\.\w+\s*\}\}`)

// HasTemplates reports whether s contains any {{.KEY}} expressions.
func HasTemplates(s string) bool {
	return placeholderRe.MatchString(s)
}

// RenderValue renders a single value string against the provided env map.
func (t *Templater) RenderValue(value string, env map[string]string) (string, error) {
	if !HasTemplates(value) {
		return value, nil
	}

	option := "missingkey=zero"
	if t.strict {
		option = "missingkey=error"
	}

	tmpl, err := template.New("").Option(option).Parse(value)
	if err != nil {
		return "", fmt.Errorf("templater: parse error in %q: %w", value, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, toTemplateData(env)); err != nil {
		return "", fmt.Errorf("templater: render error in %q: %w", value, err)
	}
	return buf.String(), nil
}

// RenderAll renders all values in the env map that contain template expressions.
// It returns a new map with rendered values; the original is not modified.
func (t *Templater) RenderAll(env map[string]string) (map[string]string, error) {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = v
	}
	for k, v := range out {
		if !HasTemplates(v) {
			continue
		}
		rendered, err := t.RenderValue(v, out)
		if err != nil {
			return nil, fmt.Errorf("templater: key %s: %w", k, err)
		}
		out[k] = rendered
	}
	return out, nil
}

// toTemplateData converts a map[string]string into a map[string]interface{}
// keyed by the original names so that {{.KEY}} syntax works correctly.
func toTemplateData(env map[string]string) map[string]interface{} {
	data := make(map[string]interface{}, len(env))
	for k, v := range env {
		// template keys must be valid identifiers; replace non-word chars.
		safe := strings.Map(func(r rune) rune {
			if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') ||
				(r >= '0' && r <= '9') || r == '_' {
				return r
			}
			return '_'
		}, k)
		data[safe] = v
	}
	return data
}
