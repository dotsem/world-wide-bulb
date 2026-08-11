package api

import (
	"world-wide-bulb/internal/api/rest"
	"world-wide-bulb/internal/api/ws"

	"github.com/gin-gonic/gin"
)

func NewRouter(restH *rest.Handler, wsH *ws.Handler) *gin.Engine {
	r := gin.Default()
	_ = r.SetTrustedProxies(nil)

	v1 := r.Group("/api/v1")
	{
		v1.GET("/state", restH.GetState)
		v1.GET("/history", restH.GetHistory)
		v1.POST("/toggle", restH.PostToggle)
	}

	r.GET("/ws", wsH.ServeWS)

	return r
}
