// Package printer provides formatted terminal and file output for resolved
// environment variable maps.
//
// Supported formats:
//
//   - table  — aligned two-column table with KEY / VALUE headers
//   - kv     — plain KEY=VALUE lines, suitable for shell sourcing
//   - json   — indented JSON object keyed by variable name
//
// An optional Masker function can be supplied via WithMasker to redact
// sensitive values (e.g. passwords, API keys) before they reach the output
// writer. The masker receives both the key and the raw value so decisions
// can be made based on naming conventions.
//
// Example usage:
//
//	p := printer.New(printer.WithMasker(myMasker))
//	if err := p.Print(os.Stdout, env, printer.FormatTable); err != nil {
//		log.Fatal(err)
//	}
package printer
