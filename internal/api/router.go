// Package api initializes HTTP routing for REST and WebSocket handlers.
package api

import (
	"io/fs"
	"net/http"
	"world-wide-bulb/internal/api/middleware"
	"world-wide-bulb/internal/api/rest"
	"world-wide-bulb/internal/api/ws"

	"github.com/gin-gonic/gin"
)

// NewRouter creates and initializes the Gin engine with all API routes and embedded static frontend.
// Any toggleGuards run before the toggle handler and may abort the request.
func NewRouter(restH *rest.Handler, wsH *ws.Handler, staticFS fs.FS, isProd bool, allowedHosts []string, toggleGuards ...gin.HandlerFunc) *gin.Engine {
	if isProd {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	_ = r.SetTrustedProxies(nil)

	publicLimiter := middleware.NewIPRateLimiter(100, 30)
	webLimiter := middleware.NewIPRateLimiter(20, 5)

	v1 := r.Group("/api/v1")
	{
		publicCORS := middleware.PublicCORS()
		public := v1.Group("")
		public.Use(publicCORS, middleware.RateLimit(publicLimiter))
		{
			public.GET("/state", restH.GetState)
			public.GET("/history", restH.GetHistory)
			public.GET("/events", restH.StreamEvents)

			public.OPTIONS("/state", publicCORS)
			public.OPTIONS("/history", publicCORS)
			public.OPTIONS("/events", publicCORS)
		}
		webCORS := middleware.WebCORS(isProd, allowedHosts)
		web := v1.Group("")
		web.Use(webCORS, middleware.RateLimit(webLimiter))
		{
			web.POST("/toggle", append(toggleGuards, restH.PostToggle)...)
			web.POST("/reason", restH.PostReason)
			web.POST("/verify", append(toggleGuards, verified)...)

			web.OPTIONS("/toggle", webCORS)
			web.OPTIONS("/reason", webCORS)
			web.OPTIONS("/verify", webCORS)
		}
	}

	wsLimiter := middleware.NewIPRateLimiter(30, 5)
	r.GET("/ws", middleware.RateLimit(wsLimiter), wsH.ServeWS)

	if staticFS != nil {
		r.NoRoute(ServeStatic(staticFS))
	}

	return r
}

// verified lets the frontend pass the toggle guards ahead of time, so a click does not wait on them.
func verified(c *gin.Context) {
	c.Status(http.StatusNoContent)
}
