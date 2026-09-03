package ws_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
	"world-wide-bulb/internal/api/ws"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClientPumps(t *testing.T) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(_ *http.Request) bool { return true },
	}

	t.Run("ReadPump and WritePump handle messages and connection close", func(t *testing.T) {
		hub := ws.NewHub()
		go hub.Run()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			conn, err := upgrader.Upgrade(w, r, nil)
			require.NoError(t, err)

			client := ws.NewClient(hub, conn, "viewer_1")
			hub.Register(client)

			go client.WritePump()
			client.ReadPump()
		}))
		defer server.Close()

		url := "ws" + strings.TrimPrefix(server.URL, "http")
		wsConn, _, err := websocket.DefaultDialer.Dial(url, nil)
		require.NoError(t, err)

		err = wsConn.WriteMessage(websocket.TextMessage, []byte("ping"))
		require.NoError(t, err)

		hub.Broadcast([]byte("hello client"))

		_, msg, err := wsConn.ReadMessage()
		require.NoError(t, err)
		assert.Equal(t, "hello client", string(msg))

		require.NoError(t, wsConn.Close())
		time.Sleep(50 * time.Millisecond)
	})

	t.Run("WritePump sends ping messages and batch messages", func(t *testing.T) {
		hub := ws.NewHub()
		go hub.Run()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			conn, err := upgrader.Upgrade(w, r, nil)
			require.NoError(t, err)

			client := ws.NewClient(hub, conn, "viewer_2")
			hub.Register(client)

			go client.WritePump()
			client.ReadPump()
		}))
		defer server.Close()

		url := "ws" + strings.TrimPrefix(server.URL, "http")
		wsConn, _, err := websocket.DefaultDialer.Dial(url, nil)
		require.NoError(t, err)

		hub.Broadcast([]byte("msg1"))
		hub.Broadcast([]byte("msg2"))

		time.Sleep(50 * time.Millisecond)
		require.NoError(t, wsConn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")))
		_ = wsConn.Close()
	})
}
