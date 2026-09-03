package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"world-wide-bulb/internal/api/middleware"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestIPRateLimiter(t *testing.T) {
	t.Run("allows requests up to burst capacity and then blocks", func(t *testing.T) {
		limiter := middleware.NewIPRateLimiter(60, 3)

		allowed, _ := limiter.Allow("192.168.1.1")
		assert.True(t, allowed)
		allowed, _ = limiter.Allow("192.168.1.1")
		assert.True(t, allowed)
		allowed, _ = limiter.Allow("192.168.1.1")
		assert.True(t, allowed)

		allowed, retryAfter := limiter.Allow("192.168.1.1")
		assert.False(t, allowed)
		assert.Greater(t, retryAfter, time.Duration(0))
	})

	t.Run("isolates rate limits between distinct client IPs", func(t *testing.T) {
		limiter := middleware.NewIPRateLimiter(60, 1)

		allowed, _ := limiter.Allow("10.0.0.1")
		assert.True(t, allowed)
		allowed, _ = limiter.Allow("10.0.0.1")
		assert.False(t, allowed)

		allowed, _ = limiter.Allow("10.0.0.2")
		assert.True(t, allowed)
	})

	t.Run("replenishes tokens over time", func(t *testing.T) {
		limiter := middleware.NewIPRateLimiter(1200, 1) // 20 tokens per second = 1 token per 50ms

		allowed, _ := limiter.Allow("10.0.0.3")
		assert.True(t, allowed)
		allowed, _ = limiter.Allow("10.0.0.3")
		assert.False(t, allowed)

		time.Sleep(60 * time.Millisecond)

		allowed, _ = limiter.Allow("10.0.0.3")
		assert.True(t, allowed)
	})
}

func TestRateLimitMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("aborts with 429 and Retry-After header when rate limit exceeded", func(t *testing.T) {
		limiter := middleware.NewIPRateLimiter(60, 1)
		r := gin.New()
		r.Use(middleware.RateLimit(limiter))
		r.GET("/ping", func(c *gin.Context) { c.String(http.StatusOK, "pong") })

		req1 := httptest.NewRequest(http.MethodGet, "/ping", nil)
		rec1 := httptest.NewRecorder()
		r.ServeHTTP(rec1, req1)
		assert.Equal(t, http.StatusOK, rec1.Code)

		req2 := httptest.NewRequest(http.MethodGet, "/ping", nil)
		rec2 := httptest.NewRecorder()
		r.ServeHTTP(rec2, req2)
		assert.Equal(t, http.StatusTooManyRequests, rec2.Code)
		assert.NotEmpty(t, rec2.Header().Get("Retry-After"))
		assert.Contains(t, rec2.Body.String(), "rate limit exceeded")
	})
}
