// Package scoper provides namespace-based scoping for environment variables.
//
// It allows callers to:
//   - Extract a subset of keys belonging to a named namespace (e.g. "APP_").
//   - Prefix a flat map of keys with a namespace to produce namespaced output.
//   - Discover which namespaces are present in a merged environment map.
//
// Example usage:
//
//	s := scoper.New()
//	appEnv := s.Scope(merged, "APP")   // {"HOST": "localhost", "PORT": "8080"}
//	prefixed := s.Namespace(flat, "SVC") // {"SVC_HOST": "...", "SVC_PORT": "..."}
//	ns := s.Namespaces(merged)           // ["APP", "DB", "SVC"]
package scoper
