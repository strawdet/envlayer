// Package cli provides the command-line interface for envlayer.
// It wires together the core internal packages into a usable tool
// that resolves, merges, validates, and exports environment variables
// based on runtime context.
package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/envlayer/envlayer/internal/exporter"
	"github.com/envlayer/envlayer/internal/linter"
	"github.com/envlayer/envlayer/internal/masker"
	"github.com/envlayer/envlayer/internal/merger"
	"github.com/envlayer/envlayer/internal/printer"
	"github.com/envlayer/envlayer/internal/resolver"
	"github.com/envlayer/envlayer/internal/transformer"
)

// Config holds the parsed CLI configuration for a single invocation.
type Config struct {
	// BaseDir is the directory containing .env layer files.
	BaseDir string

	// Context is the runtime context (e.g. "production", "staging").
	Context string

	// OutputFormat controls how the resolved env is written (dotenv, export, json).
	OutputFormat string

	// OutputFile is an optional path to write the output. Defaults to stdout.
	OutputFile string

	// MaskSensitive controls whether sensitive keys are masked in printed output.
	MaskSensitive bool

	// PrintTable renders the resolved env in a table instead of key=value.
	PrintTable bool

	// LintOnly runs linting checks without producing output.
	LintOnly bool

	// StripPrefix removes the given prefix from all resolved keys.
	StripPrefix string

	// KeyCase normalises key casing: "upper", "lower", or "" (no change).
	KeyCase string
}

// Run executes the envlayer resolution pipeline according to cfg.
// It returns a non-nil error if any required step fails.
func Run(cfg Config) error {
	if cfg.BaseDir == "" {
		cfg.BaseDir = "."
	}
	if cfg.OutputFormat == "" {
		cfg.OutputFormat = "dotenv"
	}

	// Build merger and resolver.
	m := merger.New(cfg.BaseDir)
	r := resolver.New(m)

	env, err := r.Resolve(cfg.Context)
	if err != nil {
		return fmt.Errorf("resolve: %w", err)
	}

	// Optional linting pass.
	if cfg.LintOnly {
		return runLint(cfg, m)
	}

	// Apply transformer options when requested.
	if cfg.StripPrefix != "" || cfg.KeyCase != "" {
		var opts []transformer.Option
		if cfg.StripPrefix != "" {
			opts = append(opts, transformer.WithStripPrefix(cfg.StripPrefix))
		}
		if cfg.KeyCase != "" {
			opts = append(opts, transformer.WithKeyCase(cfg.KeyCase))
		}
		t := transformer.New(opts...)
		env, err = t.Apply(env)
		if err != nil {
			return fmt.Errorf("transform: %w", err)
		}
	}

	// Print to stdout when no output file is requested.
	if cfg.OutputFile == "" {
		return printEnv(cfg, env)
	}

	// Write to file via exporter.
	e := exporter.New()
	format := exporter.Format(strings.ToLower(cfg.OutputFormat))
	return e.WriteToFile(env, format, cfg.OutputFile)
}

// runLint performs lint checks on the resolved layers and prints findings.
func runLint(cfg Config, m *merger.Merger) error {
	layers, err := m.LayersForContext(cfg.Context)
	if err != nil {
		return fmt.Errorf("layers: %w", err)
	}

	l := linter.New(
		linter.WithDuplicateCheck(),
		linter.WithOSShadowCheck(),
		linter.WithEmptyValueCheck(),
	)

	findings := l.LintLayers(layers)
	if len(findings) == 0 {
		fmt.Fprintln(os.Stdout, "No lint findings.")
		return nil
	}

	for _, f := range findings {
		fmt.Fprintln(os.Stdout, f)
	}
	return nil
}

// printEnv renders the resolved environment map to stdout.
func printEnv(cfg Config, env map[string]string) error {
	var opts []printer.Option
	if cfg.MaskSensitive {
		m := masker.New()
		opts = append(opts, printer.WithMasker(m))
	}

	p := printer.New(opts...)

	format := printer.KVFormat
	switch strings.ToLower(cfg.OutputFormat) {
	case "table":
		format = printer.TableFormat
	case "json":
		format = printer.JSONFormat
	}

	return p.Print(os.Stdout, env, format)
}
