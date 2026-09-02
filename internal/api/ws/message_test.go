package ws_test

import (
	"database/sql"
	"testing"
	"time"
	"world-wide-bulb/internal/api/ws"
	"world-wide-bulb/internal/store"

	"github.com/stretchr/testify/assert"
)

func TestFromToggle(t *testing.T) {
	now := time.Now()
	toggle := store.Toggle{
		ID:    42,
		State: true,
		Reason: sql.NullString{
			String: "test reason",
			Valid:  true,
		},
		CreatedAt: sql.NullTime{
			Time:  now,
			Valid: true,
		},
	}

	msg := ws.FromToggle(toggle)

	assert.Equal(t, "state_changed", msg.Type)
	assert.Equal(t, int64(42), msg.ID)
	assert.True(t, msg.State)
	assert.Equal(t, "test reason", msg.Reason)
	assert.Equal(t, now, msg.CreatedAt)
}

func TestFromReason(t *testing.T) {
	toggle := store.Toggle{
		ID: 42,
		Reason: sql.NullString{
			String: "test reason",
			Valid:  true,
		},
	}

	msg := ws.FromReason(toggle)

	assert.Equal(t, "reason_updated", msg.Type)
	assert.Equal(t, int64(42), msg.ID)
	assert.Equal(t, "test reason", msg.Reason)
}
