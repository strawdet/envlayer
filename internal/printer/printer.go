// Package printer provides formatted output rendering for resolved environment variables.
// It supports table, key=value, and JSON output formats with optional masking.
package printer

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
	"text/tabwriter"
)

// Format defines the output format for printing.
type Format string

const (
	FormatTable  Format = "table"
	FormatKV     Format = "kv"
	FormatJSON   Format = "json"
)

// Masker is a function that masks a value for a given key.
type Masker func(key, value string) string

// Printer renders environment variables to an io.Writer.
type Printer struct {
	masker Masker
}

// Option configures a Printer.
type Option func(*Printer)

// WithMasker sets a masking function applied to each value before output.
func WithMasker(m Masker) Option {
	return func(p *Printer) {
		p.masker = m
	}
}

// New creates a new Printer with the given options.
func New(opts ...Option) *Printer {
	p := &Printer{}
	for _, o := range opts {
		o(p)
	}
	return p
}

// Print writes the environment map to w in the specified format.
func (p *Printer) Print(w io.Writer, env map[string]string, format Format) error {
	keys := sortedKeys(env)
	switch format {
	case FormatTable:
		return p.printTable(w, env, keys)
	case FormatKV:
		return p.printKV(w, env, keys)
	case FormatJSON:
		return p.printJSON(w, env, keys)
	default:
		return fmt.Errorf("printer: unsupported format %q", format)
	}
}

func (p *Printer) mask(key, value string) string {
	if p.masker != nil {
		return p.masker(key, value)
	}
	return value
}

func (p *Printer) printTable(w io.Writer, env map[string]string, keys []string) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "KEY\tVALUE")
	fmt.Fprintln(tw, strings.Repeat("-", 20)+"\t"+strings.Repeat("-", 30))
	for _, k := range keys {
		fmt.Fprintf(tw, "%s\t%s\n", k, p.mask(k, env[k]))
	}
	return tw.Flush()
}

func (p *Printer) printKV(w io.Writer, env map[string]string, keys []string) error {
	for _, k := range keys {
		if _, err := fmt.Fprintf(w, "%s=%s\n", k, p.mask(k, env[k])); err != nil {
			return err
		}
	}
	return nil
}

func (p *Printer) printJSON(w io.Writer, env map[string]string, keys []string) error {
	masked := make(map[string]string, len(keys))
	for _, k := range keys {
		masked[k] = p.mask(k, env[k])
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(masked)
}

func sortedKeys(env map[string]string) []string {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
