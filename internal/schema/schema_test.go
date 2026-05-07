package schema_test

import (
	"testing"

	"github.com/envlayer/envlayer/internal/schema"
)

func TestValidate_RequiredMissing(t *testing.T) {
	s := schema.New([]schema.Field{
		{Key: "DB_HOST", Type: schema.TypeString, Required: true},
	})
	errs := s.Validate(map[string]string{})
	if len(errs) == 0 {
		t.Fatal("expected error for missing required key")
	}
}

func TestValidate_RequiredPresent(t *testing.T) {
	s := schema.New([]schema.Field{
		{Key: "DB_HOST", Type: schema.TypeString, Required: true},
	})
	errs := s.Validate(map[string]string{"DB_HOST": "localhost"})
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
}

func TestValidate_InvalidInt(t *testing.T) {
	s := schema.New([]schema.Field{
		{Key: "PORT", Type: schema.TypeInt},
	})
	errs := s.Validate(map[string]string{"PORT": "not-a-number"})
	if len(errs) == 0 {
		t.Fatal("expected type error for PORT")
	}
}

func TestValidate_ValidBool(t *testing.T) {
	s := schema.New([]schema.Field{
		{Key: "DEBUG", Type: schema.TypeBool},
	})
	for _, val := range []string{"true", "false", "1", "0"} {
		errs := s.Validate(map[string]string{"DEBUG": val})
		if len(errs) != 0 {
			t.Fatalf("unexpected error for bool value %q: %v", val, errs)
		}
	}
}

func TestValidate_InvalidFloat(t *testing.T) {
	s := schema.New([]schema.Field{
		{Key: "RATIO", Type: schema.TypeFloat},
	})
	errs := s.Validate(map[string]string{"RATIO": "abc"})
	if len(errs) == 0 {
		t.Fatal("expected type error for RATIO")
	}
}

func TestApplyDefaults_FillsMissing(t *testing.T) {
	s := schema.New([]schema.Field{
		{Key: "LOG_LEVEL", Type: schema.TypeString, Default: "info"},
	})
	out := s.ApplyDefaults(map[string]string{})
	if out["LOG_LEVEL"] != "info" {
		t.Fatalf("expected default 'info', got %q", out["LOG_LEVEL"])
	}
}

func TestApplyDefaults_DoesNotOverride(t *testing.T) {
	s := schema.New([]schema.Field{
		{Key: "LOG_LEVEL", Type: schema.TypeString, Default: "info"},
	})
	out := s.ApplyDefaults(map[string]string{"LOG_LEVEL": "debug"})
	if out["LOG_LEVEL"] != "debug" {
		t.Fatalf("expected 'debug', got %q", out["LOG_LEVEL"])
	}
}

func TestNew_DefaultsTypeToString(t *testing.T) {
	s := schema.New([]schema.Field{
		{Key: "APP_NAME"},
	})
	errs := s.Validate(map[string]string{"APP_NAME": "myapp"})
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
}
