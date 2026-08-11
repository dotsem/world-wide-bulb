package rest

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// HistoryItem represents a public toggle history entry.
type HistoryItem struct {
	ID        int64     `json:"id"`
	State     bool      `json:"state"`
	Reason    string    `json:"reason,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// GetHistory returns recent toggle events from the database.
func (h *Handler) GetHistory(c *gin.Context) {
	toggles, err := h.queries.GetRecentToggles(c.Request.Context(), 100)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	items := make([]HistoryItem, len(toggles))
	for i, t := range toggles {
		items[i] = HistoryItem{
			ID:        t.ID,
			State:     t.State,
			Reason:    t.Reason.String,
			CreatedAt: t.CreatedAt.Time,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"toggles": items,
	})
}
