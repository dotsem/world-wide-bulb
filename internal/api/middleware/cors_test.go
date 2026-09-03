package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"world-wide-bulb/internal/api/middleware"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestPublicCORS(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("handles OPTIONS preflight with 204 and wildcard origin", func(t *testing.T) {
		r := gin.New()
		r.Use(middleware.PublicCORS())
		r.GET("/public", func(c *gin.Context) { c.String(http.StatusOK, "ok") })

		req := httptest.NewRequest(http.MethodOptions, "/public", nil)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNoContent, rec.Code)
		assert.Equal(t, "*", rec.Header().Get("Access-Control-Allow-Origin"))
		assert.Empty(t, rec.Header().Get("Access-Control-Allow-Credentials"))
		assert.Contains(t, rec.Header().Get("Access-Control-Allow-Methods"), "GET")
	})

	t.Run("injects wildcard origin on standard GET request", func(t *testing.T) {
		r := gin.New()
		r.Use(middleware.PublicCORS())
		r.GET("/public", func(c *gin.Context) { c.String(http.StatusOK, "ok") })

		req := httptest.NewRequest(http.MethodGet, "/public", nil)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "*", rec.Header().Get("Access-Control-Allow-Origin"))
	})
}

func TestWebCORS(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("non-production allows origin and sets credentials", func(t *testing.T) {
		r := gin.New()
		r.Use(middleware.WebCORS(false, nil))
		r.POST("/toggle", func(c *gin.Context) { c.String(http.StatusOK, "ok") })

		req := httptest.NewRequest(http.MethodPost, "/toggle", nil)
		req.Header.Set("Origin", "http://localhost:5173")
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "http://localhost:5173", rec.Header().Get("Access-Control-Allow-Origin"))
		assert.Equal(t, "true", rec.Header().Get("Access-Control-Allow-Credentials"))
	})

	t.Run("production allows whitelisted host", func(t *testing.T) {
		r := gin.New()
		r.Use(middleware.WebCORS(true, []string{"example.com"}))
		r.POST("/toggle", func(c *gin.Context) { c.String(http.StatusOK, "ok") })

		req := httptest.NewRequest(http.MethodPost, "/toggle", nil)
		req.Header.Set("Origin", "https://EXAMPLE.COM")
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "https://EXAMPLE.COM", rec.Header().Get("Access-Control-Allow-Origin"))
		assert.Equal(t, "true", rec.Header().Get("Access-Control-Allow-Credentials"))
	})

	t.Run("production rejects non-whitelisted host with 403", func(t *testing.T) {
		r := gin.New()
		r.Use(middleware.WebCORS(true, []string{"example.com"}))
		r.POST("/toggle", func(c *gin.Context) { c.String(http.StatusOK, "ok") })

		req := httptest.NewRequest(http.MethodPost, "/toggle", nil)
		req.Header.Set("Origin", "https://evil.com")
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusForbidden, rec.Code)
	})

	t.Run("production rejects malformed origin URL with 403", func(t *testing.T) {
		r := gin.New()
		r.Use(middleware.WebCORS(true, []string{"example.com"}))
		r.POST("/toggle", func(c *gin.Context) { c.String(http.StatusOK, "ok") })

		req := httptest.NewRequest(http.MethodPost, "/toggle", nil)
		req.Header.Set("Origin", "http://%41:8080")
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusForbidden, rec.Code)
	})

	t.Run("passes through without CORS headers when Origin is empty", func(t *testing.T) {
		r := gin.New()
		r.Use(middleware.WebCORS(true, []string{"example.com"}))
		r.POST("/toggle", func(c *gin.Context) { c.String(http.StatusOK, "ok") })

		req := httptest.NewRequest(http.MethodPost, "/toggle", nil)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Empty(t, rec.Header().Get("Access-Control-Allow-Origin"))
	})

	t.Run("handles OPTIONS preflight with 204", func(t *testing.T) {
		r := gin.New()
		r.Use(middleware.WebCORS(true, []string{"example.com"}))
		r.POST("/toggle", func(c *gin.Context) { c.String(http.StatusOK, "ok") })

		req := httptest.NewRequest(http.MethodOptions, "/toggle", nil)
		req.Header.Set("Origin", "https://example.com")
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNoContent, rec.Code)
		assert.Equal(t, "https://example.com", rec.Header().Get("Access-Control-Allow-Origin"))
	})
}
