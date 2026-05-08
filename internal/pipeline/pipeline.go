// Package pipeline provides a composable processing pipeline for environment
// variable maps, chaining transformer, expander, masker, and validator steps.
package pipeline

import (
	"fmt"

	"github.com/envlayer/envlayer/internal/expander"
	"github.com/envlayer/envlayer/internal/masker"
	"github.com/envlayer/envlayer/internal/transformer"
	"github.com/envlayer/envlayer/internal/validator"
)

// Step is a function that transforms an env map, returning the result or an error.
type Step func(env map[string]string) (map[string]string, error)

// Pipeline executes an ordered sequence of Steps over an env map.
type Pipeline struct {
	steps []Step
}

// Option configures a Pipeline.
type Option func(*Pipeline)

// New creates a Pipeline with the given options.
func New(opts ...Option) *Pipeline {
	p := &Pipeline{}
	for _, o := range opts {
		o(p)
	}
	return p
}

// WithTransformer appends a transformer step.
func WithTransformer(t *transformer.Transformer) Option {
	return func(p *Pipeline) {
		p.steps = append(p.steps, func(env map[string]string) (map[string]string, error) {
			return t.Apply(env), nil
		})
	}
}

// WithExpander appends a variable-expansion step.
func WithExpander(e *expander.Expander) Option {
	return func(p *Pipeline) {
		p.steps = append(p.steps, func(env map[string]string) (map[string]string, error) {
			return e.ExpandAll(env), nil
		})
	}
}

// WithMasker appends a masking step (values are masked in the output map).
func WithMasker(m *masker.Masker) Option {
	return func(p *Pipeline) {
		p.steps = append(p.steps, func(env map[string]string) (map[string]string, error) {
			out := make(map[string]string, len(env))
			for k, v := range env {
				out[k] = m.Mask(k, v)
			}
			return out, nil
		})
	}
}

// WithValidator appends a validation step; returns an error if validation fails.
func WithValidator(v *validator.Validator) Option {
	return func(p *Pipeline) {
		p.steps = append(p.steps, func(env map[string]string) (map[string]string, error) {
			result := v.Validate(env)
			if len(result.Errors) > 0 {
				return nil, fmt.Errorf("pipeline validation failed: %v", result.Errors)
			}
			return env, nil
		})
	}
}

// WithStep appends a custom Step to the pipeline.
func WithStep(s Step) Option {
	return func(p *Pipeline) {
		p.steps = append(p.steps, s)
	}
}

// Run executes all steps in order, threading the env map through each one.
func (p *Pipeline) Run(env map[string]string) (map[string]string, error) {
	current := env
	for i, step := range p.steps {
		var err error
		current, err = step(current)
		if err != nil {
			return nil, fmt.Errorf("pipeline step %d: %w", i, err)
		}
	}
	return current, nil
}
