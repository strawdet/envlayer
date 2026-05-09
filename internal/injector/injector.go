// Package injector provides functionality to inject resolved environment
// variables into a child process's environment at runtime.
package injector

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Injector merges a resolved env map into a subprocess environment.
type Injector struct {
	baseEnv  []string
	override bool
}

// Option configures an Injector.
type Option func(*Injector)

// WithOverride causes resolved keys to override existing OS env variables.
// By default, existing OS variables take precedence.
func WithOverride() Option {
	return func(i *Injector) {
		i.override = true
	}
}

// WithBaseEnv sets the base environment instead of inheriting os.Environ.
func WithBaseEnv(env []string) Option {
	return func(i *Injector) {
		i.baseEnv = env
	}
}

// New creates a new Injector with the provided options.
func New(opts ...Option) *Injector {
	inj := &Injector{
		baseEnv: os.Environ(),
	}
	for _, o := range opts {
		o(inj)
	}
	return inj
}

// BuildEnv merges resolved into the base environment and returns the result.
func (inj *Injector) BuildEnv(resolved map[string]string) []string {
	existing := make(map[string]string, len(inj.baseEnv))
	for _, kv := range inj.baseEnv {
		parts := strings.SplitN(kv, "=", 2)
		if len(parts) == 2 {
			existing[parts[0]] = parts[1]
		}
	}

	for k, v := range resolved {
		if _, exists := existing[k]; !exists || inj.override {
			existing[k] = v
		}
	}

	result := make([]string, 0, len(existing))
	for k, v := range existing {
		result = append(result, fmt.Sprintf("%s=%s", k, v))
	}
	return result
}

// Run executes the given command with the merged environment.
func (inj *Injector) Run(resolved map[string]string, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Env = inj.BuildEnv(resolved)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
