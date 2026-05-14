// Package sorter provides deterministic ordering of environment variable maps.
//
// By default, keys are sorted in ascending lexicographic order. Options allow
// descending order, sorting by value instead of key, and case-insensitive
// comparisons.
//
// Example:
//
//	s := sorter.New(
//		sorter.WithOrder(sorter.Ascending),
//		sorter.WithCaseInsensitive(),
//	)
//	keys := s.SortedKeys(env)
//	for _, k := range keys {
//		fmt.Printf("%s=%s\n", k, env[k])
//	}
package sorter
