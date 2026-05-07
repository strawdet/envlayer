// Package transformer provides composable transformations for environment
// variable maps produced by the envlayer resolver.
//
// Transformations are applied in order:
//  1. Prefix filtering – retain only keys that match a given prefix.
//  2. Prefix stripping – remove the matched prefix from retained keys.
//  3. Key case normalisation – convert all keys to upper or lower case.
//  4. Value case normalisation – convert all values to upper or lower case.
//
// Example:
//
//	tr := transformer.New(
//		transformer.WithPrefixFilter("APP_"),
//		transformer.WithStripPrefix(true),
//		transformer.WithKeyCase("lower"),
//	)
//	result := tr.Apply(resolvedEnv)
package transformer
