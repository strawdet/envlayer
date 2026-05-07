// Package schema provides a declarative schema system for environment variables.
//
// A Schema is defined as a collection of Field definitions, each specifying:
//   - Key: the environment variable name
//   - Type: expected value type (string, int, bool, float)
//   - Required: whether the key must be present
//   - Default: fallback value applied when the key is absent
//   - Description: human-readable documentation for the field
//
// Example usage:
//
//	s := schema.New([]schema.Field{
//		{Key: "PORT",      Type: schema.TypeInt,    Required: true},
//		{Key: "LOG_LEVEL", Type: schema.TypeString,  Default: "info"},
//		{Key: "DEBUG",     Type: schema.TypeBool,    Default: "false"},
//	})
//
//	env := s.ApplyDefaults(rawEnv)
//	if errs := s.Validate(env); len(errs) > 0 {
//		// handle validation failures
//	}
package schema
