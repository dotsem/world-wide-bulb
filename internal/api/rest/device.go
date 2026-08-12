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

// getOrCreateDeviceID returns (primaryKey, secondaryKey).
// If the request presents a device_id cookie, returns (IP:deviceID, "").
// If no cookie is present, sets Set-Cookie and returns (IP, IP:newDeviceID).
func (h *Handler) getOrCreateDeviceID(c *gin.Context) (string, string) {
	ip := c.ClientIP()
	deviceID, err := c.Cookie(deviceCookieName)
	if err == nil && deviceID != "" {
		return ip + ":" + deviceID, ""
	}

	newID := uuid.NewString()
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(deviceCookieName, newID, deviceCookieMaxAge, "/", "", false, true)
	return ip, ip + ":" + newID
}
