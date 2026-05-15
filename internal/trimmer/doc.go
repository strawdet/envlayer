// Package trimmer provides a Trimmer that normalizes environment variable
// keys and values by removing leading/trailing whitespace and optionally
// collapsing internal whitespace runs.
//
// # Usage
//
//	tr := trimmer.New(
//		trimmer.WithTrimKeys(),
//		trimmer.WithCollapseWhitespace(),
//		trimmer.WithSkipKeys("PRIVATE_KEY", "CERT"),
//	)
//	clean := tr.Apply(resolved)
//
// By default, Trimmer trims leading and trailing whitespace from values.
// Key trimming and whitespace collapsing must be explicitly enabled.
// Keys listed via WithSkipKeys are passed through unchanged.
package trimmer
