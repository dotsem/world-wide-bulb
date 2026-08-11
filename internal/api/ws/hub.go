package ws

import (
	"log/slog"
)

// Hub coordinates active WebSocket clients and broadcasts state payloads.
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

// NewHub initializes and runs a new WebSocket Hub event loop.
func NewHub() *Hub {
	hub := &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
	go hub.Run()
	return hub
}

// Run executes the multiplexing event loop for the Hub.
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
					slog.Warn("evicted slow client")
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
