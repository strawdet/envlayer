package expander_test

import (
	"os"
	"testing"

	"github.com/yourorg/envlayer/internal/expander"
)

func TestExpand_SimpleReference(t *testing.T) {
	env := map[string]string{
		"HOME":    "/home/user",
		"CONFIG":  "${HOME}/.config",
	}
	e := expander.New(false)
	result := e.Expand(env)

	if result["CONFIG"] != "/home/user/.config" {
		t.Errorf("expected /home/user/.config, got %s", result["CONFIG"])
	}
}

func TestExpand_DollarSyntax(t *testing.T) {
	env := map[string]string{
		"BASE": "myapp",
		"NAME": "$BASE-service",
	}
	e := expander.New(false)
	result := e.Expand(env)

	if result["NAME"] != "myapp-service" {
		t.Errorf("expected myapp-service, got %s", result["NAME"])
	}
}

func TestExpand_UnresolvedVar_NoFallback(t *testing.T) {
	env := map[string]string{
		"GREETING": "Hello, ${UNKNOWN_VAR}!",
	}
	e := expander.New(false)
	result := e.Expand(env)

	if result["GREETING"] != "Hello, !" {
		t.Errorf("expected 'Hello, !', got %s", result["GREETING"])
	}
}

func TestExpand_FallbackToOS(t *testing.T) {
	os.Setenv("OS_VAR", "from-os")
	defer os.Unsetenv("OS_VAR")

	env := map[string]string{
		"VALUE": "${OS_VAR}-suffix",
	}
	e := expander.New(true)
	result := e.Expand(env)

	if result["VALUE"] != "from-os-suffix" {
		t.Errorf("expected from-os-suffix, got %s", result["VALUE"])
	}
}

func TestExpand_NoReferences(t *testing.T) {
	env := map[string]string{
		"PLAIN": "just-a-value",
	}
	e := expander.New(false)
	result := e.Expand(env)

	if result["PLAIN"] != "just-a-value" {
		t.Errorf("expected just-a-value, got %s", result["PLAIN"])
	}
}

func TestExpandValue_Convenience(t *testing.T) {
	env := map[string]string{"PORT": "8080"}
	out := expander.ExpandValue("http://localhost:${PORT}", env)
	if out != "http://localhost:8080" {
		t.Errorf("expected http://localhost:8080, got %s", out)
	}
}

func TestHasReferences(t *testing.T) {
	if !expander.HasReferences("${FOO}") {
		t.Error("expected HasReferences to return true for ${FOO}")
	}
	if expander.HasReferences("plain-value") {
		t.Error("expected HasReferences to return false for plain-value")
	}
}
