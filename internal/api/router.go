// Package api initializes HTTP routing for REST and WebSocket handlers.
package api

import (
	"io/fs"
	"world-wide-bulb/internal/api/middleware"
	"world-wide-bulb/internal/api/rest"
	"world-wide-bulb/internal/api/ws"

	"github.com/gin-gonic/gin"
)

// NewRouter creates and initializes the Gin engine with all API routes and embedded static frontend.
func NewRouter(restH *rest.Handler, wsH *ws.Handler, staticFS fs.FS, isProd bool, allowedHosts []string) *gin.Engine {
	if isProd {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	_ = r.SetTrustedProxies(nil)

	publicLimiter := middleware.NewIPRateLimiter(100, 30)
	webLimiter := middleware.NewIPRateLimiter(20, 5)

	v1 := r.Group("/api/v1")
	{
		public := v1.Group("")
		public.Use(middleware.PublicCORS(), middleware.RateLimit(publicLimiter))
		{
			public.GET("/state", restH.GetState)
			public.GET("/history", restH.GetHistory)
			public.GET("/events", restH.StreamEvents)

			public.OPTIONS("/state", func(_ *gin.Context) {})
			public.OPTIONS("/history", func(_ *gin.Context) {})
			public.OPTIONS("/events", func(_ *gin.Context) {})
		}
		web := v1.Group("")
		web.Use(middleware.WebCORS(isProd, allowedHosts), middleware.RateLimit(webLimiter))
		{
			web.POST("/toggle", restH.PostToggle)
			web.POST("/reason", restH.PostReason)

			web.OPTIONS("/toggle", func(_ *gin.Context) {})
			web.OPTIONS("/reason", func(_ *gin.Context) {})
		}
	}

	wsLimiter := middleware.NewIPRateLimiter(30, 5)
	r.GET("/ws", middleware.RateLimit(wsLimiter), wsH.ServeWS)

	if staticFS != nil {
		r.NoRoute(ServeStatic(staticFS))
	}

	return r
}
