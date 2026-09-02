// Package sse provides a fan-out event broker for Server-Sent Events subscribers.
package sse

import "sync"

// Event represents an SSE event with a name and byte payload.
type Event struct {
	Name string
	Data []byte
}

// NewEvent creates a new event.
func NewEvent(name string, data []byte) Event {
	return Event{Name: name, Data: data}
}

// Broker manages active SSE client channels and dispatches broadcast payloads.
type Broker struct {
	mu      sync.RWMutex
	clients map[chan Event]struct{}
}

// NewBroker initializes an empty SSE broker.
func NewBroker() *Broker {
	return &Broker{
		clients: make(map[chan Event]struct{}),
	}
}

// Subscribe returns a channel that receives all broadcasted events.
// Clients should close the returned channel when they unsubscribe.
func (b *Broker) Subscribe() chan Event {
	b.mu.Lock()
	defer b.mu.Unlock()
	ch := make(chan Event, 16)
	b.clients[ch] = struct{}{}
	return ch
}

// Unsubscribe closes and removes a client channel from the broker.
func (b *Broker) Unsubscribe(ch chan Event) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if _, ok := b.clients[ch]; ok {
		delete(b.clients, ch)
		close(ch)
	}
}

// Broadcast sends a message payload to all active subscriber channels.
func (b *Broker) Broadcast(eventName string, msg []byte) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for ch := range b.clients {
		select {
		case ch <- NewEvent(eventName, msg):
		default:
			// ponytail: drop slow client if buffer fills; upgrade to eviction queue if needed
		}
	}
}
