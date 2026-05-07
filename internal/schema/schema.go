// Package schema provides validation of environment variable definitions
// against a declared schema, supporting types, defaults, and descriptions.
package schema

import (
	"fmt"
	"strconv"
	"strings"
)

// FieldType represents the expected type of an env var value.
type FieldType string

const (
	TypeString  FieldType = "string"
	TypeInt     FieldType = "int"
	TypeBool    FieldType = "bool"
	TypeFloat   FieldType = "float"
)

// Field describes a single environment variable declaration.
type Field struct {
	Key         string
	Type        FieldType
	Required    bool
	Default     string
	Description string
}

// Schema holds a collection of declared fields.
type Schema struct {
	fields map[string]Field
}

// New creates a new Schema from a slice of Field definitions.
func New(fields []Field) *Schema {
	m := make(map[string]Field, len(fields))
	for _, f := range fields {
		if f.Type == "" {
			f.Type = TypeString
		}
		m[f.Key] = f
	}
	return &Schema{fields: m}
}

// Validate checks the provided env map against the schema.
// It returns a list of violation messages (empty means valid).
func (s *Schema) Validate(env map[string]string) []string {
	var errs []string
	for key, field := range s.fields {
		val, ok := env[key]
		if !ok || strings.TrimSpace(val) == "" {
			if field.Required && field.Default == "" {
				errs = append(errs, fmt.Sprintf("required key %q is missing", key))
			}
			continue
		}
		if err := validateType(key, val, field.Type); err != nil {
			errs = append(errs, err.Error())
		}
	}
	return errs
}

// ApplyDefaults fills in default values for missing keys defined in the schema.
func (s *Schema) ApplyDefaults(env map[string]string) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = v
	}
	for key, field := range s.fields {
		if _, ok := out[key]; !ok && field.Default != "" {
			out[key] = field.Default
		}
	}
	return out
}

func validateType(key, val string, t FieldType) error {
	switch t {
	case TypeInt:
		if _, err := strconv.Atoi(val); err != nil {
			return fmt.Errorf("key %q expects int, got %q", key, val)
		}
	case TypeBool:
		if _, err := strconv.ParseBool(val); err != nil {
			return fmt.Errorf("key %q expects bool, got %q", key, val)
		}
	case TypeFloat:
		if _, err := strconv.ParseFloat(val, 64); err != nil {
			return fmt.Errorf("key %q expects float, got %q", key, val)
		}
	}
	return nil
}
