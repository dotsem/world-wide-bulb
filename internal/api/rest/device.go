package rest

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	deviceCookieName   = "device_id"
	deviceCookieMaxAge = 31536000 // 1 year
)

// getOrCreateDeviceID returns (primaryKey, secondaryKey).
// If the request presents a valid HMAC-signed device_id cookie, returns (IP:deviceID, "").
// If no cookie is present or the signature is invalid, sets a signed Set-Cookie and returns (IP, IP:newDeviceID).
func (h *Handler) getOrCreateDeviceID(c *gin.Context) (string, string) {
	ip := c.ClientIP()
	hasher := h.hasher

	if rawCookie, err := c.Cookie(deviceCookieName); err == nil && rawCookie != "" {
		if id, sig, ok := strings.Cut(rawCookie, "."); ok && id != "" && sig != "" {
			if hasher.Verify(id, sig) {
				return ip + ":" + id, ""
			}
		}
	}

	newID := uuid.NewString()
	sig := hasher.Sign(newID)
	signedCookie := newID + "." + sig

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(deviceCookieName, signedCookie, deviceCookieMaxAge, "/", "", h.isProd, true)
	return ip, ip + ":" + newID
}
