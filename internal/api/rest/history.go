package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetHistory(c *gin.Context) {
	toggles, err := h.queries.GetRecentToggles(c.Request.Context(), 100)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"toggles": toggles,
	})
}
