package exporter

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

// Format represents the output format for exported variables.
type Format string

const (
	FormatDotEnv Format = "dotenv"
	FormatExport Format = "export"
	FormatJSON   Format = "json"
)

// Exporter writes merged environment variables to a target.
type Exporter struct {
	vars   map[string]string
	format Format
}

// New creates a new Exporter with the given variables and format.
func New(vars map[string]string, format Format) *Exporter {
	return &Exporter{vars: vars, format: format}
}

// Write outputs the environment variables to the given writer.
func (e *Exporter) Write(w io.Writer) error {
	keys := sortedKeys(e.vars)
	switch e.format {
	case FormatDotEnv:
		return writeDotEnv(w, keys, e.vars)
	case FormatExport:
		return writeExport(w, keys, e.vars)
	case FormatJSON:
		return writeJSON(w, keys, e.vars)
	default:
		return fmt.Errorf("unsupported format: %s", e.format)
	}
}

// WriteToFile writes the environment variables to a file at the given path.
func (e *Exporter) WriteToFile(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("exporter: create file %q: %w", path, err)
	}
	defer f.Close()
	return e.Write(f)
}

func writeDotEnv(w io.Writer, keys []string, vars map[string]string) error {
	for _, k := range keys {
		if _, err := fmt.Fprintf(w, "%s=%q\n", k, vars[k]); err != nil {
			return err
		}
	}
	return nil
}

func writeExport(w io.Writer, keys []string, vars map[string]string) error {
	for _, k := range keys {
		if _, err := fmt.Fprintf(w, "export %s=%q\n", k, vars[k]); err != nil {
			return err
		}
	}
	return nil
}

func writeJSON(w io.Writer, keys []string, vars map[string]string) error {
	var sb strings.Builder
	sb.WriteString("{\n")
	for i, k := range keys {
		comma := ","
		if i == len(keys)-1 {
			comma = ""
		}
		sb.WriteString(fmt.Sprintf("  %q: %q%s\n", k, vars[k], comma))
	}
	sb.WriteString("}\n")
	_, err := fmt.Fprint(w, sb.String())
	return err
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
