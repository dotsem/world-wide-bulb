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
func NewRouter(restH *rest.Handler, wsH *ws.Handler, staticFS fs.FS, isProd bool) *gin.Engine {
	if isProd {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	_ = r.SetTrustedProxies(nil)
	r.Use(middleware.CORS())

	v1 := r.Group("/api/v1")
	{
		v1.GET("/state", restH.GetState)
		v1.GET("/history", restH.GetHistory)
		v1.POST("/toggle", restH.PostToggle)
		v1.POST("/reason", restH.PostReason)
	}

	r.GET("/ws", wsH.ServeWS)

	if staticFS != nil {
		r.NoRoute(ServeStatic(staticFS))
	}

	return r
}
