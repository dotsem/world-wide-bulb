package api_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"
	"world-wide-bulb/internal/api"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestServeStatic(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("serves html extension fallback", func(t *testing.T) {
		fsys := fstest.MapFS{
			"about.html": &fstest.MapFile{Data: []byte("<h1>About</h1>")},
		}
		router := gin.New()
		router.Use(api.ServeStatic(fsys))

		req := httptest.NewRequest(http.MethodGet, "/about", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "<h1>About</h1>", rec.Body.String())
	})

	t.Run("sets long cache control for immutable assets", func(t *testing.T) {
		fsys := fstest.MapFS{
			"_app/immutable/bundle.js": &fstest.MapFile{Data: []byte("console.log(1);")},
		}
		router := gin.New()
		router.Use(api.ServeStatic(fsys))

		req := httptest.NewRequest(http.MethodGet, "/_app/immutable/bundle.js", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "public, max-age=31536000, immutable", rec.Header().Get("Cache-Control"))
	})

	t.Run("returns 404 when asset and index.html are missing", func(t *testing.T) {
		fsys := fstest.MapFS{}
		router := gin.New()
		router.Use(api.ServeStatic(fsys))

		req := httptest.NewRequest(http.MethodGet, "/missing", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Contains(t, rec.Body.String(), "404 page not found")
	})

	t.Run("serves index.html for root path request", func(t *testing.T) {
		fsys := fstest.MapFS{
			"index.html": &fstest.MapFile{Data: []byte("<html>Home</html>")},
		}
		router := gin.New()
		router.Use(api.ServeStatic(fsys))

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "no-cache, must-revalidate", rec.Header().Get("Cache-Control"))
		assert.Equal(t, "<html>Home</html>", rec.Body.String())
	})
}
