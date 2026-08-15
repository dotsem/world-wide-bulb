package ws

import (
	"time"
	"world-wide-bulb/internal/store"
)

// StateMessage represents a state change message payload for WebSocket broadcasting.
type StateMessage struct {
	ID        int64     `json:"id"`
	State     bool      `json:"state"`
	Reason    string    `json:"reason,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// FromToggle maps a database store.Toggle record into a StateMessage DTO.
func FromToggle(t store.Toggle) StateMessage {
	return StateMessage{
		ID:        t.ID,
		State:     t.State,
		Reason:    t.Reason.String,
		CreatedAt: t.CreatedAt.Time,
	}
}

