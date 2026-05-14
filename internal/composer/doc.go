// Package composer provides a high-level orchestration layer for envlayer.
//
// It wires together the resolver, expander, validator, transformer, and
// exporter packages into a single, declarative Run call, removing the
// boilerplate of manually chaining each stage.
//
// # Basic usage
//
//	c := composer.New(composer.Options{
//		BaseDir:      "/app/config",
//		Context:      "staging",
//		RequiredKeys: []string{"DATABASE_URL", "SECRET_KEY"},
//		KeyCase:      "upper",
//		OutputFormat: "dotenv",
//		OutputPath:   ".env.merged",
//	})
//	result, err := c.Run()
//
// The returned Result.Env map contains the final, fully-processed
// environment variables. Non-fatal validation findings are surfaced
// through Result.Warnings rather than causing an error.
package composer
