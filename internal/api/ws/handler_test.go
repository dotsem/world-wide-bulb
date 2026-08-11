package ws_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
	"world-wide-bulb/internal/api/ws"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOriginValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	hub := ws.NewHub()

	t.Run("dev mode allows any origin", func(t *testing.T) {
		handler := ws.NewHandler(hub, false, nil)
		r := gin.New()
		r.GET("/ws", handler.ServeWS)
		s := httptest.NewServer(r)
		defer s.Close()

		u := "ws" + strings.TrimPrefix(s.URL, "http") + "/ws"
		header := make(http.Header)
		header.Set("Origin", "http://untrusted-origin.local")

		conn, resp, err := websocket.DefaultDialer.Dial(u, header)
		require.NoError(t, err)
		assert.Equal(t, http.StatusSwitchingProtocols, resp.StatusCode)
		defer conn.Close()
	})

	t.Run("prod mode rejects unauthorized origin with 403", func(t *testing.T) {
		handler := ws.NewHandler(hub, true, []string{"my-site.com"})
		r := gin.New()
		r.GET("/ws", handler.ServeWS)
		s := httptest.NewServer(r)
		defer s.Close()

		u := "ws" + strings.TrimPrefix(s.URL, "http") + "/ws"
		header := make(http.Header)
		header.Set("Origin", "http://evil.com")

		_, resp, err := websocket.DefaultDialer.Dial(u, header)
		assert.Error(t, err)
		if resp != nil {
			assert.Equal(t, http.StatusForbidden, resp.StatusCode)
		}
	})

	t.Run("prod mode allows matching origin with port", func(t *testing.T) {
		handler := ws.NewHandler(hub, true, []string{"my-site.com"})
		r := gin.New()
		r.GET("/ws", handler.ServeWS)
		s := httptest.NewServer(r)
		defer s.Close()

		u := "ws" + strings.TrimPrefix(s.URL, "http") + "/ws"
		header := make(http.Header)
		header.Set("Origin", "https://my-site.com:8443")

		conn, resp, err := websocket.DefaultDialer.Dial(u, header)
		require.NoError(t, err)
		assert.Equal(t, http.StatusSwitchingProtocols, resp.StatusCode)
		defer conn.Close()
	})
}

func TestWebSocketBroadcast(t *testing.T) {
	gin.SetMode(gin.TestMode)
	hub := ws.NewHub()
	handler := ws.NewHandler(hub, false, nil)

	r := gin.New()
	r.GET("/ws", handler.ServeWS)
	s := httptest.NewServer(r)
	defer s.Close()

	u := "ws" + strings.TrimPrefix(s.URL, "http") + "/ws"
	conn, _, err := websocket.DefaultDialer.Dial(u, nil)
	require.NoError(t, err)
	defer conn.Close()

	time.Sleep(20 * time.Millisecond)

	testPayload := []byte(`{"state":true,"reason":"test broadcast"}`)
	hub.Broadcast(testPayload)

	_ = conn.SetReadDeadline(time.Now().Add(1 * time.Second))
	_, msg, err := conn.ReadMessage()
	require.NoError(t, err)
	assert.Equal(t, testPayload, msg)
}
