package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetState(c *gin.Context) {
	state := h.engine.GetState()

	c.JSON(http.StatusOK, gin.H{
		"state": state,
	})
}
