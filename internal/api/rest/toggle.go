package rest

import (
	"encoding/json"
	"errors"
	"io"
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

// ToggleRequest defines the expected JSON payload for toggling the bulb.
type ToggleRequest struct {
	Reason string `json:"reason" binding:"max=60"`
}

// ToggleResponse defines the response payload for toggling the bulb.
type ToggleResponse struct {
	State      bool      `json:"state"`
	Reason     string    `json:"reason,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	CooldownMs int64     `json:"cooldown_ms"`
}

// PostToggle flips the state of the bulb and broadcasts the change.
func (h *Handler) PostToggle(c *gin.Context) {
	var req ToggleRequest
	if err := c.ShouldBindJSON(&req); err != nil && err != io.EOF {
		c.JSON(http.StatusBadRequest, gin.H{errKey: "invalid payload (reason max 60 chars)"})
		return
	}

	primaryKey, secondaryKey := h.getOrCreateDeviceID(c)

	toggle, remainingTime, err := h.engine.Toggle(c.Request.Context(), req.Reason, h.hasher.Hash(primaryKey))
	if err != nil {
		if errors.Is(err, bulb.ErrCooldown) {
			c.JSON(http.StatusTooManyRequests, gin.H{errKey: err.Error()})
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
		State:      toggle.State,
		Reason:     toggle.Reason.String,
		CreatedAt:  toggle.CreatedAt.Time,
		CooldownMs: remainingTime.Milliseconds(),
	})
}
