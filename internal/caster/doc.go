// Package caster converts raw string values stored in an environment variable
// map into typed Go values (int, bool, float64, []string, string).
//
// Basic usage:
//
//	env := map[string]string{
//		"PORT":  "8080",
//		"DEBUG": "true",
//		"RATIO": "0.75",
//		"TAGS":  "web,api,grpc",
//	}
//
//	c := caster.New(env)
//	port, _  := c.Int("PORT", 3000)
//	debug, _ := c.Bool("DEBUG", false)
//	tags     := c.Strings("TAGS", ",", nil)
//
// Strict mode:
//
//	c := caster.New(env, caster.WithStrict())
//	// Returns an error instead of the default when a value cannot be parsed.
package caster
