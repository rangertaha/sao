package cot

import (
	"encoding/xml"
	"fmt"
	"strings"
	"time"
)

// ParseEvent parses CoT XML into a normalized Event.
func ParseEvent(data []byte) (Event, error) {
	var payload struct {
		XMLName xml.Name `xml:"event"`
		UID     string   `xml:"uid,attr"`
		Type    string   `xml:"type,attr"`
		Time    string   `xml:"time,attr"`
	}

	if err := xml.Unmarshal(data, &payload); err != nil {
		return Event{}, fmt.Errorf("decode cot xml: %w", err)
	}
	if payload.XMLName.Local != "event" {
		return Event{}, fmt.Errorf("root element must be <event>")
	}
	if strings.TrimSpace(payload.UID) == "" {
		return Event{}, fmt.Errorf("cot uid is required")
	}
	if strings.TrimSpace(payload.Type) == "" {
		return Event{}, fmt.Errorf("cot type is required")
	}

	var eventTime time.Time
	if strings.TrimSpace(payload.Time) != "" {
		parsed, err := time.Parse(time.RFC3339, payload.Time)
		if err != nil {
			return Event{}, fmt.Errorf("parse cot time: %w", err)
		}
		eventTime = parsed
	}

	return Event{
		UID:    payload.UID,
		Type:   payload.Type,
		Time:   eventTime,
		RawXML: append([]byte(nil), data...),
	}, nil
}
