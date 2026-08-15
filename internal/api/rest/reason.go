package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"world-wide-bulb/internal/bulb"

	"github.com/gin-gonic/gin"
)

// ReasonRequest defines the expected JSON payload for attaching a reason to a toggle.
type ReasonRequest struct {
	ID     string `json:"id" binding:"required,uuid"`
	Reason string `json:"reason" binding:"required,max=100"`
}

// PostReason attaches a text reason to an existing toggle by its UUID.
func (h *Handler) PostReason(c *gin.Context) {
	var req ReasonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{errKey: "invalid payload (valid uuid and reason max 100 chars required)"})
		return
	}

	trimmedReason := strings.TrimSpace(req.Reason)
	if trimmedReason == "" {
		c.JSON(http.StatusBadRequest, gin.H{errKey: "reason cannot be empty or whitespace only"})
		return
	}

	toggle, err := h.engine.UpdateReason(c.Request.Context(), req.ID, trimmedReason)
	if err != nil {
		if errors.Is(err, bulb.ErrReasonAlreadySet) {
			c.JSON(http.StatusBadRequest, gin.H{errKey: "reason already set for this toggle"})
			return
		}
		if errors.Is(err, bulb.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{errKey: "toggle not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{errKey: errInternalServer})
		return
	}

	wsMsg := gin.H{
		"type":      "reason_updated",
		"id":        req.ID,
		"toggle_id": toggle.ID,
		"reason":    trimmedReason,
	}
	if payload, err := json.Marshal(wsMsg); err == nil {
		h.hub.Broadcast(payload)
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
