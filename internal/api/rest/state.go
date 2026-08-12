package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type StateResponse struct {
	State      bool  `json:"state"`
	CooldownMs int64 `json:"cooldown_ms"`
}

// GetState returns the current state of the bulb.
func (h *Handler) GetState(c *gin.Context) {
	state := h.engine.GetState()

	primaryKey, _ := h.getOrCreateDeviceID(c)
	remaining := h.engine.GetRemainingCooldown(h.hasher.Hash(primaryKey))

	c.JSON(http.StatusOK, StateResponse{
		State:      state,
		CooldownMs: remaining.Milliseconds(),
	})
}
