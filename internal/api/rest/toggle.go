package rest

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
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

// PostToggle flips the state of the bulb and broadcasts the change.
func (h *Handler) PostToggle(c *gin.Context) {
	var req ToggleRequest
	if err := c.ShouldBindJSON(&req); err != nil && err != io.EOF {
		c.JSON(http.StatusBadRequest, gin.H{errKey: "invalid payload (reason max 60 chars)"})
		return
	}

	toggle, err := h.engine.Toggle(c.Request.Context(), req.Reason, h.hasher.Hash(c.ClientIP()))
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

	res := ws.FromToggle(toggle)

	if payload, err := json.Marshal(res); err == nil {
		h.hub.Broadcast(payload)
	}

	c.JSON(http.StatusOK, res)
}
