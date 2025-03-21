package events

import "encoding/json"

type InternalEvent struct {
	EventType string          `json:"event_type"`
	Data      json.RawMessage `json:"data"`
}
