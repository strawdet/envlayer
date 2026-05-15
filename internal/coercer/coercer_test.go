package coercer_test

import (
	"testing"
	"time"

	"github.com/nicholasgasior/envlayer/internal/coercer"
)

func baseEnv() map[string]string {
	return map[string]string{
		"PORT":    "8080",
		"RATIO":   "3.14",
		"ENABLED": "true",
		"TIMEOUT": "30s",
		"LABEL":   "production",
	}
}

func TestString_Found(t *testing.T) {
	c := coercer.New(baseEnv())
	if got := c.String("LABEL", "default"); got != "production" {
		t.Errorf("expected production, got %q", got)
	}
}

func TestString_Missing_ReturnsDefault(t *testing.T) {
	c := coercer.New(baseEnv())
	if got := c.String("MISSING", "fallback"); got != "fallback" {
		t.Errorf("expected fallback, got %q", got)
	}
}

func TestInt_Valid(t *testing.T) {
	c := coercer.New(baseEnv())
	v, err := c.Int("PORT", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != 8080 {
		t.Errorf("expected 8080, got %d", v)
	}
}

func TestInt_Missing_ReturnsDefault(t *testing.T) {
	c := coercer.New(baseEnv())
	v, err := c.Int("MISSING", 42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != 42 {
		t.Errorf("expected 42, got %d", v)
	}
}

func TestInt_Invalid_ReturnsError(t *testing.T) {
	c := coercer.New(map[string]string{"X": "notanint"})
	_, err := c.Int("X", 0)
	if err == nil {
		t.Fatal("expected error for invalid int")
	}
}

func TestFloat64_Valid(t *testing.T) {
	c := coercer.New(baseEnv())
	v, err := c.Float64("RATIO", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != 3.14 {
		t.Errorf("expected 3.14, got %f", v)
	}
}

func TestFloat64_Invalid_ReturnsError(t *testing.T) {
	c := coercer.New(map[string]string{"X": "abc"})
	_, err := c.Float64("X", 0)
	if err == nil {
		t.Fatal("expected error for invalid float64")
	}
}

func TestBool_TrueVariants(t *testing.T) {
	for _, val := range []string{"true", "1", "yes", "on", "TRUE", "Yes"} {
		c := coercer.New(map[string]string{"FLAG": val})
		v, err := c.Bool("FLAG", false)
		if err != nil || !v {
			t.Errorf("expected true for value %q, got %v err=%v", val, v, err)
		}
	}
}

func TestBool_FalseVariants(t *testing.T) {
	for _, val := range []string{"false", "0", "no", "off", "FALSE"} {
		c := coercer.New(map[string]string{"FLAG": val})
		v, err := c.Bool("FLAG", true)
		if err != nil || v {
			t.Errorf("expected false for value %q, got %v err=%v", val, v, err)
		}
	}
}

func TestBool_Invalid_ReturnsError(t *testing.T) {
	c := coercer.New(map[string]string{"FLAG": "maybe"})
	_, err := c.Bool("FLAG", false)
	if err == nil {
		t.Fatal("expected error for invalid bool")
	}
}

func TestDuration_Valid(t *testing.T) {
	c := coercer.New(baseEnv())
	v, err := c.Duration("TIMEOUT", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != 30*time.Second {
		t.Errorf("expected 30s, got %v", v)
	}
}

func TestDuration_Invalid_ReturnsError(t *testing.T) {
	c := coercer.New(map[string]string{"T": "nope"})
	_, err := c.Duration("T", 0)
	if err == nil {
		t.Fatal("expected error for invalid duration")
	}
}
