package ws

import (
	"encoding/json"
	"log/slog"
	"sync/atomic"
	"time"
)

// ViewerMessage represents a viewer count payload for WebSocket broadcast.
type ViewerMessage struct {
	Type  string `json:"type"`
	Count int    `json:"count"`
}

// Hub coordinates active WebSocket clients and broadcasts state payloads.
type Hub struct {
	clients      map[*Client]bool
	viewers      map[string]int
	broadcast    chan []byte
	register     chan *Client
	unregister   chan *Client
	count        atomic.Int64
	pushInterval time.Duration
}

// NewHub initializes and runs a new WebSocket Hub event loop.
func NewHub() *Hub {
	return newHubWithInterval(1 * time.Second)
}

func newHubWithInterval(interval time.Duration) *Hub {
	hub := &Hub{
		clients:      make(map[*Client]bool),
		viewers:      make(map[string]int),
		broadcast:    make(chan []byte),
		register:     make(chan *Client),
		unregister:   make(chan *Client),
		pushInterval: interval,
	}
	go hub.Run()
	return hub
}

func (h *Hub) removeClient(client *Client) {
	delete(h.clients, client)
	close(client.send)
	if client.viewerID != "" {
		h.viewers[client.viewerID]--
		if h.viewers[client.viewerID] <= 0 {
			delete(h.viewers, client.viewerID)
		}
	}
	h.count.Store(int64(len(h.viewers)))
}

// Run executes the multiplexing event loop for the Hub.
func (h *Hub) Run() {
	ticker := time.NewTicker(h.pushInterval)
	defer ticker.Stop()
	var lastPushedCount int

	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			if client.viewerID != "" {
				h.viewers[client.viewerID]++
			}
			h.count.Store(int64(len(h.viewers)))
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				h.removeClient(client)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					h.removeClient(client)
					slog.Warn("evicted slow client")
				}
			}
		case <-ticker.C:
			current := len(h.viewers)
			if current != lastPushedCount {
				lastPushedCount = current
				payload, err := json.Marshal(ViewerMessage{Type: "viewer_count", Count: current})
				if err == nil {
					for client := range h.clients {
						select {
						case client.send <- payload:
						default:
							h.removeClient(client)
							slog.Warn("evicted slow client")
						}
					}
				}
			}
		}
	}
}

// Broadcast sends a raw payload to all registered clients.
func (h *Hub) Broadcast(payload []byte) {
	h.broadcast <- payload
}

// Register adds a client connection to the hub.
func (h *Hub) Register(c *Client) {
	h.register <- c
}

// Unregister removes a client connection from the hub.
func (h *Hub) Unregister(c *Client) {
	h.unregister <- c
}

// ClientCount returns the current number of active connected clients.
func (h *Hub) ClientCount() int {
	return int(h.count.Load())
}
