package ws

import (
	"time"
	"world-wide-bulb/internal/store"
)

// StateMessage represents a state change message payload for WebSocket broadcasting.
type StateMessage struct {
	Type      string    `json:"type"`
	ID        int64     `json:"id"`
	State     bool      `json:"state"`
	Reason    string    `json:"reason,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// FromToggle maps a database store.Toggle record into a StateMessage DTO.
func FromToggle(t store.Toggle) StateMessage {
	return StateMessage{
		Type:      "state_changed",
		ID:        t.ID,
		State:     t.State,
		Reason:    t.Reason.String,
		CreatedAt: t.CreatedAt.Time,
	}
}

// ReasonMessage represents a reason update message payload for WebSocket broadcasting.
type ReasonMessage struct {
	Type   string `json:"type"`
	ID     int64  `json:"id"`
	Reason string `json:"reason"`
}

// FromReason maps a database store.Toggle record into a ReasonMessage DTO.
func FromReason(t store.Toggle) ReasonMessage {
	return ReasonMessage{
		Type:   "reason_updated",
		ID:     t.ID,
		Reason: t.Reason.String,
	}
}
