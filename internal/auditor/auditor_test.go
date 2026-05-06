package auditor_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/your-org/envlayer/internal/auditor"
)

func TestRecord_AndEvents(t *testing.T) {
	a := auditor.New(nil)
	a.Record(auditor.EventLoaded, "APP_ENV", "source=base.env")
	a.Record(auditor.EventOverride, "APP_ENV", "old=dev new=prod")

	events := a.Events()
	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}
	if events[0].Kind != auditor.EventLoaded {
		t.Errorf("expected LOADED, got %s", events[0].Kind)
	}
	if events[1].Key != "APP_ENV" {
		t.Errorf("expected key APP_ENV, got %s", events[1].Key)
	}
}

func TestFlush_WritesAllEvents(t *testing.T) {
	var buf bytes.Buffer
	a := auditor.New(&buf)
	a.Record(auditor.EventMissing, "DB_URL", "not found in any layer")
	a.Record(auditor.EventExpanded, "BASE_URL", "expanded from $HOST")

	if err := a.Flush(); err != nil {
		t.Fatalf("Flush error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "MISSING") {
		t.Error("expected MISSING in output")
	}
	if !strings.Contains(out, "EXPANDED") {
		t.Error("expected EXPANDED in output")
	}
	if !strings.Contains(out, "DB_URL") {
		t.Error("expected DB_URL in output")
	}
}

func TestReset_ClearsEvents(t *testing.T) {
	a := auditor.New(nil)
	a.Record(auditor.EventLoaded, "FOO", "")
	a.Reset()

	if len(a.Events()) != 0 {
		t.Error("expected events to be cleared after Reset")
	}
}

func TestEvent_String_Format(t *testing.T) {
	var buf bytes.Buffer
	a := auditor.New(&buf)
	a.Record(auditor.EventOverride, "MY_KEY", "old=a new=b")
	_ = a.Flush()

	line := buf.String()
	if !strings.Contains(line, "[OVERRIDE]") {
		t.Errorf("unexpected format: %s", line)
	}
	if !strings.Contains(line, `key="MY_KEY"`) {
		t.Errorf("expected key in output: %s", line)
	}
}

func TestEvents_ReturnsCopy(t *testing.T) {
	a := auditor.New(nil)
	a.Record(auditor.EventLoaded, "X", "")
	events := a.Events()
	events[0].Key = "MUTATED"

	original := a.Events()
	if original[0].Key == "MUTATED" {
		t.Error("Events() should return a copy, not a reference")
	}
}
