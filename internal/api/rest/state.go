package rest

import (
	"net/http"
	"world-wide-bulb/internal/api/middleware"

	"github.com/gin-gonic/gin"
)

// StateResponse defines the JSON response payload for the bulb state endpoint.
type StateResponse struct {
	State      bool  `json:"state"`
	CooldownMs int64 `json:"cooldown_ms"`
	Viewers    int   `json:"viewers"`
	// TurnstileSiteKey is set only when Turnstile is enabled and this client has not passed it yet.
	TurnstileSiteKey string `json:"turnstile_site_key,omitempty"`
}

// GetState returns the current state of the bulb.
func (h *Handler) GetState(c *gin.Context) {
	state := h.engine.GetState()

	primaryKey, _ := h.getOrCreateDeviceID(c)
	remaining := h.engine.GetRemainingCooldown(h.hasher.Hash(primaryKey))

	res := StateResponse{
		State:      state,
		CooldownMs: remaining.Milliseconds(),
		Viewers:    h.hub.ClientCount(),
	}
	if h.turnstileSiteKey != "" && middleware.HumanID(c, h.hasher) == "" {
		res.TurnstileSiteKey = h.turnstileSiteKey
	}
	c.JSON(http.StatusOK, res)
}
