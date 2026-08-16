// Package api initializes HTTP routing for REST, WebSocket, and static asset handlers.
package api

import (
	"io/fs"
	"net/http"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
)

// ServeStatic creates a gin.HandlerFunc that serves static assets from staticFS with SPA fallback.
func ServeStatic(staticFS fs.FS) gin.HandlerFunc {
	fileServer := http.FileServer(http.FS(staticFS))

	return func(c *gin.Context) {
		reqPath := path.Clean(strings.TrimPrefix(c.Request.URL.Path, "/"))
		if reqPath == "." || reqPath == "" {
			reqPath = "index.html"
		}

		if strings.HasPrefix(reqPath, "_app/immutable/") {
			c.Header("Cache-Control", "public, max-age=31536000, immutable")
		} else if reqPath == "index.html" {
			c.Header("Cache-Control", "no-cache, must-revalidate")
		}

		if stat, err := fs.Stat(staticFS, reqPath); err == nil && !stat.IsDir() {
			fileServer.ServeHTTP(c.Writer, c.Request)
			c.Abort()
			return
		}

		htmlPath := reqPath + ".html"
		if stat, err := fs.Stat(staticFS, htmlPath); err == nil && !stat.IsDir() {
			c.Request.URL.Path = "/" + htmlPath
			fileServer.ServeHTTP(c.Writer, c.Request)
			c.Abort()
			return
		}

		if stat, err := fs.Stat(staticFS, "index.html"); err == nil && !stat.IsDir() {
			c.Header("Cache-Control", "no-cache, must-revalidate")
			indexContent, readErr := fs.ReadFile(staticFS, "index.html")
			if readErr == nil {
				c.Data(http.StatusOK, "text/html; charset=utf-8", indexContent)
				c.Abort()
				return
			}
		}

		c.String(http.StatusNotFound, "404 page not found")
		c.Abort()
	}
}
