package rest

import (
	"errors"
	"net/http"
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

	err := h.engine.UpdateReason(c.Request.Context(), req.ID, req.Reason)
	if err != nil {
		if errors.Is(err, bulb.ErrInvalidUUID) {
			c.JSON(http.StatusBadRequest, gin.H{errKey: "invalid uuid format"})
			return
		}
		if errors.Is(err, bulb.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{errKey: "toggle not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{errKey: errInternalServer})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
