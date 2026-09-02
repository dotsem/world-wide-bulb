package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strings"
	"world-wide-bulb/internal/bulb"

	goaway "github.com/TwiN/go-away"
	"github.com/gin-gonic/gin"
)

// TODO: is there a beter way to do this?
var urlPattern = regexp.MustCompile(`(?i)(https?://|ftp://|www\.|\b[a-z0-9-]+(?:\.[a-z0-9-]+)*\.(com|net|org|io|xyz|ai|co|app|dev|me|info|biz|tv|cc|gg|ly|be|de|uk|ru|cn|nl|eu|site|online|top|link|store|shop|live|tech|space|fun)\b)`)

// ReasonRequest defines the expected JSON payload for attaching a reason to a toggle.
type ReasonRequest struct {
	ID     string `json:"id" binding:"required,uuid"`
	Reason string `json:"reason" binding:"required,max=60"`
}

// PostReason attaches a text reason to an existing toggle by its UUID.
func (h *Handler) PostReason(c *gin.Context) {
	var req ReasonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{errKey: "invalid payload (valid uuid and reason max 60 chars required)"})
		return
	}

	trimmedReason := strings.TrimSpace(req.Reason)
	if trimmedReason == "" {
		c.JSON(http.StatusBadRequest, gin.H{errKey: "reason cannot be empty or whitespace only"})
		return
	}

	if urlPattern.MatchString(trimmedReason) {
		c.JSON(http.StatusBadRequest, gin.H{errKey: "urls are not allowed in reasons"})
		return
	}

	if goaway.IsProfane(trimmedReason) {
		c.JSON(http.StatusBadRequest, gin.H{errKey: "profanity is not allowed in reasons"})
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
		"toggle_id": toggle.ID,
		"reason":    trimmedReason,
	}
	if payload, err := json.Marshal(wsMsg); err == nil {
		h.hub.Broadcast(payload)
		h.broker.Broadcast("reason_updated", payload)
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
