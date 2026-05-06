// Package auditor provides a simple audit log for tracking environment
// variable resolution events such as overrides, missing keys, and expansions.
package auditor

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// EventKind classifies the type of audit event.
type EventKind string

const (
	EventOverride  EventKind = "OVERRIDE"
	EventMissing   EventKind = "MISSING"
	EventExpanded  EventKind = "EXPANDED"
	EventLoaded    EventKind = "LOADED"
)

// Event represents a single audit record.
type Event struct {
	Timestamp time.Time
	Kind      EventKind
	Key       string
	Detail    string
}

func (e Event) String() string {
	return fmt.Sprintf("%s [%s] key=%q %s",
		e.Timestamp.Format(time.RFC3339), e.Kind, e.Key, e.Detail)
}

// Auditor collects audit events and can write them to a writer.
type Auditor struct {
	events []Event
	out    io.Writer
}

// New creates a new Auditor. If w is nil, os.Stdout is used.
func New(w io.Writer) *Auditor {
	if w == nil {
		w = os.Stdout
	}
	return &Auditor{out: w}
}

// Record appends an event to the audit log.
func (a *Auditor) Record(kind EventKind, key, detail string) {
	a.events = append(a.events, Event{
		Timestamp: time.Now(),
		Kind:      kind,
		Key:       key,
		Detail:    detail,
	})
}

// Events returns a copy of all recorded events.
func (a *Auditor) Events() []Event {
	out := make([]Event, len(a.events))
	copy(out, a.events)
	return out
}

// Flush writes all recorded events to the configured writer.
func (a *Auditor) Flush() error {
	var sb strings.Builder
	for _, e := range a.events {
		sb.WriteString(e.String())
		sb.WriteByte('\n')
	}
	_, err := fmt.Fprint(a.out, sb.String())
	return err
}

// Reset clears all recorded events.
func (a *Auditor) Reset() {
	a.events = a.events[:0]
}
