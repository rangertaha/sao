package cot

import "time"

// Event is the normalized CoT envelope used inside SAO.
type Event struct {
	UID    string
	Type   string
	Time   time.Time
	RawXML []byte
}
