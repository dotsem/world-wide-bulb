package ws

import (
	"time"
	"world-wide-bulb/internal/store"
)

type StateMessage struct {
	State     bool      `json:"state"`
	Reason    string    `json:"reason,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

func FromToggle(t store.Toggle) StateMessage {
	return StateMessage{
		State:     t.State,
		Reason:    t.Reason.String,
		CreatedAt: t.CreatedAt.Time,
	}
}
