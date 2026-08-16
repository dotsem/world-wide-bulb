package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// StateResponse defines the JSON response payload for the bulb state endpoint.
type StateResponse struct {
	State      bool  `json:"state"`
	CooldownMs int64 `json:"cooldown_ms"`
	Viewers    int   `json:"viewers"`
}

// GetState returns the current state of the bulb.
func (h *Handler) GetState(c *gin.Context) {
	state := h.engine.GetState()

	primaryKey, _ := h.getOrCreateDeviceID(c)
	remaining := h.engine.GetRemainingCooldown(h.hasher.Hash(primaryKey))

	c.JSON(http.StatusOK, StateResponse{
		State:      state,
		CooldownMs: remaining.Milliseconds(),
		Viewers:    h.hub.ClientCount(),
	})
}
