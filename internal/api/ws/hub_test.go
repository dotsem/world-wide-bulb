package ws

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHub(t *testing.T) {
	t.Run("broadcasts payload to all registered clients", func(t *testing.T) {
		hub := NewHub()
		c1 := NewClient(hub, nil)
		c2 := NewClient(hub, nil)

		hub.Register(c1)
		hub.Register(c2)
		time.Sleep(10 * time.Millisecond)

		msg := []byte("hello")
		hub.Broadcast(msg)

		select {
		case received := <-c1.send:
			assert.Equal(t, msg, received)
		case <-time.After(100 * time.Millisecond):
			t.Fatal("client 1 did not receive broadcast")
		}

		select {
		case received := <-c2.send:
			assert.Equal(t, msg, received)
		case <-time.After(100 * time.Millisecond):
			t.Fatal("client 2 did not receive broadcast")
		}
	})

	t.Run("evicts slow client with full send buffer without blocking others", func(t *testing.T) {
		hub := NewHub()
		slowClient := &Client{hub: hub, send: make(chan []byte, 1)}
		fastClient := &Client{hub: hub, send: make(chan []byte, 10)}

		hub.Register(slowClient)
		hub.Register(fastClient)
		time.Sleep(10 * time.Millisecond)

		slowClient.send <- []byte("backlogged")

		hub.Broadcast([]byte("new message"))
		time.Sleep(10 * time.Millisecond)

		select {
		case received := <-fastClient.send:
			assert.Equal(t, []byte("new message"), received)
		case <-time.After(100 * time.Millisecond):
			t.Fatal("fast client was blocked")
		}

		<-slowClient.send
		_, ok := <-slowClient.send
		assert.False(t, ok, "slow client send channel should be closed")
	})

	t.Run("unregisters client cleanly", func(t *testing.T) {
		hub := NewHub()
		client := NewClient(hub, nil)

		hub.Register(client)
		time.Sleep(10 * time.Millisecond)

		hub.Unregister(client)
		time.Sleep(10 * time.Millisecond)

		_, ok := <-client.send
		assert.False(t, ok, "unregistered client send channel should be closed")
	})
}
