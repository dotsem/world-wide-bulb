// Package middleware provides HTTP middlewares for the Gin router.
package middleware

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

// PublicCORS returns a middleware permitting open cross-origin access for public endpoints without credentials.
func PublicCORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "GET, OPTIONS, HEAD")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

// WebCORS returns a middleware that restricts cross-origin access to allowed hosts in production.
func WebCORS(isProd bool, allowedHosts []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			if isProd {
				u, err := url.Parse(origin)
				if err != nil {
					c.AbortWithStatus(http.StatusForbidden)
					return
				}
				hostname := u.Hostname()
				allowed := false
				for _, host := range allowedHosts {
					if strings.EqualFold(hostname, host) {
						allowed = true
						break
					}
				}
				if !allowed {
					c.AbortWithStatus(http.StatusForbidden)
					return
				}
			}
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
			c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		}
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
