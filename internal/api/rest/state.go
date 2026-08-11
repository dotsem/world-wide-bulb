package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetState returns the current state of the bulb.
func (h *Handler) GetState(c *gin.Context) {
	state := h.engine.GetState()

	c.JSON(http.StatusOK, gin.H{
		"state": state,
	})
}
