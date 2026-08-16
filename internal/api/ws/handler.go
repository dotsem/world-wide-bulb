package ws

import (
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	// ReadBufferSize defines default read buffer size for websockets.
	ReadBufferSize = 1024
	// WriteBufferSize defines default write buffer size for websockets.
	WriteBufferSize = 1024
)

// Handler manages WebSocket upgrades and client registrations.
type Handler struct {
	hub      *Hub
	upgrader websocket.Upgrader
}

// NewHandler creates a configured WebSocket Handler with origin checking.
func NewHandler(hub *Hub, isProd bool, allowedHosts []string) *Handler {
	return &Handler{
		hub: hub,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  ReadBufferSize,
			WriteBufferSize: WriteBufferSize,
			CheckOrigin: func(r *http.Request) bool {
				if isProd {
					origin := r.Header.Get("Origin")
					if origin == "" {
						return false
					}
					u, err := url.Parse(origin)
					if err != nil {
						return false
					}
					hostname := u.Hostname()
					for _, host := range allowedHosts {
						if strings.EqualFold(hostname, host) {
							return true
						}
					}
					return false
				}
				return true
			},
		},
	}
}

// ServeWS handles incoming HTTP GET requests and upgrades them to WebSocket connections.
func (h *Handler) ServeWS(c *gin.Context) {
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		slog.Error("failed to upgrade websocket connection", slog.Any("err", err))
		return
	}
	deviceID, _ := c.Cookie("device_id")
	if deviceID == "" {
		deviceID = c.ClientIP()
	}

	client := NewClient(h.hub, conn, deviceID)
	h.hub.Register(client)
	go client.WritePump()
	go client.ReadPump()
}
