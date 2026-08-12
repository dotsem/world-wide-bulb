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

	assert.True(t, msg.State)
	assert.Equal(t, "test reason", msg.Reason)
	assert.Equal(t, now, msg.CreatedAt)
}
