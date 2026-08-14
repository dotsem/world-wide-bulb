package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
	"world-wide-bulb/internal/api/ws"
	"world-wide-bulb/internal/bulb"

	"github.com/gin-gonic/gin"
)

const (
	errKey            = "error"
	errInternalServer = "internal server error"
)

// ToggleResponse defines the response payload for toggling the bulb.
type ToggleResponse struct {
	ID         string    `json:"id"`
	State      bool      `json:"state"`
	CreatedAt  time.Time `json:"created_at"`
	CooldownMs int64     `json:"cooldown_ms"`
}

// PostToggle flips the state of the bulb and broadcasts the change.
func (h *Handler) PostToggle(c *gin.Context) {
	primaryKey, secondaryKey := h.getOrCreateDeviceID(c)

	toggle, remainingTime, err := h.engine.Toggle(c.Request.Context(), h.hasher.Hash(primaryKey))
	if err != nil {
		if errors.Is(err, bulb.ErrCooldown) {
			remaining := h.engine.GetRemainingCooldown(h.hasher.Hash(primaryKey))
			c.JSON(http.StatusTooManyRequests, gin.H{
				errKey:        err.Error(),
				"cooldown_ms": remaining.Milliseconds(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			errKey: errInternalServer,
		})
		return
	}

	if secondaryKey != "" {
		remainingTime = h.engine.RecordCooldown(h.hasher.Hash(secondaryKey))
	}

	wsMsg := ws.FromToggle(toggle)

	if payload, err := json.Marshal(wsMsg); err == nil {
		h.hub.Broadcast(payload)
	}

	c.JSON(http.StatusOK, ToggleResponse{
		ID:         toggle.Uuid,
		State:      toggle.State,
		CreatedAt:  toggle.CreatedAt.Time,
		CooldownMs: remainingTime.Milliseconds(),
	})
}
