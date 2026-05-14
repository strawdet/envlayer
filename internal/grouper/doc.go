// Package grouper partitions a flat environment variable map into named
// groups based on key prefixes and a configurable separator.
//
// # Overview
//
// Given an env map such as:
//
//	DB_HOST=localhost
//	DB_PORT=5432
//	APP_ENV=production
//
// Grouper produces:
//
//	{
//	  "DB":  {"HOST": "localhost", "PORT": "5432"},
//	  "APP": {"ENV":  "production"},
//	}
//
// Keys that contain no separator are placed under the empty-string group.
//
// # Usage
//
//	g := grouper.New(grouper.WithSeparator("_"))
//	groups := g.Group(env)
//	flat   := g.Flatten(groups)
package grouper
