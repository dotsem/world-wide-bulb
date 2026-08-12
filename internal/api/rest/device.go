package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	deviceCookieName   = "device_id"
	deviceCookieMaxAge = 31536000 // 1 year
)

// getOrCreateDeviceID reads an HttpOnly device cookie if present, or sets a new one on the response.
// Returns (deviceID, true) if request presented a cookie, or (newDeviceID, false) if no cookie was sent.
func (h *Handler) getOrCreateDeviceID(c *gin.Context) (string, bool) {
	deviceID, err := c.Cookie(deviceCookieName)
	if err == nil && deviceID != "" {
		return deviceID, true
	}

	newID := uuid.NewString()
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(deviceCookieName, newID, deviceCookieMaxAge, "/", "", false, true)
	return newID, false
}
