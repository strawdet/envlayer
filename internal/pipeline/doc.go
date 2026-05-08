// Package pipeline provides a composable, ordered processing pipeline for
// environment variable maps in envlayer.
//
// A Pipeline chains together one or more Steps — each a function that
// receives a map[string]string and returns a transformed map or an error.
// Built-in helpers integrate with the transformer, expander, masker, and
// validator packages:
//
//	p := pipeline.New(
//		pipeline.WithTransformer(tr),   // key/value case & prefix filters
//		pipeline.WithExpander(e),        // variable reference expansion
//		pipeline.WithMasker(m),          // sensitive-value masking
//		pipeline.WithValidator(v),        // schema / required-key validation
//	)
//
//	out, err := p.Run(env)
//
// Custom steps can be added with WithStep for any additional transformation
// logic not covered by the built-in integrations.
package pipeline
