// Package sse provides a fan-out event broker for Server-Sent Events subscribers.
package sse

import "sync"

// Broker manages active SSE client channels and dispatches broadcast payloads.
type Broker struct {
	mu      sync.RWMutex
	clients map[chan []byte]struct{}
}

// NewBroker initializes an empty SSE broker.
func NewBroker() *Broker {
	return &Broker{
		clients: make(map[chan []byte]struct{}),
	}
}

// Subscribe returns a channel that receives all broadcasted events.
// Clients should close the returned channel when they unsubscribe.
func (b *Broker) Subscribe() chan []byte {
	b.mu.Lock()
	defer b.mu.Unlock()
	ch := make(chan []byte, 16)
	b.clients[ch] = struct{}{}
	return ch
}

// Unsubscribe closes and removes a client channel from the broker.
func (b *Broker) Unsubscribe(ch chan []byte) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if _, ok := b.clients[ch]; ok {
		delete(b.clients, ch)
		close(ch)
	}
}

// Broadcast sends a message payload to all active subscriber channels.
func (b *Broker) Broadcast(msg []byte) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for ch := range b.clients {
		select {
		case ch <- msg:
		default:
			// ponytail: drop slow client if buffer fills; upgrade to eviction queue if needed
		}
	}
}
