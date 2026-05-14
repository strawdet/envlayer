// Package tagger attaches arbitrary string labels (tags) to environment variable
// keys, enabling downstream components to filter, group, or annotate variables
// by category.
//
// # Overview
//
// A Tagger maintains a mapping from env key names to a set of string tags.
// Tags are user-defined and can represent any semantic category, such as:
//
//   - "secret"  — values that should be masked or encrypted
//   - "infra"   — infrastructure-level variables (hosts, ports, regions)
//   - "app"     — application-level configuration
//   - "debug"   — variables only relevant in development
//
// # Usage
//
//	tgr := tagger.New()
//	tgr.Tag("DB_PASSWORD", "secret", "infra")
//	tgr.Tag("APP_NAME", "app")
//
//	secrets := tgr.FilterByTag(env, "secret")
//
// Tags can also be pre-loaded via the WithTags option:
//
//	tgr := tagger.New(tagger.WithTags(map[string][]string{
//		"API_KEY": {"secret"},
//	}))
package tagger
