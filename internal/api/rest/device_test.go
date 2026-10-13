package rest

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"world-wide-bulb/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetOrCreateDeviceID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	hasher := utils.NewHasher("test_salt")

	t.Run("returns ip:deviceID when valid signed device_id cookie is present", func(t *testing.T) {
		h := NewHandler(nil, nil, nil, nil, hasher, false)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		uuidVal := "existing-device-uuid"
		signedCookie := uuidVal + "." + hasher.Sign(uuidVal)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = "192.168.1.100:12345"
		req.AddCookie(&http.Cookie{
			Name:     deviceCookieName,
			Value:    signedCookie,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
		})
		c.Request = req

		primary, secondary := h.getOrCreateDeviceID(c)

		assert.Equal(t, "192.168.1.100:"+uuidVal, primary)
		assert.Empty(t, secondary)
		assert.Empty(t, w.Header().Get("Set-Cookie"))
	})

	t.Run("rejects forged or unsigned cookie and falls back to IP", func(t *testing.T) {
		h := NewHandler(nil, nil, nil, nil, hasher, false)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = "192.168.1.100:12345"
		req.AddCookie(&http.Cookie{
			Name:     deviceCookieName,
			Value:    "forged-uuid-without-signature",
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
		})
		c.Request = req

		primary, secondary := h.getOrCreateDeviceID(c)

		assert.Equal(t, "192.168.1.100", primary)
		assert.NotEmpty(t, secondary)
		assert.NotEmpty(t, w.Header().Get("Set-Cookie"))
	})

	t.Run("generates new device_id and sets signed cookie when cookie is missing", func(t *testing.T) {
		h := NewHandler(nil, nil, nil, nil, hasher, false)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = "192.168.1.100:12345"
		c.Request = req

		primary, secondary := h.getOrCreateDeviceID(c)

		assert.Equal(t, "192.168.1.100", primary)
		require.True(t, len(secondary) > len("192.168.1.100:"))
		assert.Equal(t, "192.168.1.100:", secondary[:len("192.168.1.100:")])

		newID := secondary[len("192.168.1.100:"):]
		_, err := uuid.Parse(newID)
		require.NoError(t, err)

		setCookie := w.Header().Get("Set-Cookie")
		require.Contains(t, setCookie, "device_id="+newID+"."+hasher.Sign(newID))
		require.Contains(t, setCookie, "Max-Age=31536000")
		require.Contains(t, setCookie, "HttpOnly")
		require.Contains(t, setCookie, "SameSite=Lax")
		require.NotContains(t, setCookie, "Secure")
	})

	t.Run("sets secure cookie in production mode", func(t *testing.T) {
		h := NewHandler(nil, nil, nil, nil, hasher, true)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = "192.168.1.100:12345"
		c.Request = req

		_, _ = h.getOrCreateDeviceID(c)

		setCookie := w.Header().Get("Set-Cookie")
		require.Contains(t, setCookie, "Secure")
	})

	t.Run("generates new ID when cookie is empty string", func(t *testing.T) {
		h := NewHandler(nil, nil, nil, nil, nil, false)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = "10.0.0.1:12345"
		req.AddCookie(&http.Cookie{
			Name:     deviceCookieName,
			Value:    "",
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
		})
		c.Request = req

		primary, secondary := h.getOrCreateDeviceID(c)

		assert.Equal(t, "10.0.0.1", primary)
		assert.True(t, len(secondary) > len("10.0.0.1:"))
		assert.Contains(t, w.Header().Get("Set-Cookie"), "device_id=")
	})
}
