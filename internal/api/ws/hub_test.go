package ws

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHub(t *testing.T) {
	t.Run("broadcasts payload to all registered clients", func(t *testing.T) {
		hub := NewHub()
		c1 := NewClient(hub, nil, "viewer_1")
		c2 := NewClient(hub, nil, "viewer_2")

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
		slowClient := &Client{hub: hub, send: make(chan []byte, 1), viewerID: "viewer_slow"}
		fastClient := &Client{hub: hub, send: make(chan []byte, 10), viewerID: "viewer_fast"}

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
		client := NewClient(hub, nil, "viewer_1")

		hub.Register(client)
		time.Sleep(10 * time.Millisecond)

		hub.Unregister(client)
		time.Sleep(10 * time.Millisecond)

		_, ok := <-client.send
		assert.False(t, ok, "unregistered client send channel should be closed")
	})

	t.Run("tracks client count accurately with viewer deduplication", func(t *testing.T) {
		hub := NewHub()
		assert.Equal(t, 0, hub.ClientCount())

		c1Tab1 := NewClient(hub, nil, "user_A")
		c1Tab2 := NewClient(hub, nil, "user_A")
		c2Tab1 := NewClient(hub, nil, "user_B")

		hub.Register(c1Tab1)
		hub.Register(c1Tab2)
		time.Sleep(10 * time.Millisecond)
		assert.Equal(t, 1, hub.ClientCount(), "multiple tabs from user_A should count as 1 viewer")

		hub.Register(c2Tab1)
		time.Sleep(10 * time.Millisecond)
		assert.Equal(t, 2, hub.ClientCount(), "user_A + user_B should count as 2 viewers")

		hub.Unregister(c1Tab1)
		time.Sleep(10 * time.Millisecond)
		assert.Equal(t, 2, hub.ClientCount(), "closing one tab from user_A should keep viewer count at 2")

		hub.Unregister(c1Tab2)
		time.Sleep(10 * time.Millisecond)
		assert.Equal(t, 1, hub.ClientCount(), "closing all tabs from user_A should decrement count to 1")

		hub.Unregister(c2Tab1)
		time.Sleep(10 * time.Millisecond)
		assert.Equal(t, 0, hub.ClientCount(), "all viewers disconnected")
	})
}
