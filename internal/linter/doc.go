// Package linter provides static analysis for environment variable maps
// loaded by envlayer.
//
// It checks for common issues such as:
//   - Empty values that may indicate misconfiguration
//   - Keys that shadow existing OS environment variables
//   - Duplicate keys across multiple .env layers
//
// Example usage:
//
//	l := linter.New(
//		linter.WithEmptyValueCheck(),
//		linter.WithOSShadowCheck(),
//		linter.WithDuplicateCheck(),
//	)
//
//	findings := l.Lint(env)
//	for _, f := range findings {
//		fmt.Println(f)
//	}
//
//	layerFindings := l.LintLayers(layers)
package linter
