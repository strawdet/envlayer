// Package composer provides a high-level orchestration layer that chains
// multiple envlayer stages (resolve → expand → validate → transform → export)
// into a single, configurable pipeline run.
package composer

import (
	"fmt"

	"github.com/envlayer/envlayer/internal/expander"
	"github.com/envlayer/envlayer/internal/exporter"
	"github.com/envlayer/envlayer/internal/resolver"
	"github.com/envlayer/envlayer/internal/transformer"
	"github.com/envlayer/envlayer/internal/validator"
)

// Options controls which stages are active and how they behave.
type Options struct {
	// BaseDir is the root directory containing .env layer files.
	BaseDir string
	// Context selects the active environment (e.g. "production").
	Context string
	// RequiredKeys lists keys that must be present after resolution.
	RequiredKeys []string
	// StripPrefix removes a common prefix from all keys after transformation.
	StripPrefix string
	// KeyCase normalises key casing: "upper", "lower", or "" for no change.
	KeyCase string
	// OutputFormat is passed to the exporter: "dotenv", "export", or "json".
	OutputFormat string
	// OutputPath writes the result to a file when non-empty.
	OutputPath string
}

// Result holds the final merged environment and any non-fatal warnings.
type Result struct {
	Env      map[string]string
	Warnings []string
}

// Composer orchestrates the full envlayer processing chain.
type Composer struct {
	opts Options
}

// New creates a Composer with the supplied options.
func New(opts Options) *Composer {
	return &Composer{opts: opts}
}

// Run executes every active stage in order and returns the final Result.
func (c *Composer) Run() (*Result, error) {
	// 1. Resolve layers.
	res := resolver.New(c.opts.BaseDir)
	env, err := res.Resolve(c.opts.Context)
	if err != nil {
		return nil, fmt.Errorf("composer: resolve: %w", err)
	}

	// 2. Expand variable references.
	exp := expander.New(env)
	env = exp.ExpandAll(env)

	// 3. Validate required keys.
	var warnings []string
	if len(c.opts.RequiredKeys) > 0 {
		v := validator.New(c.opts.RequiredKeys)
		findings := v.Validate(env)
		for _, f := range findings {
			if f.Fatal {
				return nil, fmt.Errorf("composer: validation: %s", f.Message)
			}
			warnings = append(warnings, f.Message)
		}
	}

	// 4. Transform keys.
	var txOpts []transformer.Option
	if c.opts.StripPrefix != "" {
		txOpts = append(txOpts, transformer.WithStripPrefix(c.opts.StripPrefix))
	}
	if c.opts.KeyCase != "" {
		txOpts = append(txOpts, transformer.WithKeyCase(c.opts.KeyCase))
	}
	if len(txOpts) > 0 {
		tx := transformer.New(txOpts...)
		env = tx.Apply(env)
	}

	// 5. Export (optional).
	if c.opts.OutputPath != "" {
		exp2 := exporter.New()
		if err := exp2.WriteToFile(env, c.opts.OutputFormat, c.opts.OutputPath); err != nil {
			return nil, fmt.Errorf("composer: export: %w", err)
		}
	}

	return &Result{Env: env, Warnings: warnings}, nil
}
