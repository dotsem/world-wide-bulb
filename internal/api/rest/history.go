package rest

import (
	"net/http"
	"strconv"
	"time"

	"world-wide-bulb/internal/store"

	"github.com/gin-gonic/gin"
)

// HistoryItem represents a public toggle history entry.
type HistoryItem struct {
	ID        int64     `json:"id"`
	State     bool      `json:"state"`
	Reason    string    `json:"reason,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// GetHistory returns recent toggle events from the database with position-pointer pagination.
func (h *Handler) GetHistory(c *gin.Context) {
	limit := int64(20)
	if lStr := c.Query("limit"); lStr != "" {
		if parsed, err := strconv.ParseInt(lStr, 10, 64); err == nil && parsed > 0 {
			limit = min(parsed, 100)
		}
	}

	var before int64
	if bStr := c.Query("before"); bStr != "" {
		if parsed, err := strconv.ParseInt(bStr, 10, 64); err == nil && parsed > 0 {
			before = parsed
		}
	}

	var toggles []store.Toggle
	var err error

	if before > 0 {
		toggles, err = h.queries.GetTogglesBefore(c.Request.Context(), store.GetTogglesBeforeParams{
			ID:    before,
			Limit: limit,
		})
	} else {
		toggles, err = h.queries.GetRecentToggles(c.Request.Context(), limit)
	}

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

	var nextCursor int64
	hasMore := false
	if len(toggles) > 0 {
		nextCursor = toggles[len(toggles)-1].ID
		hasMore = int64(len(toggles)) == limit
	}

	c.JSON(http.StatusOK, gin.H{
		"toggles":     items,
		"next_cursor": nextCursor,
		"has_more":    hasMore,
	})
}
