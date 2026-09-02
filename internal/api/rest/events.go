package rest

import (
	"encoding/json"
	"io"
	"time"

	"github.com/gin-gonic/gin"
)

// StreamEvents serves a public Server-Sent Events stream for real-time bulb state changes.
func (h *Handler) StreamEvents(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("X-Accel-Buffering", "no") // why: revents Nginx/proxy buffer delay

	ch := h.broker.Subscribe()
	defer h.broker.Unsubscribe(ch)

	initialPayload, err := json.Marshal(gin.H{
		"event":  stateKey,
		stateKey: h.engine.GetState(),
	})

	if err == nil {
		c.SSEvent("state_changed", string(initialPayload))
		c.Writer.Flush()
	}

	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	c.Stream(func(_ io.Writer) bool {
		select {
		case <-c.Request.Context().Done():
			return false
		case <-ticker.C:
			c.SSEvent("ping", "")
			return true
		case evt, ok := <-ch:
			if !ok {
				return false
			}
			c.SSEvent(evt.Name, string(evt.Data))
			return true
		}
	})

}
