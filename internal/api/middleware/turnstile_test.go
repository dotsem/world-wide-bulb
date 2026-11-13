package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"world-wide-bulb/internal/api/middleware"
	"world-wide-bulb/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHuman(t *testing.T) {
	gin.SetMode(gin.TestMode)

	verify := func(_ context.Context, token string) bool { return token == "good" }
	r := gin.New()
	r.POST("/toggle", middleware.Human("site-key", utils.NewHasher("salt"), false, verify), func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	do := func(remoteAddr, token, humanCookie string) *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodPost, "/toggle", nil)
		req.RemoteAddr = remoteAddr
		if token != "" {
			req.Header.Set("X-Turnstile-Token", token)
		}
		if humanCookie != "" {
			req.Header.Set("Cookie", "human="+humanCookie)
		}
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		return rec
	}

	t.Run("rejects a bare request and advertises the site key", func(t *testing.T) {
		rec := do("1.1.1.1:1000", "", "")
		assert.Equal(t, http.StatusForbidden, rec.Code)
		assert.Contains(t, rec.Body.String(), `"turnstile_site_key":"site-key"`)
	})

	t.Run("rejects an invalid token", func(t *testing.T) {
		assert.Equal(t, http.StatusForbidden, do("1.1.1.1:1000", "bad", "").Code)
	})

	t.Run("valid token issues a cookie that only works from the same IP", func(t *testing.T) {
		rec := do("1.1.1.1:1000", "good", "")
		require.Equal(t, http.StatusOK, rec.Code)
		cookies := rec.Result().Cookies()
		require.Len(t, cookies, 1)
		assert.True(t, cookies[0].HttpOnly)

		assert.Equal(t, http.StatusOK, do("1.1.1.1:2000", "", cookies[0].Value).Code)
		assert.Equal(t, http.StatusForbidden, do("2.2.2.2:1000", "", cookies[0].Value).Code)
	})

	t.Run("each solve on the same IP gets its own ID", func(t *testing.T) {
		first := strings.Split(do("1.1.1.1:1000", "good", "").Result().Cookies()[0].Value, ".")[0]
		second := strings.Split(do("1.1.1.1:1000", "good", "").Result().Cookies()[0].Value, ".")[0]
		assert.NotEmpty(t, first)
		assert.NotEqual(t, first, second)
	})

	t.Run("rejects a tampered or expired cookie", func(t *testing.T) {
		value := do("1.1.1.1:1000", "good", "").Result().Cookies()[0].Value
		parts := strings.Split(value, ".")
		require.Len(t, parts, 3)
		assert.Equal(t, http.StatusForbidden, do("1.1.1.1:1000", "", parts[0]+".9999999999."+parts[2]).Code)
		assert.Equal(t, http.StatusForbidden, do("1.1.1.1:1000", "", "other-id."+parts[1]+"."+parts[2]).Code)

		expired := "id.1." + utils.NewHasher("salt").Sign("human|1.1.1.1|id|1")
		assert.Equal(t, http.StatusForbidden, do("1.1.1.1:1000", "", expired).Code)
	})
}
