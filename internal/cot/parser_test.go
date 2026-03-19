package cot

import "testing"

func TestParseEvent(t *testing.T) {
	t.Parallel()

	raw := []byte(`<event uid="A-1" type="a-f-G-U-C" time="2026-03-18T12:00:00Z"></event>`)
	evt, err := ParseEvent(raw)
	if err != nil {
		t.Fatalf("ParseEvent() error: %v", err)
	}
	if evt.UID != "A-1" {
		t.Fatalf("unexpected uid: %s", evt.UID)
	}
	if evt.Type != "a-f-G-U-C" {
		t.Fatalf("unexpected type: %s", evt.Type)
	}
}

func TestParseEventMissingUID(t *testing.T) {
	t.Parallel()

	_, err := ParseEvent([]byte(`<event type="a-f-G-U-C"></event>`))
	if err == nil {
		t.Fatal("expected error for missing uid")
	}
}
