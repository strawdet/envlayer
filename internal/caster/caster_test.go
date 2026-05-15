package caster_test

import (
	"testing"

	"github.com/yourorg/envlayer/internal/caster"
)

func TestString_Found(t *testing.T) {
	c := caster.New(map[string]string{"HOST": "localhost"})
	if got := c.String("HOST", "default"); got != "localhost" {
		t.Fatalf("want localhost, got %q", got)
	}
}

func TestString_Missing(t *testing.T) {
	c := caster.New(map[string]string{})
	if got := c.String("HOST", "fallback"); got != "fallback" {
		t.Fatalf("want fallback, got %q", got)
	}
}

func TestInt_Valid(t *testing.T) {
	c := caster.New(map[string]string{"PORT": "8080"})
	v, err := c.Int("PORT", 0)
	if err != nil {
		t.Fatal(err)
	}
	if v != 8080 {
		t.Fatalf("want 8080, got %d", v)
	}
}

func TestInt_Missing_ReturnsDefault(t *testing.T) {
	c := caster.New(map[string]string{})
	v, err := c.Int("PORT", 3000)
	if err != nil || v != 3000 {
		t.Fatalf("want 3000/nil, got %d/%v", v, err)
	}
}

func TestInt_Invalid_LenientMode(t *testing.T) {
	c := caster.New(map[string]string{"PORT": "abc"})
	v, err := c.Int("PORT", 42)
	if err != nil || v != 42 {
		t.Fatalf("want 42/nil in lenient mode, got %d/%v", v, err)
	}
}

func TestInt_Invalid_StrictMode(t *testing.T) {
	c := caster.New(map[string]string{"PORT": "abc"}, caster.WithStrict())
	_, err := c.Int("PORT", 0)
	if err == nil {
		t.Fatal("expected error in strict mode")
	}
}

func TestBool_Valid(t *testing.T) {
	c := caster.New(map[string]string{"DEBUG": "true"})
	v, err := c.Bool("DEBUG", false)
	if err != nil || !v {
		t.Fatalf("want true/nil, got %v/%v", v, err)
	}
}

func TestBool_Invalid_StrictMode(t *testing.T) {
	c := caster.New(map[string]string{"DEBUG": "yes_please"}, caster.WithStrict())
	_, err := c.Bool("DEBUG", false)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestFloat64_Valid(t *testing.T) {
	c := caster.New(map[string]string{"RATIO": "3.14"})
	v, err := c.Float64("RATIO", 0)
	if err != nil || v != 3.14 {
		t.Fatalf("want 3.14/nil, got %f/%v", v, err)
	}
}

func TestStrings_Split(t *testing.T) {
	c := caster.New(map[string]string{"TAGS": "a,b,c"})
	got := c.Strings("TAGS", ",", nil)
	if len(got) != 3 || got[1] != "b" {
		t.Fatalf("unexpected result: %v", got)
	}
}

func TestStrings_Missing_ReturnsDefault(t *testing.T) {
	c := caster.New(map[string]string{})
	def := []string{"x"}
	got := c.Strings("TAGS", ",", def)
	if len(got) != 1 || got[0] != "x" {
		t.Fatalf("want default, got %v", got)
	}
}
