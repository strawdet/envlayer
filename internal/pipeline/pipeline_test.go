package pipeline_test

import (
	"strings"
	"testing"

	"github.com/envlayer/envlayer/internal/expander"
	"github.com/envlayer/envlayer/internal/masker"
	"github.com/envlayer/envlayer/internal/pipeline"
	"github.com/envlayer/envlayer/internal/transformer"
)

func TestPipeline_Empty(t *testing.T) {
	p := pipeline.New()
	env := map[string]string{"FOO": "bar"}
	out, err := p.Run(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["FOO"] != "bar" {
		t.Errorf("expected bar, got %s", out["FOO"])
	}
}

func TestPipeline_TransformerStep(t *testing.T) {
	tr := transformer.New(
		transformer.WithKeyCase("upper"),
	)
	p := pipeline.New(pipeline.WithTransformer(tr))
	env := map[string]string{"foo": "hello", "bar": "world"}
	out, err := p.Run(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["FOO"]; !ok {
		t.Error("expected key FOO after upper-case transform")
	}
	if _, ok := out["foo"]; ok {
		t.Error("expected original key foo to be gone")
	}
}

func TestPipeline_ExpanderStep(t *testing.T) {
	e := expander.New()
	p := pipeline.New(pipeline.WithExpander(e))
	env := map[string]string{
		"BASE": "/app",
		"PATH": "${BASE}/bin",
	}
	out, err := p.Run(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["PATH"] != "/app/bin" {
		t.Errorf("expected /app/bin, got %s", out["PATH"])
	}
}

func TestPipeline_MaskerStep(t *testing.T) {
	m := masker.New(masker.WithSensitiveKeys([]string{"SECRET"}))
	p := pipeline.New(pipeline.WithMasker(m))
	env := map[string]string{"SECRET": "topsecret", "APP": "envlayer"}
	out, err := p.Run(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["SECRET"] == "topsecret" {
		t.Error("expected SECRET to be masked")
	}
	if out["APP"] != "envlayer" {
		t.Errorf("expected APP to be unchanged, got %s", out["APP"])
	}
}

func TestPipeline_CustomStep(t *testing.T) {
	prefixStep := pipeline.Step(func(env map[string]string) (map[string]string, error) {
		out := make(map[string]string, len(env))
		for k, v := range env {
			out[k] = "prefix_" + v
		}
		return out, nil
	})
	p := pipeline.New(pipeline.WithStep(prefixStep))
	env := map[string]string{"KEY": "value"}
	out, err := p.Run(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(out["KEY"], "prefix_") {
		t.Errorf("expected prefix_, got %s", out["KEY"])
	}
}

func TestPipeline_StepError_Propagates(t *testing.T) {
	failStep := pipeline.Step(func(env map[string]string) (map[string]string, error) {
		return nil, fmt.Errorf("intentional failure")
	})
	p := pipeline.New(pipeline.WithStep(failStep))
	_, err := p.Run(map[string]string{"K": "v"})
	if err == nil {
		t.Fatal("expected error from failing step")
	}
	if !strings.Contains(err.Error(), "pipeline step 0") {
		t.Errorf("unexpected error message: %v", err)
	}
}
