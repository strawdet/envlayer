package tagger_test

import (
	"testing"

	"github.com/your/envlayer/internal/tagger"
)

func TestTag_AddAndRetrieve(t *testing.T) {
	tgr := tagger.New()
	tgr.Tag("DB_PASSWORD", "secret", "infra")

	tags := tgr.Tags("DB_PASSWORD")
	if len(tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(tags))
	}
	if tags[0] != "infra" || tags[1] != "secret" {
		t.Errorf("unexpected tags order: %v", tags)
	}
}

func TestTag_NoDuplicates(t *testing.T) {
	tgr := tagger.New()
	tgr.Tag("KEY", "app", "app", "secret")

	if got := len(tgr.Tags("KEY")); got != 2 {
		t.Errorf("expected 2 unique tags, got %d", got)
	}
}

func TestUntag_RemovesTag(t *testing.T) {
	tgr := tagger.New()
	tgr.Tag("API_KEY", "secret", "app")
	tgr.Untag("API_KEY", "secret")

	if tgr.HasTag("API_KEY", "secret") {
		t.Error("expected 'secret' tag to be removed")
	}
	if !tgr.HasTag("API_KEY", "app") {
		t.Error("expected 'app' tag to remain")
	}
}

func TestHasTag_True(t *testing.T) {
	tgr := tagger.New()
	tgr.Tag("PORT", "infra")
	if !tgr.HasTag("PORT", "infra") {
		t.Error("expected HasTag to return true")
	}
}

func TestHasTag_False(t *testing.T) {
	tgr := tagger.New()
	if tgr.HasTag("UNKNOWN", "secret") {
		t.Error("expected HasTag to return false for unknown key")
	}
}

func TestFilterByTag(t *testing.T) {
	tgr := tagger.New()
	tgr.Tag("DB_PASS", "secret")
	tgr.Tag("APP_NAME", "app")
	tgr.Tag("API_TOKEN", "secret", "app")

	env := map[string]string{
		"DB_PASS":  "hunter2",
		"APP_NAME": "envlayer",
		"API_TOKEN": "tok123",
	}

	secrets := tgr.FilterByTag(env, "secret")
	if len(secrets) != 2 {
		t.Fatalf("expected 2 secret keys, got %d", len(secrets))
	}
	if _, ok := secrets["DB_PASS"]; !ok {
		t.Error("expected DB_PASS in secrets")
	}
	if _, ok := secrets["API_TOKEN"]; !ok {
		t.Error("expected API_TOKEN in secrets")
	}
}

func TestAllTags(t *testing.T) {
	tgr := tagger.New()
	tgr.Tag("A", "secret", "infra")
	tgr.Tag("B", "app", "infra")

	all := tgr.AllTags()
	if len(all) != 3 {
		t.Fatalf("expected 3 distinct tags, got %d: %v", len(all), all)
	}
	if all[0] != "app" || all[1] != "infra" || all[2] != "secret" {
		t.Errorf("unexpected AllTags order: %v", all)
	}
}

func TestWithTags_Option(t *testing.T) {
	pre := map[string][]string{
		"HOST": {"infra"},
		"TOKEN": {"secret"},
	}
	tgr := tagger.New(tagger.WithTags(pre))

	if !tgr.HasTag("HOST", "infra") {
		t.Error("expected HOST to have 'infra' tag from WithTags")
	}
	if !tgr.HasTag("TOKEN", "secret") {
		t.Error("expected TOKEN to have 'secret' tag from WithTags")
	}
}
